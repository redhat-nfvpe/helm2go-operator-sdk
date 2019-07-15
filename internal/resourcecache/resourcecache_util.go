package resourcecache

import (
	"github.com/iancoleman/strcase"
)

// FileExtension Used for "type safe file extension use"
type FileExtension string

var (
	FileExtensionGo   FileExtension = ".go"
	FileExtensionYAML               = ".yaml"
)

func (f *FileExtension) String() string {
	return string(*f)
}

func nameToFileName(name string, extension FileExtension) string {
	// Takes a name and returns the appropriate filename
	return strcase.ToLowerCamel(name) + extension.String()
}
