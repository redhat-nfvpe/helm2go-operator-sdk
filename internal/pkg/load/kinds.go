package load

import (
	"regexp"

	"github.com/redhat-nfvpe/helm2go-operator-sdk/internal/resourcecache"
)

// AcceptedK8sTypes contains the supported core Kubernetes types
var AcceptedK8sTypes = regexp.MustCompile(`(Role|ClusterRole|RoleBinding|ClusterRoleBinding|ServiceAccount|Service|Deployment)`)

type resourceConfig struct {
	resource    interface{}
	kindType    resourcecache.KindType
	packageType resourcecache.PackageType
}

func (rc *resourceConfig) toResourceCache() *resourcecache.Resource {
	var resource resourcecache.Resource

	resource.PackageName = rc.packageType
	resource.FileName = string(rc.kindType)

	return &resource
}
