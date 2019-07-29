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

var cd string

// PerformResourceValidation validates templated resources to identify deprecated API versions
func PerformResourceValidation(rp string) (*validatemap.ValidateMap, error) {
	cd = rp
	// collect all templated files in directory
	files, err := ioutil.ReadDir(rp)
	if err != nil {
		return nil, fmt.Errorf("error reading files in directory %s: %v", rp, err)
	}

	var validMap validatemap.ValidateMap

	for _, f := range files {
		rconfig, err := yamlUnmarshalSingleResource(filepath.Join(rp, f.Name()))
		if err != nil {
			reader := bufio.NewReader(os.Stdin)
			if e := err.Error(); e == "deprecated" {
				fmt.Printf("Resource: %v has a deprecated API version. Please enter 'continue' to proceed without this resource or 'stop' to exit the program: ", reflect.TypeOf(rconfig.r))
				text, _ := reader.ReadString('\n')
				if isStop(text) {
					cleanUpAndExit()
				}
				if isContinue(text) {
					addResourceToContinue(&validMap, f.Name())
				}
			} else if e == "unsupported" {

				fmt.Printf("Resource: %v is unsupported. Please enter 'continue' to proceed without this resource, or 'stop' to exit the program: ", reflect.TypeOf(rconfig.r))
				text, _ := reader.ReadString('\n')

				if isStop(text) {
					cleanUpAndExit()
				}
				if isContinue(text) {
					addResourceToContinue(&validMap, f.Name())
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

func isStop(text string) bool {
	return matchString(text, "stop")
}

func isContinue(text string) bool {
	return matchString(strings.ToLower(text), "continue")
}

func matchString(text string, contains string) bool {
	return strings.Contains(text, contains)
}

func cleanUpAndExit() {
	// log exit begin
	log.Println("Cleanup Temp Started")
	// clean up temp folder
	temp := filepath.Dir(filepath.Dir(cd))
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
