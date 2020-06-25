package api_http

import (
	"github.com/bmsandoval/vehicle-api/pkg/appcontext"
	"github.com/bmsandoval/vehicle-api/pkg/services"
	"github.com/gorilla/mux"
	"net/http"
)

// Create the http server
func AcquireVehicleServer(ctx appcontext.Context, servicesBundle services.Bundle) *http.Server {
	router := mux.NewRouter()

	MakeVehicleCreateHttpHandler(ctx, router, servicesBundle)
	MakeVehicleReadHttpHandler(ctx, router, servicesBundle)
	MakeVehicleUpdateHttpHandler(ctx, router, servicesBundle)
	MakeVehicleDeleteHttpHandler(ctx, router, servicesBundle)

	httpS := &http.Server{
		Handler: router,
	}

	return httpS
}
