package load

import (
	"fmt"
	"regexp"
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
	resourceKind string
	apiVersion   string
}

func (d *DeprecatedResourceAPIVersionError) Error() string {
	return fmt.Sprintf("Resource of Kind: %s has deprecated API Version: %s\n", d.resourceKind, d.apiVersion)
}
