package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	// "fmt"
	"gopkg.in/yaml.v2"
	"log"
	"text/template"
)

type Test struct {
	Name             string          `yaml:"name"`
	Instructions     [][]interface{} `yaml:"instructions"`
	InstructionsJSON string
	Tests            []struct {
		Name   string                 `yaml:"name"`
		Result bool                   `yaml:"result"`
		Data   map[string]interface{} `yaml:"context"`
	} `yaml:"tests"`
}

func main() {

	tmpl := `
	package predicator

	import (
		"testing"
		"encoding/json"
	)

	// Test_{{.Name}} -- AUTOGENERATED DO NOT EDIT
	func Test_{{.Name}}(t *testing.T) {
		instructions := [][]interface{}{}
		b := []byte(` + "`{{.InstructionsJSON}}`" + `)
		err := json.Unmarshal(b, &instructions)
		if err != nil {
			t.Fatal()
		}
		tt := []struct{
			Name string
			Result bool
			Data   map[string]interface{} 
		}{
		{{ range .Tests}}
			{
				"{{ .Name }}",
				{{ .Result }},
				{{if .Data }}{{.Data}}{{else}} map[string]interface{}{}{{end}},
			},
		{{end}}
		}
		for _, test := range tt {
			e := NewEvaluator(instructions, test.Data)
			got := e.result()
			if  got != test.Result {
				t.Logf("FAILED %s_%s expected %v got %v", "{{.Name}}", test.Name, test.Result, got)
				t.Fail()
			}
			if e.stack.count > 0 {
				t.Logf("FAILED %s_%s expected empty stack",  "{{.Name}}", test.Name)
				t.Fail()
			}
		}
	}
	`
	compiledTmpl, err := template.New("testtmpl").Parse(tmpl)
	if err != nil {
		log.Fatal("failed to parse template ", err.Error())
	}
	fileData, err := ioutil.ReadDir("evaluator_spec")
	if err != nil {
		log.Fatal("failed to find dir ", err.Error())
	}
	for _, info := range fileData {
		if strings.Contains(info.Name(), ".yml") {
			fileContent, err := ioutil.ReadFile("evaluator_spec/" + info.Name())
			if err != nil {
				log.Fatal("failed to find dir ", err.Error())
			}
			test := Test{}
			err = yaml.Unmarshal([]byte(fileContent), &test)
			if err != nil {
				log.Fatal("failed to parse yaml ", err.Error())
			}
			instructionsJSON := []byte{}
			instructionsJSON, err = json.Marshal(test.Instructions)
			test.InstructionsJSON = string(instructionsJSON)
			if err != nil {
				log.Fatal("failed to convert instruction to json  ", err.Error())
			}
			testFilename := strings.Replace(info.Name(), ".yml", "_test.go", -1)
			f, err := os.Create("../../" + testFilename)
			if err != nil {
				log.Fatal("failed to create test file with name: ", testFilename, " error: ", err.Error())
			}
			defer f.Close()
			err = compiledTmpl.Execute(f, test)
			if err != nil {
				log.Fatal("failed to render template ", err.Error())
			}
		}
	}

}
