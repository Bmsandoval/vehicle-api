package api_grpc

import (
	"context"
	"github.com/bmsandoval/vehicle-api/internal/endpoints"
	"github.com/bmsandoval/vehicle-api/internal/interchange"
	pb "github.com/bmsandoval/vehicle-api/protos"
)

// GRPC Delete endpoint wrapper
func (s GrpcServer) VehicleDelete(ctx context.Context, in *pb.VehicleDeleteRequest) (*pb.VehicleDeleteResponse, error) {
	endpoint := endpoints.MakeVehicleDeleteEndpoint(s.Context, s.ServicesBundle)
	requestObj := interchange.VehicleDeleteRequest{
			Id: in.Id,
		}
	response, err := endpoint(s.Context.GoContext, requestObj)
	if err != nil {
		return nil, err
	}

	resolvedResponse := response.(interchange.VehicleDeleteResponse)

	return &pb.VehicleDeleteResponse{
		CountDeleted: resolvedResponse.CountDeleted,
	}, err
}
