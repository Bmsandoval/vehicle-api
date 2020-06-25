package services

import (
	"database/sql"
	"github.com/bmsandoval/vehicle-api/pkg/appcontext"
	"github.com/bmsandoval/vehicle-api/pkg/services/cache_service"
	"github.com/bmsandoval/vehicle-api/pkg/services/vehicle_service"
	gocache "github.com/pmylund/go-cache"
)

type Bundle struct {
	VehicleSvc vehicle_service.VehicleService
	CacheSvc cache_service.CacheService
}

func NewBundle(appCtx appcontext.Context, db *sql.DB, cache *gocache.Cache) (*Bundle, error) {
	bundle := &Bundle{}

	bundle.VehicleSvc = vehicle_service.ServiceImplementation{
		AppCtx: appCtx,
		DB: db,
	}
	bundle.CacheSvc = cache_service.ServiceImplementation{
		AppCtx: appCtx,
		Cache: cache,
	}

	return bundle, nil
}
