package load

import (
	"log"
	"regexp"

	"github.com/redhat-nfvpe/helm2go-operator-sdk/internal/validatemap"
)

func isEmptyFile(f []byte) bool {
	if len(f) == 0 {
		return true
	}
	isAlpha := regexp.MustCompile(`[A-Za-z]+`).MatchString
	if !isAlpha(string(f)) {
		return true
	}
	return false
}

// DeprecatedResourceAPIVersionError is used to alert that a resource api version is deprecated
type DeprecatedResourceAPIVersionError struct {
	message      string
	resourceKind string
	apiVersion   string
}

func (d *DeprecatedResourceAPIVersionError) Error() string {
	return d.message
}

func containValidMap(fileName string, validMap *validatemap.ValidateMap) bool {
	// handle empty map lookups
	if len(validMap.Map) == 0 {
		return false
	}
	if _, ok := validMap.Map[fileName]; ok {
		return true
	}
	return false
}

func logContinue() {
	log.Println("Continuing")
}
