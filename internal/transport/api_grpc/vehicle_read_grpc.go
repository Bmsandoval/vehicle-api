package api_grpc

import (
	"context"
	"github.com/bmsandoval/vehicle-api/internal/endpoints"
	"github.com/bmsandoval/vehicle-api/internal/interchange"
	pb "github.com/bmsandoval/vehicle-api/protos"
)

// GRPC Read endpoint wrapper
func (s GrpcServer) VehicleRead(ctx context.Context, in *pb.VehicleReadRequest) (*pb.VehicleReadResponse, error) {
	endpoint := endpoints.MakeVehicleReadEndpoint(s.Context, s.ServicesBundle)
	requestObj := interchange.VehicleReadRequest{ Id: in.Id }
	response, err := endpoint(s.Context.GoContext, requestObj)
	if err != nil {
		return nil, err
	}

	resolvedResponse := response.(interchange.VehicleReadResponse)

	return &pb.VehicleReadResponse{
		Id: resolvedResponse.Id,
		Make:  resolvedResponse.Make,
		Model: resolvedResponse.Model,
		Vin:   resolvedResponse.Vin,
	}, nil
}
