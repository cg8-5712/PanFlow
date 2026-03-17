package cache

import (
	"context"
	"time"

	"github.com/maypok86/otter"
)

var otterCache *otter.CacheWithVariableTTL[string, any]

// InitOtter initializes the L1 memory cache (Otter) with variable TTL support
func InitOtter(maxCapacity int) error {
	c, err := otter.MustBuilder[string, any](maxCapacity).
		CollectStats().
		WithVariableTTL().
		Build()
	if err != nil {
		return err
	}
	otterCache = &c
	return nil
}

// OtterGet retrieves a value from L1 cache
func OtterGet(key string) (any, bool) {
	if otterCache == nil {
		return nil, false
	}
	return otterCache.Get(key)
}

// OtterSet stores a value in L1 cache with TTL
func OtterSet(key string, value any, ttl time.Duration) bool {
	if otterCache == nil {
		return false
	}
	return otterCache.Set(key, value, ttl)
}

// OtterDelete removes a value from L1 cache
func OtterDelete(key string) {
	if otterCache != nil {
		otterCache.Delete(key)
	}
}

// OtterClear clears all entries from L1 cache
func OtterClear() {
	if otterCache != nil {
		otterCache.Clear()
	}
}

// OtterClose closes the L1 cache
func OtterClose() {
	if otterCache != nil {
		otterCache.Close()
	}
}

// OtterSize returns the number of entries in L1 cache
func OtterSize() int {
	if otterCache != nil {
		return otterCache.Size()
	}
	return 0
}

// OtterHas checks if a key exists in L1 cache
func OtterHas(key string) bool {
	if otterCache != nil {
		return otterCache.Has(key)
	}
	return false
}

// OtterDeleteByFunc deletes entries matching a predicate
func OtterDeleteByFunc(_ context.Context, fn func(key string, value any) bool) {
	if otterCache != nil {
		otterCache.DeleteByFunc(fn)
	}
}
