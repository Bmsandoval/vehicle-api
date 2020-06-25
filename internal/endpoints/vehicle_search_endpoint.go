package endpoints

import (
	"context"
	"github.com/bmsandoval/vehicle-api/internal/interchange"
	"github.com/bmsandoval/vehicle-api/pkg/appcontext"
	"github.com/bmsandoval/vehicle-api/pkg/services"
	pb "github.com/bmsandoval/vehicle-api/protos"
)

func MakeVehicleSearchEndpoint(appCtx appcontext.Context, server pb.Vehicles_VehicleSearchServer, servicesBundle services.Bundle) func(context.Context, interface{})error {
	return func(ctx context.Context, request interface{}) error {
		requestObj := request.(interchange.VehicleSearchRequest)

		// Track groups as we pull from DB
		limit := 10
		offset := 0

		for ;; {
			// Search DB for the (NEXT) group of size LIMIT
			vehicles, err := servicesBundle.VehicleSvc.Search(limit, offset, requestObj.Make, requestObj.Model)
			if err != nil {
				return err
			}

			// Stream a response for each vehicle we retrieve
			for _, vehicle := range vehicles {
				err = server.Send(&pb.VehicleSearchResponse{
					Id:    vehicle.Id,
					Make:  vehicle.Make,
					Model: vehicle.Model,
					Vin:   vehicle.Vin,
				})
				if err != nil {
					return err
				}
			}

			// If the amount returned is less than the limit, than this is the last group
			if len(vehicles) < limit {
				break
			}
		}

		return nil
	}
}
