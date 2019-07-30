package templating

import (
	"bytes"
	"io/ioutil"
	"reflect"
	"text/template"

	"gopkg.in/yaml.v2"
)

// NewSpecConfig returns a spec config
func NewSpecConfig(name, valuesPath string) *SpecConfig {
	yamlMap, err := getYAMLMap(valuesPath)
	if err != nil {
		panic(err)
	}
	return &SpecConfig{
		Name:          name,
		InstanceValue: expandInstanceValue(yamlMap),
	}
}

// SpecConfig is the YAML spec config
type SpecConfig struct {
	Name          string
	InstanceValue string
}

// GetTemplate returns the necessary template
func (s *SpecConfig) GetTemplate() string {
	structTemplate := `
	type {{ .Name }} struct {
		{{ .InstanceValue }}
	}
	`
	return structTemplate
}

// Execute renders the template and returns the templated string
func (s *SpecConfig) Execute() (string, error) {
	temp, err := template.New("TypeSpec").Parse(s.GetTemplate())
	if err != nil {
		return "", err
	}

	s.Name = s.Name + "Type"

	var wr bytes.Buffer
	err = temp.Execute(&wr, s)
	if err != nil {
		return "", err
	}
	return wr.String(), nil
}

func getYAMLMap(yamlInput string) (map[interface{}]interface{}, error) {
	var output map[interface{}]interface{}
	//read bytes yaml
	bytes, err := ioutil.ReadFile(yamlInput)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(bytes, &output)
	if err != nil {
		return nil, err
	}
	return output, nil
}

func expandInstanceValue(instance map[interface{}]interface{}) string {
	var result string
	for k, v := range instance {
		strK, _ := k.(string)
		if reflect.TypeOf(v).Kind() != reflect.Map {
			result = result + "\t" + strK + " string " + getJSONTag(strK) + "\n"
		} else {
			mapV := v.(map[interface{}]interface{})
			result = result + "\t" + strK + " struct { " + "\n" + "\t" + expandInstanceValue(mapV) + "}" + "\n"
		}
	}
	return result
}

func getJSONTag(input string) string {
	return "`" + "json:" + `"` + input + `,omitempty"` + "`"
}
