package templating

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"text/template"

	"github.com/redhat-nfvpe/helm2go-operator-sdk/internal/resourcecache"
)

// TemplateConfig is used as payload for executing the templating functionality
type TemplateConfig struct {
	ImportStatements []string
	PackageName      string
	Kind             string
	APIVersion       string
	ResourceName     string
	ResourceType     string
	Resource         interface{}
	ResourceJSON     string
	//resource     TemplateResource *need to create a template resource that works with all kubernetes resources*
}

// ControllerWatchFuncConfig ...
type ControllerWatchFuncConfig struct {
	APIVersion            string
	ResourceImportPackage string
	ResourceType          string
	Kind                  string
	LowerKind             string
	OwnerAPIVersion       string
}

// ControllerTemplateConfig ...
// TODO figure out what is needed for this struct
type ControllerTemplateConfig struct {
	Kind            string
	LowerKind       string
	OwnerAPIVersion string
	ImportMap       map[string]string
	ResourceWatches []string
}

// CacheTemplating takes a resource cache and renders its respective template values
func CacheTemplating(rcache *resourcecache.ResourceCache, kind, apiVersion string) map[string]string {
	var t map[string]string
	t = make(map[string]string)
	c := rcache.PrepareCacheForFile()

	for kt, r := range c {
		// TODO need to remove the hardcoded first resourcefunction
		kt = filepath.Join(filepath.Dir(kt), r.PackageName.String(), filepath.Base(kt))
		bytes, err := json.MarshalIndent(r.GetResourceFunctions()[0].Data, "", "\t")
		if err != nil {
			_ = fmt.Errorf("%v", err)
		}
		conf := TemplateConfig{
			ImportStatements: getTemplateImports(&r.GetResourceFunctions()[0].Data),
			PackageName:      string(r.PackageName),
			Kind:             kind,
			ResourceName:     getTemplateResourceTitle(&r.GetResourceFunctions()[0].Data),
			ResourceType:     reflect.TypeOf((r.GetResourceFunctions()[0].Data)).String(),
			Resource:         getTemplateResourceTitle(&r.GetResourceFunctions()[0].Data),
			ResourceJSON:     string(bytes),
			APIVersion:       apiVersion,
		}
		tmpl, err := configDeclarationTemplate(conf)
		t[string(kt)] = tmpl
	}

	return t
}

// ResourceObjectToGoResourceFile ...
func ResourceObjectToGoResourceFile(resource interface{}) (string, error) {

	switch t := reflect.TypeOf(resource); t.String() {
	case "*v1.Deployment":
		return "placeholder", nil
	default:
		return "", errors.New("unsupported kubernetes resource type")
	}
}

func configDeclarationTemplate(c TemplateConfig) (string, error) {

	cwd, err := os.Getwd()
	temp, err := template.New("resourceFunc.tmpl").ParseFiles(filepath.Join(cwd, "pkg", "templating", "resourceFunc.tmpl"))
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

// TemplatesToFiles takes the templated files and writes them to strings in the specified directory
func TemplatesToFiles(templates map[string]string, outputDir string) bool {
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		err = os.MkdirAll(outputDir, 0700)
		if err != nil {
			fmt.Println(err)
			return false
		}
	}

	for filename, file := range templates {
		// create the file to be written to
		f, err := os.Create(filepath.Join(outputDir, filename))
		if err != nil {
			fmt.Println(err)
			return false
		}
		// write the file content to the actual file
		_, err = f.WriteString(file)
		if err != nil {
			fmt.Println(err)
			return false
		}
		f.Close()
	}
	return true
}

// ResourceFileStructure creates the correct file structure based on the spcified kind types
func ResourceFileStructure(rcache *resourcecache.ResourceCache, outputDir string) bool {
	// iterate through the kind types
	for _, r := range *rcache.GetCache() {
		// create the neccessary folder for the kind type
		newOutput := filepath.Join(outputDir, r.PackageName.String())
		if _, err := os.Stat(newOutput); os.IsNotExist(err) {
			err = os.MkdirAll(newOutput, 0700)
			if err != nil {
				fmt.Println(err)
				return false
			}
		}
	}
	return true
}

// OverwriteController takes in a templated file and writes it over the original controller file
func OverwriteController(outputDir, kind, apiVersion string, rcache *resourcecache.ResourceCache) bool {

	var watchFuncs []string
	var wr bytes.Buffer
	var temp *template.Template
	var err error
	var c ControllerTemplateConfig
	var w ControllerWatchFuncConfig

	ownerAPIVersion := getOwnerAPIVersion(apiVersion, kind)
	lowerKind := kindToLowerCamel(kind)

	cwd, err := os.Getwd()
	if err != nil {
		log.Println(err)
		return false
	}

	// TODO I call this twice; I can probably merge the functions
	for _, r := range rcache.PrepareCacheForFile() {
		wr = bytes.Buffer{}
		temp, err = template.New("controllerFunc.tmpl").ParseFiles(filepath.Join(cwd, "pkg", "templating", "controllerFunc.tmpl"))
		if err != nil {
			log.Println(err)
			return false
		}

		resourceType := getTemplateResourceTitle(&r.GetResourceFunctions()[0].Data)

		w = ControllerWatchFuncConfig{
			APIVersion:            apiVersion,
			OwnerAPIVersion:       ownerAPIVersion,
			Kind:                  kind,
			LowerKind:             lowerKind,
			ResourceImportPackage: getResourceImportPackage(resourceType),
			ResourceType:          resourceType,
		}

		err = temp.Execute(&wr, w)
		if err != nil {
			log.Println(err)
			return false
		}

		watchFuncs = append(watchFuncs, wr.String())
	}
	wr = bytes.Buffer{}
	temp, err = template.New("controller.tmpl").ParseFiles(filepath.Join(cwd, "pkg", "templating", "controller.tmpl"))
	c = ControllerTemplateConfig{
		kind,
		lowerKind,
		ownerAPIVersion,
		getImportMap(outputDir, kind, apiVersion),
		watchFuncs,
	}
	err = temp.Execute(&wr, c)
	if err != nil {
		log.Println(err)
		return false
	}

	// overwrite original file
	outFile := filepath.Join(outputDir, "pkg", "controller", kindToLower(kind), fmt.Sprintf("%s_controller.go", kindToLower(kind)))
	f, err := os.OpenFile(outFile, os.O_WRONLY, 0600)
	if err != nil {
		log.Println(err)
		return false
	}

	defer f.Close()

	n, err := f.WriteString(wr.String())
	if err != nil {
		log.Printf("Unexpected Error When Writing To File: %v", err)
		return false
	}
	log.Printf("Wrote %d Bytes!", n)
	return true

}
