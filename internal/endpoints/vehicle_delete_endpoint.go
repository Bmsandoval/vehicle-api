package endpoints

import (
	"context"
	"github.com/bmsandoval/vehicle-api/internal/interchange"
	"github.com/bmsandoval/vehicle-api/pkg/appcontext"
	"github.com/bmsandoval/vehicle-api/pkg/services"
	"github.com/go-kit/kit/endpoint"
	//"github.com/pquerna/cachecontrol"
)

func MakeVehicleDeleteEndpoint(appCtx appcontext.Context, servicesBundle services.Bundle) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		requestObj := request.(interchange.VehicleDeleteRequest)

		// Delete an entry and collect count
		count, err := servicesBundle.VehicleSvc.Delete(requestObj.Id)
		if err != nil {
			return nil, err
		}

		// Delete entry from cache
		servicesBundle.CacheSvc.Delete(string(requestObj.Id))

		return interchange.VehicleDeleteResponse{
			StatusCode: 200,
			CountDeleted: count,
		}, nil
	}
}
