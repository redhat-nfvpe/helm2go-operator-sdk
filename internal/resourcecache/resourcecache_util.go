package resourcecache

import (
	"github.com/iancoleman/strcase"
)

// FileExtension Used for "type safe file extension use"
type FileExtension string

var (
	// FileExtensionGo specifies the golang file extension
	FileExtensionGo FileExtension = ".go"
	// FileExtensionYAML specifies the yaml file extension
	FileExtensionYAML = ".yaml"
)

// KindTypeLookup is used to specify an order in resources
var KindTypeLookup = map[int]KindType{
	0:  KindTypeConfigMap,
	1:  KindTypeServiceAccount,
	2:  KindTypeRole,
	3:  KindTypeClusterRole,
	4:  KindTypeRoleBinding,
	5:  KindTypeClusterRoleBinding,
	6:  KindTypeService,
	7:  KindTypeSecret,
	8:  KindTypeVolume,
	9:  KindTypeDaemonSet,
	10: KindTypeDeployment,
}

func (f *FileExtension) String() string {
	return string(*f)
}

func nameToFileName(name string, extension FileExtension) string {
	// Takes a name and returns the appropriate filename
	return strcase.ToLowerCamel(name) + extension.String()
}
