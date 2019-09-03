package load

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/redhat-nfvpe/helm2go-operator-sdk/internal/validatemap"
)

// PerformResourceValidation validates templated resources to identify deprecated API versions
func PerformResourceValidation(resourcesPath string) (*validatemap.ValidateMap, error) {
	// collect all templated files in directory
	// TODO: need to test with nested templates
	files, err := ioutil.ReadDir(resourcesPath)
	if err != nil {
		return nil, fmt.Errorf("error reading files in directory %s: %v", resourcesPath, err)
	}

	var validMap validatemap.ValidateMap

	for _, file := range files {
		rconfig, err := yamlUnmarshalSingleResource(filepath.Join(resourcesPath, file.Name()))
		if err != nil {
			reader := bufio.NewReader(os.Stdin)
			if e := err.Error(); e == "deprecated" {
				fmt.Printf("Resource: %v has a deprecated API version. Please enter 'Y' to proceed without this resource or 'N' to exit the program: ", reflect.TypeOf(rconfig.resource))
				c, err := reader.ReadByte()
				if err != nil {
					fmt.Println(err)
					cleanUpAndExit(resourcesPath)
				}
				if c == []byte("Y")[0] || c == []byte("y")[0] {
					addResourceToContinue(&validMap, file.Name())
				} else {
					cleanUpAndExit(resourcesPath)
				}
			} else if e == "unsupported" {

				fmt.Printf("Resource: %v is unsupported. Please enter 'Y' to proceed without this resource, or 'N' to exit the program: ", reflect.TypeOf(rconfig.resource))
				c, err := reader.ReadByte()
				if err != nil {
					fmt.Println(err)
					cleanUpAndExit(resourcesPath)
				}
				if c == []byte("Y")[0] || c == []byte("y")[0] {
					addResourceToContinue(&validMap, file.Name())
				} else {
					cleanUpAndExit(resourcesPath)
				}

			} else if e == "not yaml" || e == "empty" {
				continue
			} else {
				return nil, fmt.Errorf("uncaught error: %v", err)
			}
		}
	}
	return &validMap, nil
}

func matchString(text string, contains string) bool {
	return strings.Contains(text, contains)
}

func cleanUpAndExit(resourcesPath string) {
	// log exit begin
	log.Println("Cleanup Temp Started")
	// clean up temp folder
	temp := filepath.Dir(filepath.Dir(resourcesPath))
	err := os.RemoveAll(temp)
	if err != nil {
		log.Printf("Error While Cleaning Up Temp: %v", err)
	} else {
		// log cleanup
		log.Println("Cleanup Temp Completed")
	}
	// log exit
	log.Println("Will Exit")
	// exit
	os.Exit(1)
}

func addResourceToContinue(validMap *validatemap.ValidateMap, fileName string) {
	// creates the map in the case that it does not exist
	if len(validMap.Map) == 0 {
		validMap.Map = make(map[string]bool)
	}
	validMap.Map[fileName] = true
}
