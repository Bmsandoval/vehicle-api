package test

import (
	"context"
	"github.com/bmsandoval/vehicle-api/internal/endpoints"
	"github.com/bmsandoval/vehicle-api/internal/interchange"
	"github.com/bmsandoval/vehicle-api/pkg/appcontext"
	"github.com/bmsandoval/vehicle-api/pkg/models"
	"github.com/bmsandoval/vehicle-api/pkg/services"
	"github.com/bmsandoval/vehicle-api/pkg/services/cache_service"
	"github.com/bmsandoval/vehicle-api/pkg/services/vehicle_service"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

type vehicleCreateEndpointTestData struct {
	Description string
	Request interchange.VehicleCreateRequest
	MockVehicleCreateServiceResponse *models.Vehicle
	Response interchange.VehicleCreateResponse
}
func TestVehicleCreateEndpoint(t *testing.T) {
	TestData := []vehicleCreateEndpointTestData{
		{
			Description: "happy path",
			Request: interchange.VehicleCreateRequest{
				Make: "Make",
				Model: "Model",
				Vin: "Vin",
			},
			MockVehicleCreateServiceResponse: &models.Vehicle{
				Id: 1,
				Make: "Make",
				Model: "Model",
				Vin: "Vin",
			},
			Response: interchange.VehicleCreateResponse{
				StatusCode: 200,
				Vehicle: models.Vehicle{
					Id: 1,
					Make:  "Make",
					Model: "Model",
					Vin:   "Vin",
				},
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	for _, data := range TestData {
		appCtx := appcontext.Context{
			GoContext: context.Background(),
		}

		bundle := MockVehicleCreateRequiredServices(mockCtrl, data)

		t.Run(data.Description, func(t *testing.T) {
			f := endpoints.MakeVehicleCreateEndpoint(appCtx, bundle)

			responseData, _ := f(context.Background(), data.Request)

			assert.Equal(t, data.Response, responseData)
		})
	}
}

func MockVehicleCreateRequiredServices(mockCtrl *gomock.Controller, data vehicleCreateEndpointTestData) services.Bundle {
	vehicleMock := vehicle_service.NewMockVehicleService(mockCtrl)
	vehicleExpect := vehicleMock.EXPECT()
	cacheMock := cache_service.NewMockCacheService(mockCtrl)
	cacheExpect := cacheMock.EXPECT()

	vehicleExpect.Create(&models.Vehicle{
		Make: data.Request.Make,
		Model: data.Request.Model,
		Vin: data.Request.Vin,
	}).Return(data.MockVehicleCreateServiceResponse, nil)
	cacheExpect.Set(string(data.MockVehicleCreateServiceResponse.Id), data.MockVehicleCreateServiceResponse).Return(nil)

	serviceBundle := services.Bundle{
		VehicleSvc: vehicleMock,
		CacheSvc: cacheMock,
	}

	return serviceBundle
}
