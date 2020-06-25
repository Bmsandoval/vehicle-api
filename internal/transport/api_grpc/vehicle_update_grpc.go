package api_grpc

import (
	"context"
	"github.com/bmsandoval/vehicle-api/internal/endpoints"
	"github.com/bmsandoval/vehicle-api/internal/interchange"
	"github.com/bmsandoval/vehicle-api/pkg/models"
	pb "github.com/bmsandoval/vehicle-api/protos"
)

// GRPC Update endpoint wrapper
func (s GrpcServer) VehicleUpdate(ctx context.Context, in *pb.VehicleUpdateRequest) (*pb.VehicleUpdateResponse, error) {
	endpoint := endpoints.MakeVehicleUpdateEndpoint(s.Context, s.ServicesBundle)
	requestObj := interchange.VehicleUpdateRequest{
			Vehicle: models.Vehicle{
				Id:    in.Id,
				Make:  in.Make,
				Model: in.Model,
				Vin:   in.Vin,
			},
		}
	response, err := endpoint(s.Context.GoContext, requestObj)
	if err != nil {
		return nil, err
	}

	resolvedResponse := response.(interchange.VehicleUpdateResponse)

	return &pb.VehicleUpdateResponse{
		CountUpdated: resolvedResponse.CountUpdated,
	}, err
}
