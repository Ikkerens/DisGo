package main

// +build ignore

import (
	"bytes"
	"go/format"
	"io/ioutil"
	"os"
	"strings"
	"text/template"

	"github.com/ikkerens/disgo/generate"
	"github.com/slf4go/logger"
)

type eventDeclaration struct {
	Name      string
	EventName string
	Embed     string

	StarTypes  []registeredEventField
	ArrayTypes []registeredEventField
}

type registeredEventField struct {
	FieldName string
	TypeName  string
}

func main() {
	logger.Infof("Generating GO file")
	var tpl = template.Must(template.New("state").Funcs(template.FuncMap{
		"ToLower": strings.ToLower,
	}).Parse(`package disgo

		// Warning: This file has been automatically generated by generate/state/main.go
		// Do NOT make changes here, instead adapt generate/generate.go and run go generate

		import "sync"

		type state struct { {{range .}}
			{{. | ToLower}}s map[Snowflake]*{{.}} {{end}}
		}

		var objects = &state{ {{range .}}
			{{. | ToLower}}s: make(map[Snowflake]*{{.}}), {{end}}
		}

		{{range .}}
			func (s *state) register{{.}}(id identifiableObject) *{{.}} {
				if registered, exists := s.{{. | ToLower}}s[id.ID()]; exists {
					return registered
				} else {
					{{. | ToLower}}, typeOk := id.(*{{.}})
					if !typeOk {
						{{. | ToLower}} = &{{.}}{internal: &internal{{.}}{}, lock: new(sync.RWMutex)}
					}
					s.{{. | ToLower}}s[id.ID()] = {{. | ToLower}}
					return {{. | ToLower}}
				}
			}
		{{end}}
	`))

	var result []byte
	var buf bytes.Buffer
	err := tpl.Execute(&buf, generate.RegisteredTypes)

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

	ioutil.WriteFile("state_generated.go", result, 0644)
}
