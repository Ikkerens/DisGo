package main

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"html/template"
	"io/ioutil"

	"github.com/slf4go/logger"
)

type internalField struct {
	Parent  *internalType
	Name    string
	TypeStr string
}

type internalType struct {
	Name     string
	Exported string
	Fields   []internalField
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
			structDef, ok := typ.Type.(*ast.StructType)

			if ok {
				typeDef := internalType{name, name[8:], make([]internalField, 0)}
				for _, field := range structDef.Fields.List {
					var typeStr string
					ident, isIdent := field.Type.(*ast.Ident)
					if isIdent {
						typeStr = ident.Name
					} else {
						starExpr, isStarExpr := field.Type.(*ast.StarExpr)
						if isStarExpr {
							starName := starExpr.X.(*ast.Ident)
							typeStr = "*" + starName.Name
						} else {
							logger.Errorf("Unknown field type (%T)", field.Type)
							continue
						}
					}

					typeDef.Fields = append(typeDef.Fields, internalField{&typeDef, field.Names[0].Name, typeStr})
				}

				types = append(types, typeDef)
			}
		}
	}

	var tpl = template.Must(template.New("apimodel").Parse(`package disgo

		// Warning: This file has been automatically generated by generate/apimodel/main.go
		// Do NOT make changes here, instead adapt model.go and run go generate

		import "encoding/json"

		{{range .}}
			type {{.Exported}} struct {
				discordObject *{{.Name}}
			}

			func (s *{{.Exported}}) MarshalJSON() ([]byte, error) {
				return json.Marshal(s.discordObject)
			}

			func (s *{{.Exported}}) UnmarshalJSON(b []byte) error {
				return json.Unmarshal(b, &s.discordObject)
			}

			{{range .Fields}}
				func (s *{{.Parent.Exported}}) {{.Name}}() {{.TypeStr}} {
					return s.discordObject.{{.Name}}
				}
			{{end}}
		{{end}}
	`))

	var result []byte
	var buf bytes.Buffer
	err = tpl.Execute(&buf, types)

	if err == nil {
		var formatted []byte
		formatted, err = format.Source(buf.Bytes())
		result = formatted
	}

	if err != nil {
		logger.ErrorE(err)
		return
	}

	ioutil.WriteFile("apimodel.go", result, 0644)
}