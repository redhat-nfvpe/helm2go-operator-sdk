package resourcecache

import (
	"fmt"
	"strings"
	"testing"
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
