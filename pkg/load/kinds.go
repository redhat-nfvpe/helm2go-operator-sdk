package load

import (
	"regexp"

	"github.com/redhat-nfvpe/helm2go-operator-sdk/internal/resourcecache"
)

// AcceptedK8sTypes contains the supported core Kubernetes types
var AcceptedK8sTypes = regexp.MustCompile(`(Role|ClusterRole|RoleBinding|ClusterRoleBinding|ServiceAccount|Service|Deployment)`)

type resourceConfig struct {
	r  interface{}
	kt resourcecache.KindType
	pt resourcecache.PackageType
}

func (rc *resourceConfig) toResourceCache() *resourcecache.Resource {
	var r resourcecache.Resource

	r.PackageName = rc.pt
	r.FileName = string(rc.kt)

	// put things into resourcecache format

	return &r
}

// func getKindType(kind string) *resourcecache.KindType {
// 	var kt resourcecache.KindType
// 	switch kind {
// 	case resourcecache.KindTypeConfigMap.(string):
// 		kt =
// 	}
// }
