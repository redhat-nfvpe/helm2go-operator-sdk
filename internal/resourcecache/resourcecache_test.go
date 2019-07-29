package resourcecache

import (
	"fmt"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const defaultkindtype = KindTypeConfigMap

func TestResourceCache(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Resource Cache")
}

var _ = Describe("Resource Cache", func() {
	var resourceCache *ResourceCache
	It("Has Zero Size When Empty", func() {
		resourceCache = NewResourceCache()
		Expect(resourceCache.Size()).To(Equal(0))
	})
	It("Sets Correct Filename", func() {
		filename := fmt.Sprintf("%s.go", strings.ToLower(string(defaultkindtype)))
		resourceCache.SetResourceForKindType(defaultkindtype, PackageTypeConfigMaps)
		Expect(resourceCache.GetResourceForKindType(defaultkindtype).FileName).To(Equal(filename))
	})
	It("Sets Correct Type", func() {
		resourceCache := GetResourceCacheWithKindOne()
		Expect(resourceCache.Size()).To(Equal(1))

		resourceCache.SetResourceForKindType(defaultkindtype, PackageTypeConfigMaps)
		configmapResource := resourceCache.GetResourceForKindType(defaultkindtype)
		configmapResource.SetResourceFunctions("NewConfigMapForCR", getConfigMap())
		_, ok := configmapResource.GetResourceFunctions()[0].Data.(*v1.ConfigMap)
		Expect(ok).To(BeTrue())
	})
	It("Fails On Wrong Type", func() {
		resourceCache := GetResourceCacheWithKindOne()
		size := resourceCache.Size()
		Expect(size).To(Equal(1))

		resourceCache.SetResourceForKindType(defaultkindtype, PackageTypeConfigMaps)
		configmapResource := resourceCache.GetResourceForKindType(defaultkindtype)
		configmapResource.SetResourceFunctions("NewConfigMapForCR", getServiceAccount())

		_, ok := configmapResource.GetResourceFunctions()[0].Data.(*v1.ConfigMap)
		Expect(ok).To(BeFalse())
	})
})

var _ = Describe("Filetype Extensions", func() {
	It("Creates Go File Extension", func() {
		typeString := "Role"
		var n string
		n = nameToFileName(typeString, FileExtensionGo)
		Expect(n).To(Equal("role.go"))

		typeString = "RoleBinding"
		n = nameToFileName(typeString, FileExtensionGo)
		Expect(n).To(Equal("roleBinding.go"))
	})
})

func GetResourceCacheWithKindOne() *ResourceCache {
	resourceCache := NewResourceCache()
	resourceCache.SetKindType(KindTypeConfigMap)
	return resourceCache

}

func getConfigMap() *v1.ConfigMap {

	configmap := &v1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "ConfigMap",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "CRD_NAME",
			Namespace: "CRD_NAMESPACE",
		},
		Data: map[string]string{
			"config": "config",
		},
	}
	return configmap
}

func getServiceAccount() *v1.ServiceAccount {
	serviceaccount := &v1.ServiceAccount{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ServiceAccount",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "CRD_NAME",
			Namespace: "CRD_NAMESPACE",
		},
	}
	return serviceaccount
}
