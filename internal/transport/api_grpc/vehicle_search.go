package api_grpc

import (
	"github.com/bmsandoval/vehicle-api/internal/endpoints"
	"github.com/bmsandoval/vehicle-api/internal/interchange"
	pb "github.com/bmsandoval/vehicle-api/protos"
)

// GRPC Search endpoint wrapper
func (s GrpcServer) VehicleSearch(in *pb.VehicleSearchRequest, server pb.Vehicles_VehicleSearchServer) error {
	endpoint := endpoints.MakeVehicleSearchEndpoint(s.Context, server, s.ServicesBundle)
	err := endpoint(s.Context.GoContext, interchange.VehicleSearchRequest{
		Make:  in.Make,
		Model: in.Model,
	})

	if err != nil {
		return err
	}

	return nil
}
