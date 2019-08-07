package templating

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/redhat-nfvpe/helm2go-operator-sdk/internal/resourcecache"
)

// CacheTemplating takes a resource cache and renders its respective template values
func CacheTemplating(rcache *resourcecache.ResourceCache, outputDir, kind, apiVersion string) map[string]string {
	var templates map[string]string
	templates = make(map[string]string)
	cache := rcache.PrepareCacheForFile()

	for kindType, resource := range cache {
		// TODO need to remove the hardcoded first resourcefunction
		kindType = filepath.Join(filepath.Dir(kindType), resource.PackageName.String(), filepath.Base(kindType))
		bytes, err := json.MarshalIndent(resource.GetResourceFunctions()[0].Data, "", "\t")
		if err != nil {
			_ = fmt.Errorf("%v", err)
		}
		conf := NewResourceTemplateConfig(outputDir, apiVersion, kind, resource, bytes)
		tmpl, err := conf.Execute()
		templates[string(kindType)] = tmpl
	}

	return templates
}

// TemplatesToFiles takes the templated files and writes them to strings in the specified directory
func TemplatesToFiles(templates map[string]string, outputDir string) error {
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		err = os.MkdirAll(outputDir, 0700)
		if err != nil {
			return err
		}
	}

	for filename, file := range templates {
		// create the file to be written to
		f, err := os.Create(filepath.Join(outputDir, filename))
		if err != nil {
			return err
		}
		// write the file content to the actual file
		_, err = f.WriteString(file)
		if err != nil {
			return err
		}
		f.Close()
	}
	return nil
}

// ResourceFileStructure creates the correct file structure based on the spcified kind types
func ResourceFileStructure(rcache *resourcecache.ResourceCache, outputDir string) error {
	// iterate through the kind types
	for _, r := range *rcache.GetCache() {
		// create the neccessary folder for the kind type
		newOutput := filepath.Join(outputDir, r.PackageName.String())
		if _, err := os.Stat(newOutput); os.IsNotExist(err) {
			err = os.MkdirAll(newOutput, 0700)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// OverwriteController takes in a templated file and writes it over the original controller file
func OverwriteController(outputDir, kind, apiVersion string, rcache *resourcecache.ResourceCache) error {

	var watchFuncs []string
	var reconcileBlocks []string
	var resourceControllerImports map[string]string
	var err error

	resourceControllerImports = make(map[string]string)
	ownerAPIVersion := getOwnerAPIVersion(apiVersion, kind)
	lowerKind := kindToLowerCamel(kind)

	for _, resource := range rcache.PrepareCacheForFile() {

		resourceControllerImports[getResourceCRImport(outputDir, strings.ToLower(getTemplateResourceTitle(&resource.GetResourceFunctions()[0].Data)))] = ""
		watchFunctionConfig := NewControllerWatchFuncConfig(apiVersion, ownerAPIVersion, kind, lowerKind, resource)

		tmpl, err := watchFunctionConfig.Execute()
		if err != nil {
			return err
		}

		watchFuncs = append(watchFuncs, tmpl)
	}

	cacheToLookup := *rcache.GetCache()

	// for resourcetype in lookup
	for _, kindtype := range resourcecache.KindTypeLookup {
		// get resource of kindType
		if r, ok := cacheToLookup[kindtype]; ok {
			// create reconcile block template
			reconcileConfig := NewReconcileTemplateConfig(kind, lowerKind, r)
			temp, err := reconcileConfig.Execute()
			if err != nil {
				return err
			}
			reconcileBlocks = append(reconcileBlocks, temp)
		}
	}

	// get map of appropriate imports for controller file
	importMap := getImportMap(outputDir, kind, apiVersion)
	for k, v := range resourceControllerImports {
		importMap[k] = v
	}
	controllerConfig := NewControllerTemplateConfig(kind, lowerKind, ownerAPIVersion, importMap, watchFuncs, reconcileBlocks)

	tmpl, err := controllerConfig.Execute()

	// overwrite original file
	outFile := filepath.Join(outputDir, "pkg", "controller", kindToLower(kind), fmt.Sprintf("%s_controller.go", kindToLower(kind)))

	var f *os.File
	if f, err = os.OpenFile(outFile, os.O_WRONLY, 0600); err != nil {
		return err
	}
	// delete original content
	f.Truncate(0)

	defer f.Close()

	numBytes, err := f.WriteString(tmpl)
	if err != nil {
		return err
	}
	log.Printf("Wrote %d Bytes!", numBytes)
	return nil

}
