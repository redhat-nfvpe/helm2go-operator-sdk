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
	"strings"
	"text/template"

	"github.com/redhat-nfvpe/helm2go-operator-sdk/internal/resourcecache"
)

// TemplateConfig is used as payload for executing the templating functionality
type TemplateConfig struct {
	ImportStatements          []string
	PackageName               string
	Kind                      string
	APIVersion                string
	OwnerAPIVersion           string
	ResourceName              string
	ResourceType              string
	ResourceTypeInstantiation string
	Resource                  interface{}
	ResourceJSON              string
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
	Kind                 string
	LowerKind            string
	OwnerAPIVersion      string
	ImportMap            map[string]string
	ResourceWatches      []string
	ResourceForReconcile []string
}

// ReconcileTemplateConfig ...
type ReconcileTemplateConfig struct {
	ResourceImportPackage   string
	ResourceType            string
	ResourceName            string
	LowerResourceName       string
	Kind                    string
	LowerKind               string
	OwnerAPIVersion         string
	PluralLowerResourceName string
}

// CacheTemplating takes a resource cache and renders its respective template values
func CacheTemplating(rcache *resourcecache.ResourceCache, outputDir, kind, apiVersion string) map[string]string {
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
			ImportStatements:          append(getTemplateImports(&r.GetResourceFunctions()[0].Data), getOwnerAPIVersion(apiVersion, kind)+` "`+getAppTypeImport(outputDir, apiVersion)+`"`),
			PackageName:               string(r.PackageName),
			Kind:                      kind,
			ResourceName:              getTemplateResourceTitle(&r.GetResourceFunctions()[0].Data),
			ResourceType:              reflect.TypeOf((r.GetResourceFunctions()[0].Data)).String(),
			ResourceTypeInstantiation: getInstantiationString(reflect.TypeOf((r.GetResourceFunctions()[0].Data)).String()),
			Resource:                  getTemplateResourceTitle(&r.GetResourceFunctions()[0].Data),
			ResourceJSON:              string(bytes),
			APIVersion:                apiVersion,
			OwnerAPIVersion:           getOwnerAPIVersion(apiVersion, kind),
		}
		tmpl, err := configDeclarationTemplate(conf)
		t[string(kt)] = tmpl
	}

	return t
}

func getInstantiationString(resourceType string) string {
	base := resourceType[1:]
	return "&" + base + "{}"
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

	temp, err := template.New("resourceFuncTemplate").Parse(getResourceTemplate())
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
	var reconcileBlocks []string
	var resourceControllerImports map[string]string
	var wr bytes.Buffer
	var temp *template.Template
	var err error
	var c ControllerTemplateConfig
	var w ControllerWatchFuncConfig

	resourceControllerImports = make(map[string]string)
	ownerAPIVersion := getOwnerAPIVersion(apiVersion, kind)
	lowerKind := kindToLowerCamel(kind)

	for _, r := range rcache.PrepareCacheForFile() {

		resourceControllerImports[getResourceCRImport(outputDir, strings.ToLower(getTemplateResourceTitle(&r.GetResourceFunctions()[0].Data)))] = ""

		wr = bytes.Buffer{}
		temp, err = template.New("controllerFuncTemplate").Parse(getControllerFunctionTemplate())
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

	cacheToLookup := *rcache.GetCache()

	// for resourcetype in lookup
	for _, kindtype := range resourcecache.KindTypeLookup {
		if r, ok := cacheToLookup[kindtype]; ok {
			wr = bytes.Buffer{}
			temp, err = template.New("reconcileResourceTemplate").Parse(getReconcileBlockTemplate())
			if err != nil {
				log.Println(err)
				return false
			}

			resourceType := getTemplateResourceTitle(&r.GetResourceFunctions()[0].Data)

			reconcileConfig := ReconcileTemplateConfig{
				OwnerAPIVersion:         ownerAPIVersion,
				Kind:                    kind,
				LowerKind:               lowerKind,
				ResourceImportPackage:   getResourceImportPackage(resourceType),
				ResourceType:            resourceType,
				ResourceName:            getTemplateResourceTitle(&r.GetResourceFunctions()[0].Data),
				LowerResourceName:       strings.ToLower(getTemplateResourceTitle(&r.GetResourceFunctions()[0].Data)),
				PluralLowerResourceName: getNamePlural(strings.ToLower(getTemplateResourceTitle(&r.GetResourceFunctions()[0].Data))),
			}

			err = temp.Execute(&wr, reconcileConfig)
			if err != nil {
				log.Println(err)
				return false
			}

			reconcileBlocks = append(reconcileBlocks, wr.String())
		}
	}
	wr = bytes.Buffer{}
	temp, err = template.New("controllerTemplate").Parse(getControllerTemplate())
	importMap := getImportMap(outputDir, kind, apiVersion)
	for k, v := range resourceControllerImports {
		importMap[k] = v
	}
	c = ControllerTemplateConfig{
		kind,
		lowerKind,
		ownerAPIVersion,
		importMap,
		watchFuncs,
		reconcileBlocks,
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

func getReconcileBlockTemplate() string {
	return `
	{{.LowerResourceName}} := {{.PluralLowerResourceName}}.New{{.ResourceName}}ForCR(instance)
	// Set {{ .Kind }} instance as the owner and controller
	if err{{.ResourceName}} := controllerutil.SetControllerReference(instance, {{.LowerResourceName}}, r.scheme); err{{.ResourceName}} != nil {
		return reconcile.Result{}, err{{.ResourceName}}
	}
	// Check if this {{.ResourceName}} already exists
	found{{.ResourceName}} := &{{.ResourceImportPackage}}.{{.ResourceType}}{}
	err{{.ResourceName}} := r.client.Get(context.TODO(), types.NamespacedName{Name: {{.LowerResourceName}}.Name, Namespace: {{.LowerResourceName}}.Namespace}, found{{.ResourceName}})
	if err{{.ResourceName}} != nil && errors.IsNotFound(err{{.ResourceName}}) {
		reqLogger.Info("Creating a new {{.ResourceName}}", "{{.LowerResourceName}}.Namespace", {{.LowerResourceName}}.Namespace, "{{.LowerResourceName}}.Name", {{.LowerResourceName}}.Name)
		err{{.ResourceName}} = r.client.Create(context.TODO(), found{{.ResourceName}})
		if err{{.ResourceName}} != nil {
			return reconcile.Result{}, err{{.ResourceName}}
		}
		// Pod created successfully - don't requeue
		return reconcile.Result{}, nil
	} else if err{{.ResourceName}} != nil {
		return reconcile.Result{}, err{{.ResourceName}}
	}
	`
}
