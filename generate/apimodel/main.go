package main

// +build ignore

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"sort"
	"text/template"

	"github.com/ikkerens/disgo/generate"
	"github.com/slf4go/logger"
)

type internalField struct {
	Parent  *internalType
	Name    string
	TypeStr string
}

type internalType struct {
	Name      string
	Exported  string
	StateType bool
	Fields    []internalField
}

func main() {
	file, err := parser.ParseFile(token.NewFileSet(), "model.go", nil, 0)
	if err != nil {
		logger.ErrorE(err)
		return
	}

	var types = make([]internalType, 0)
	for name, object := range file.Scope.Objects {
		if object.Kind == ast.Typ {
			typ := object.Decl.(*ast.TypeSpec)

			if len(typ.Name.Name) < 8 || typ.Name.Name[:8] != "internal" {
				continue
			}

			structDef, ok := typ.Type.(*ast.StructType)

			if ok {
				typeDef := internalType{name, name[8:], generate.IsRegisteredType(name[8:]), make([]internalField, 0)}
				logger.Infof("Creating API struct for %s with name %s.", typeDef.Name, typeDef.Exported)
				for _, field := range structDef.Fields.List {
					typeStr, err := determineType(field.Type)
					if err != nil {
						logger.ErrorE(err)
						continue
					}

					iField := internalField{&typeDef, field.Names[0].Name, typeStr}
					typeDef.Fields = append(typeDef.Fields, iField)
					logger.Infof("Adding func %s() with return type %s to %s.", iField.Name, iField.TypeStr, iField.Parent.Exported)
				}

				types = append(types, typeDef)
			}
		}
	}

	sort.SliceStable(types, func(i, j int) bool {
		return types[i].Name < types[j].Name
	})

	logger.Infof("Generating GO file")
	var tpl = template.Must(template.New("apimodel").Parse(`package disgo

		// Warning: This file has been automatically generated by generate/apimodel/main.go
		// Do NOT make changes here, instead adapt model.go and run go generate

		import (
			"encoding/json"
			"time"
		)

		{{range .}}
			// {{.Exported}} is based on the Discord object with the same name.
			// Any fields can be obtained by calling the respective getters.
			type {{.Exported}} struct {
				session *Session
				internal *{{.Name}}
			}

			// MarshalJSON is used to convert this object into its json representation for Discord
			func (s *{{.Exported}}) MarshalJSON() ([]byte, error) {
				return json.Marshal(s.internal)
			}

			// UnmarshalJSON is used to convert json discord objects back into their respective structs
			func (s *{{.Exported}}) UnmarshalJSON(b []byte) error { {{if .StateType}}
				id := IDObject{}
				if err := json.Unmarshal(b, &id); err != nil {
					return err
				}

				*s = *objects.register{{.Exported}}(&id) {{else}}
					s.internal = &{{.Name}}{} {{end}}
				return json.Unmarshal(b, &s.internal)
			}

			{{range .Fields}}
				// {{ .Name}} is used to export the {{.Name}} from this struct.
				func (s *{{.Parent.Exported}}) {{.Name}}() {{.TypeStr}} {
					return s.internal.{{.Name}}
				}
			{{end}}
		{{end}}
	`))

	var result []byte
	var buf bytes.Buffer
	err = tpl.Execute(&buf, types)

	if err == nil {
		var formatted []byte
		logger.Infof("Formatting GO file")
		formatted, err = format.Source(buf.Bytes())
		result = formatted
	}

	if err != nil {
		logger.ErrorE(err)
		os.Exit(1)
	}

	ioutil.WriteFile("model_generated.go", result, 0644)
}

func determineType(field ast.Expr) (typeStr string, err error) {
	switch v := field.(type) {
	case *ast.Ident:
		typeStr = v.Name
	case *ast.StarExpr:
		subType, err := determineType(v.X)
		if err != nil {
			return "", err
		}
		typeStr = "*" + subType
	case *ast.SelectorExpr:
		subType, err := determineType(v.X)
		if err != nil {
			return "", err
		}
		typeStr = subType + "." + v.Sel.Name
	case *ast.ArrayType:
		subType, err := determineType(v.Elt)
		if err != nil {
			return "", err
		}
		typeStr = "[]" + subType
	default:
		err = fmt.Errorf("Unknown field type (%T): %+v", v, v)
	}

	return
}
