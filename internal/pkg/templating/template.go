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
		conf := NewResourceTemplateConfig(outputDir, apiVersion, kind, r, bytes)
		tmpl, err := conf.Execute()
		t[string(kt)] = tmpl
	}

	return t
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
	var err error

	resourceControllerImports = make(map[string]string)
	ownerAPIVersion := getOwnerAPIVersion(apiVersion, kind)
	lowerKind := kindToLowerCamel(kind)

	for _, r := range rcache.PrepareCacheForFile() {

		resourceControllerImports[getResourceCRImport(outputDir, strings.ToLower(getTemplateResourceTitle(&r.GetResourceFunctions()[0].Data)))] = ""
		w := NewControllerWatchFuncConfig(apiVersion, ownerAPIVersion, kind, lowerKind, r)

		tmpl, err := w.Execute()
		if err != nil {
			log.Println(err)
			return false
		}

		watchFuncs = append(watchFuncs, tmpl)
	}

	cacheToLookup := *rcache.GetCache()

	// for resourcetype in lookup
	for _, kindtype := range resourcecache.KindTypeLookup {
		if r, ok := cacheToLookup[kindtype]; ok {
			reconcileConfig := NewReconcileTemplateConfig(kind, lowerKind, r)
			temp, err := reconcileConfig.Execute()
			if err != nil {
				log.Println(err)
				return false
			}
			reconcileBlocks = append(reconcileBlocks, temp)
		}
	}

	importMap := getImportMap(outputDir, kind, apiVersion)
	for k, v := range resourceControllerImports {
		importMap[k] = v
	}
	c := NewControllerTemplateConfig(kind, lowerKind, ownerAPIVersion, importMap, watchFuncs, reconcileBlocks)

	tmpl, err := c.Execute()

	// overwrite original file
	outFile := filepath.Join(outputDir, "pkg", "controller", kindToLower(kind), fmt.Sprintf("%s_controller.go", kindToLower(kind)))

	f, err := os.OpenFile(outFile, os.O_WRONLY, 0600)
	if err != nil {
		log.Println(err)
		return false
	}
	// delete original content
	f.Truncate(0)

	defer f.Close()

	n, err := f.WriteString(tmpl)
	if err != nil {
		log.Printf("Unexpected Error When Writing To File: %v", err)
		return false
	}
	log.Printf("Wrote %d Bytes!", n)
	return true

}
