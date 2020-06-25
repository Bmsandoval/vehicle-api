package endpoints

import (
	"bytes"
	"context"
	"encoding/gob"
	"github.com/bmsandoval/vehicle-api/internal/interchange"
	"github.com/bmsandoval/vehicle-api/pkg/appcontext"
	"github.com/bmsandoval/vehicle-api/pkg/models"
	"github.com/bmsandoval/vehicle-api/pkg/services"
	"github.com/go-kit/kit/endpoint"
)

func MakeVehicleReadEndpoint(appCtx appcontext.Context, servicesBundle services.Bundle) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		requestObj := request.(interchange.VehicleReadRequest)

		// Check cache for the result we want
		cachedResult, err := servicesBundle.CacheSvc.Get(string(requestObj.Id))
		if err != nil {
			return nil, err
		}

		if cachedResult != nil {
			// Decode result if cache returned something
			var resolvedCacheResult models.Vehicle

			buffer := bytes.NewBuffer(cachedResult)
			dec := gob.NewDecoder(buffer)
			if err := dec.Decode(&resolvedCacheResult); err != nil {
				return nil, err
			}

			return interchange.VehicleReadResponse{
				StatusCode: 200,
				Vehicle:    resolvedCacheResult,
			}, nil
		}

		// Read from DB. We won't get here if the cache returned something
		result, err := servicesBundle.VehicleSvc.Read(requestObj.Id)
		if err != nil {
			return nil, err
		}

		var output models.Vehicle
		var statusCode int
		if result != nil {
			// return OK if entry found
			statusCode = 200
			output = *result

			// If we got this far than it was a cache miss
			// update cache with the entry that was missing
			err := servicesBundle.CacheSvc.Set(string(result.Id), result)
			if err != nil {
				return nil, err
			}
		} else {
			// return NO CONTENT if entry not found
			statusCode = 204
		}

		return interchange.VehicleReadResponse{
			StatusCode: statusCode,
			Vehicle: output,
		}, nil
	}
}
