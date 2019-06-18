package resourcecache

import (
	"fmt"
	"strings"
	"testing"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const defaultkindtype = KindTypeConfigMap

func TestZeroResourceCache(t *testing.T) {
	resourceCache := NewResourceCache()
	if size := resourceCache.Size(); size != 0 {
		t.Errorf("wrong count of cache size , expected 0 and got %d", size)
	}
}
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
func TestOneResourceCache(t *testing.T) {
	resourceCache := GetResourceCacheWithKindOne()
	if size := resourceCache.Size(); size != 1 {
		t.Errorf("wrong count of cache size , expected 1 and got %d", size)
	}
}

func TestZeroResource(t *testing.T) {
	resourceCache := GetResourceCacheWithKindOne()
	if size := resourceCache.Size(); size != 1 {
		t.Errorf("wrong count of cache size , expected 0 and got %d", size)
	}
	filename := fmt.Sprintf("%s.go", strings.ToLower(string(defaultkindtype)))
	resourceCache.SetResourceForKindType(defaultkindtype, PackageTypeConfigMaps)

	if resourceCache.GetResourceForKindType(defaultkindtype).FileName != filename {
		t.Errorf("wrong filename  , expected %s and got %s", filename, resourceCache.GetResourceForKindType(defaultkindtype).FileName)
	}
}

func TestDataInterfaceForCorrectType(t *testing.T) {
	resourceCache := GetResourceCacheWithKindOne()
	if size := resourceCache.Size(); size != 1 {
		t.Errorf("wrong count of cache size , expected 0 and got %d", size)
	}
	resourceCache.SetResourceForKindType(defaultkindtype, PackageTypeConfigMaps)
	configmapResource := resourceCache.GetResourceForKindType(defaultkindtype)

	configmapResource.SetResourceFunctions("NewConfigMapForCR", getConfigMap())

	if r, ok := configmapResource.GetResourceFunctions()[0].Data.(*v1.ConfigMap); !ok {
		t.Errorf("wrong data type  , expected %T and got %T", configmapResource.GetResourceFunctions()[0].Data, r)
	}
}

func TestDataInterfaceForWrongType(t *testing.T) {
	resourceCache := GetResourceCacheWithKindOne()
	if size := resourceCache.Size(); size != 1 {
		t.Errorf("wrong count of cache size , expected 0 and got %d", size)
	}
	resourceCache.SetResourceForKindType(defaultkindtype, PackageTypeConfigMaps)
	configmapResource := resourceCache.GetResourceForKindType(defaultkindtype)
	configmapResource.SetResourceFunctions("NewConfigMapForCR", getServiceAccount())

	if r, ok := configmapResource.GetResourceFunctions()[0].Data.(*v1.ConfigMap); ok {
		t.Errorf("wrong data type  , expected %T and got %T", configmapResource.GetResourceFunctions()[0].Data, r)
	}
}
