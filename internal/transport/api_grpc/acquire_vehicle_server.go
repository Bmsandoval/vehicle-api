package api_grpc

import (
	"github.com/bmsandoval/vehicle-api/pkg/appcontext"
	"github.com/bmsandoval/vehicle-api/pkg/services"
	pb "github.com/bmsandoval/vehicle-api/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GrpcServer struct{
	Context appcontext.Context
	ServicesBundle services.Bundle
}

// Create the GRPC server
func AcquireVehicleServer(ctx appcontext.Context, servicesBundle services.Bundle) *grpc.Server {
	grpcS := grpc.NewServer()

	pb.RegisterVehiclesServer(grpcS, &GrpcServer{
		Context: ctx,
		ServicesBundle: servicesBundle,
	})

	reflection.Register(grpcS)

	return grpcS
}

