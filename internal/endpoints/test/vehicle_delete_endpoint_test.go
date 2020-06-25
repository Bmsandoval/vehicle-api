package test

import (
	"context"
	"github.com/bmsandoval/vehicle-api/internal/endpoints"
	"github.com/bmsandoval/vehicle-api/internal/interchange"
	"github.com/bmsandoval/vehicle-api/pkg/appcontext"
	"github.com/bmsandoval/vehicle-api/pkg/services"
	"github.com/bmsandoval/vehicle-api/pkg/services/cache_service"
	"github.com/bmsandoval/vehicle-api/pkg/services/vehicle_service"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

type vehicleDeleteEndpointTestData struct {
	Description string
	Request interchange.VehicleDeleteRequest
	MockVehicleDeleteServiceResponse int64
	Response interchange.VehicleDeleteResponse
}
func TestVehicleDeleteEndpoint(t *testing.T) {
	TestData := []vehicleDeleteEndpointTestData{
		{
			Description: "happy path",
			Request: interchange.VehicleDeleteRequest{
				Id: 1,
			},
			MockVehicleDeleteServiceResponse: 1,
			Response: interchange.VehicleDeleteResponse{
				StatusCode: 200,
				CountDeleted: 1,
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	for _, data := range TestData {
		appCtx := appcontext.Context{
			GoContext: context.Background(),
		}

		bundle := MockVehicleDeleteRequiredServices(mockCtrl, data)

		t.Run(data.Description, func(t *testing.T) {
			f := endpoints.MakeVehicleDeleteEndpoint(appCtx, bundle)

			responseData, _ := f(context.Background(), data.Request)

			assert.Equal(t, data.Response, responseData)
		})
	}
}

func MockVehicleDeleteRequiredServices(mockCtrl *gomock.Controller, data vehicleDeleteEndpointTestData) services.Bundle {
	vehicleMock := vehicle_service.NewMockVehicleService(mockCtrl)
	vehicleExpect := vehicleMock.EXPECT()
	cacheMock := cache_service.NewMockCacheService(mockCtrl)
	cacheExpect := cacheMock.EXPECT()

	cacheExpect.Delete(string(data.Request.Id))
	vehicleExpect.Delete(data.Request.Id).Return(data.MockVehicleDeleteServiceResponse, nil)

	serviceBundle := services.Bundle{
		VehicleSvc: vehicleMock,
		CacheSvc: cacheMock,
	}

	return serviceBundle
}
