package templating

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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
	ServiceName      string
	ResourceName     string
	ResourceType     string
	Resource         interface{}
	ResourceJSON     string
	//resource     TemplateResource *need to create a template resource that works with all kubernetes resources*
}

// CacheTemplating takes a resource cache and renders its respective template values
func CacheTemplating(rcache *resourcecache.ResourceCache, serviceName string) map[string]string {
	var t map[string]string
	t = make(map[string]string)
	c := rcache.GetCache()
	for kt, r := range *c {

		// TODO need to remove the hardcoded first resourcefunction

		bytes, err := json.Marshal(r.GetResourceFunctions()[0].Data)
		if err != nil {
			_ = fmt.Errorf("%v", err)
		}
		conf := TemplateConfig{
			ImportStatements: getTemplateImports(&r.GetResourceFunctions()[0].Data),
			PackageName:      string(r.PackageName),
			ServiceName:      serviceName,
			ResourceName:     getTemplateResourceTitle(&r.GetResourceFunctions()[0].Data),
			ResourceType:     reflect.TypeOf(&r.GetResourceFunctions()[0].Data).String(),
			Resource:         getTemplateResourceTitle(&r.GetResourceFunctions()[0].Data),
			ResourceJSON:     string(bytes),
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

//FOLLOWING FUNCTION (TestDeclarationTemplate) USED FOR TESTING ONLY

// TestDeclarationTemplate used to debug template output
func TestDeclarationTemplate(resource interface{}) error {
	tpl, err := funcDeclarationTemplate(resource)
	fmt.Println(tpl)
	return err
}

// CreateConfig creates a configuration based on the resource
func CreateConfig(resource interface{}) (TemplateConfig, error) {
	bytes, err := json.Marshal(resource)
	if err != nil {
		return TemplateConfig{}, err
	}

	c := TemplateConfig{
		ImportStatements: getTemplateImports(&resource),
		ServiceName:      "",
		Resource:         resource,
		PackageName:      getTemplatePackageName(&resource),
		ResourceType:     reflect.TypeOf(resource).String(),
		ResourceName:     getTemplateResourceTitle(&resource),
		ResourceJSON:     string(bytes),
	}

	return c, nil
}

func funcDeclarationTemplate(resource interface{}) (string, error) {

	cwd, err := os.Getwd()
	temp, err := template.New("resourceFunc.tmpl").ParseFiles(filepath.Join(cwd, "pkg", "templating", "resourceFunc.tmpl"))
	if err != nil {
		return "", err
	}

	c, err := CreateConfig(resource)
	if err != nil {
		return "", err
	}
	c.ServiceName = "*v1alpha1.Collectd"

	var wr bytes.Buffer
	err = temp.Execute(&wr, c)
	if err != nil {
		return "", err
	}

	return wr.String(), nil

}

func configDeclarationTemplate(c TemplateConfig) (string, error) {

	cwd, err := os.Getwd()
	temp, err := template.New("resourceFunc.tmpl").ParseFiles(filepath.Join(cwd, "pkg", "templating", "resourceFunc.tmpl"))
	if err != nil {
		return "", err
	}

	c.ServiceName = "*v1alpha1.Collectd"

	var wr bytes.Buffer
	err = temp.Execute(&wr, c)
	if err != nil {
		return "", err
	}

	return wr.String(), nil
}
