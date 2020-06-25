package cache_service

import (
	"github.com/bmsandoval/vehicle-api/pkg/appcontext"
	_ "github.com/lib/pq"
	gocache "github.com/pmylund/go-cache"
	"time"
)

type ServiceImplementation struct {
	AppCtx appcontext.Context
	Cache *gocache.Cache
	Timeout  time.Duration
}

type CacheService interface {
	Set(string, interface{}) error
	Get(string) ([]byte, error)
	Delete(string)
}
