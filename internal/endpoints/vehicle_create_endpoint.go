package endpoints

import (
	"context"
	"github.com/bmsandoval/vehicle-api/internal/interchange"
	"github.com/bmsandoval/vehicle-api/pkg/appcontext"
	"github.com/bmsandoval/vehicle-api/pkg/models"
	"github.com/bmsandoval/vehicle-api/pkg/services"
	"github.com/go-kit/kit/endpoint"
)

func MakeVehicleCreateEndpoint(appCtx appcontext.Context, servicesBundle services.Bundle) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		requestObj := request.(interchange.VehicleCreateRequest)

		// Create a new vehicle
		result, err := servicesBundle.VehicleSvc.Create(&models.Vehicle{
			Make:  requestObj.Make,
			Model: requestObj.Model,
			Vin:   requestObj.Vin,
		})
		if err != nil {
			return nil, err
		}

		var output models.Vehicle
		if result != nil {
			output = *result
			// If creation is successful, add entry to cache
			err := servicesBundle.CacheSvc.Set(string(output.Id), &output)
			if err != nil {
				return nil, err
			}
		}

		return interchange.VehicleCreateResponse{
			StatusCode: 200,
			Vehicle: output,
		}, nil
	}
}
