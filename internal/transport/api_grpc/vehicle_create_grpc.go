package api_grpc

import (
	"context"
	"github.com/bmsandoval/vehicle-api/internal/endpoints"
	"github.com/bmsandoval/vehicle-api/internal/interchange"
	pb "github.com/bmsandoval/vehicle-api/protos"
)

// GRPC Create endpoint wrapper
func (s GrpcServer) VehicleCreate(ctx context.Context, in *pb.VehicleCreateRequest) (*pb.VehicleCreateResponse, error) {
	endpoint := endpoints.MakeVehicleCreateEndpoint(s.Context, s.ServicesBundle)
	requestObj := interchange.VehicleCreateRequest{
			Make:  in.Make,
			Model: in.Model,
			Vin:   in.Vin,
		}
	response, err := endpoint(s.Context.GoContext, requestObj)
	if err != nil {
		return nil, err
	}

	if response == nil {
		return nil, nil
	}

	resolvedResponse := response.(interchange.VehicleCreateResponse)

	return &pb.VehicleCreateResponse{
		Id: resolvedResponse.Id,
		Make:  resolvedResponse.Make,
		Model: resolvedResponse.Model,
		Vin:   resolvedResponse.Vin,
	}, err
}
