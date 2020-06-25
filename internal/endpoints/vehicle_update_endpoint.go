package endpoints

import (
	"context"
	"github.com/bmsandoval/vehicle-api/internal/interchange"
	"github.com/bmsandoval/vehicle-api/pkg/appcontext"
	"github.com/bmsandoval/vehicle-api/pkg/models"
	"github.com/bmsandoval/vehicle-api/pkg/services"
	"github.com/go-kit/kit/endpoint"
	//"github.com/pquerna/cachecontrol"
)

func MakeVehicleUpdateEndpoint(appCtx appcontext.Context, servicesBundle services.Bundle) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		requestObj := request.(interchange.VehicleUpdateRequest)

		// Read in existing record from DB
		existingRecord, err := servicesBundle.VehicleSvc.Read(requestObj.Id)
		if err != nil {
			return nil, err
		}

		// return NO CONTENT if updating a non-existent entry
		if existingRecord == nil {
			return interchange.VehicleUpdateResponse {
				StatusCode: 204,
				CountUpdated: 0,
			}, nil
		}

		// Check if anything in the request is different from existing object
		if existingRecord.Id == requestObj.Id &&
			existingRecord.Make == requestObj.Make &&
			existingRecord.Model == requestObj.Model &&
			existingRecord.Vin == requestObj.Vin {

			// return NOT MODIFIED if existing entry matches updates
			return interchange.VehicleUpdateResponse{
				StatusCode: 204,
				CountUpdated: 0,
			}, nil
		}

		updateObj := &models.Vehicle{
			Id: requestObj.Id,
			Make:  requestObj.Make,
			Model: requestObj.Model,
			Vin:   requestObj.Vin,
		}

		// If there are differences between current record and new record, update
		count, err := servicesBundle.VehicleSvc.Update(updateObj)
		if err != nil {
			return nil, err
		}

		// If record changed, update cache
		err = servicesBundle.CacheSvc.Set(string(requestObj.Id), updateObj)
		if err != nil {
			return nil, err
		}

		// return OK and 1 to say successfully updated a single entry
		return interchange.VehicleUpdateResponse{
			StatusCode: 200,
			CountUpdated: count,
		}, nil
	}
}
