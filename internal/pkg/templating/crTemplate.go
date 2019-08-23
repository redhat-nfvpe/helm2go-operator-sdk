package templating

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

// CRTemplateConfig is used to template the appropriate CR Scaffold
type CRTemplateConfig struct {
	APIVersion string
	Kind       string
	LowerKind  string
	Spec       CRSpec
}

// CRSpec contains the specific values needed for the CRTemplate
type CRSpec struct {
	placeholder string
}

// NewCRTemplateConfig ...
// TODO : need to implement the function to create an actual spec struct
func NewCRTemplateConfig(apiVersion, kind, lowerKind string, spec CRSpec) *CRTemplateConfig {
	return &CRTemplateConfig{
		APIVersion: apiVersion,
		Kind:       kind,
		LowerKind:  lowerKind,
		Spec:       spec,
	}
}

// Execute renders the template and returns the templated string
func (c *CRTemplateConfig) Execute() (string, error) {
	temp, err := template.New("CRTemplate").Parse(c.GetTemplate())
	if err != nil {
		return "", err
	}

	var wr bytes.Buffer
	err = temp.Execute(&wr, c)
	if err != nil {
		return "", err
	}
	return wr.String(), nil
}

// GetTemplate returns the necessary template for the CR
func (c *CRTemplateConfig) GetTemplate() string {

	return `apiVersion: {{ .Resource.APIVersion }}
	kind: {{ .Resource.Kind }}
	metadata:
	  name: example-{{ .Resource.LowerKind }}
	spec:
	{{- with .Spec }}
	{{ . | indent 2 }}
	{{- else }}
	  # Add fields here
	  size: 3
	{{- end }}
	`
}

//OverwriteCR rewrites the operator CR
func OverwriteCR(outputDir, kind, apiVersion, valuesPath string) bool {
	var err error

	// read values
	bytes, err := ioutil.ReadFile(valuesPath)
	if err != nil {
		log.Printf("error reading values file %s: %v", valuesPath, err)
		return false
	}
	values := string(bytes)
	apiComponents, err := getAPIVersionComponents(apiVersion)
	version := apiComponents.Version
	group := apiComponents.Subdomain

	// overwrite the original file
	outFile := filepath.Join(outputDir, "deploy", "crds", fmt.Sprintf("%s_%s_%s_cr.yaml", group, version, kindToLower(kind)))

	err = appendWrite(outFile, values)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

func appendWrite(filePath string, values string) error {
	err := writeHelper(filePath, values, false)
	return err
}

func writeHelper(filePath, input string, truncate bool) error {
	// open the file with write permissions
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return fmt.Errorf("error opening file %s: %v", filePath, err)
	}
	if truncate {
		// delete original content
		f.Truncate(0)

	}
	defer f.Close()
	// write the new inputted content
	_, err = f.WriteString(input)
	if err != nil {
		return fmt.Errorf("error writing to file %s: %v", filePath, err)
	}

	return nil
}
