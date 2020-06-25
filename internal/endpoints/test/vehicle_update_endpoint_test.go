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

type vehicleUpdateEndpointTestData struct {
	Description string
	Request interchange.VehicleUpdateRequest
	MockVehicleReadServiceResponse *models.Vehicle
	MockVehicleUpdateServiceResponse int64
	Response interchange.VehicleUpdateResponse
}
func TestVehicleUpdateEndpoint(t *testing.T) {
	TestData := []vehicleUpdateEndpointTestData{
		{
			Description: "happy path",
			Request: interchange.VehicleUpdateRequest{
				Vehicle: models.Vehicle{
					Id:    1,
					Make:  "Make",
					Model: "Model",
					Vin:   "Vin",
				},
			},
			MockVehicleReadServiceResponse: &models.Vehicle{
				Id:    1,
				Make:  "Maker",
				Model: "Modeler",
				Vin:   "Vinny",
			},
			MockVehicleUpdateServiceResponse: 1,
			Response: interchange.VehicleUpdateResponse{
				StatusCode: 200,
				CountUpdated: 1,
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	for _, data := range TestData {
		appCtx := appcontext.Context{
			GoContext: context.Background(),
		}

		bundle := MockVehicleUpdateRequiredServices(mockCtrl, data)

		t.Run(data.Description, func(t *testing.T) {
			f := endpoints.MakeVehicleUpdateEndpoint(appCtx, bundle)

			responseData, _ := f(context.Background(), data.Request)

			assert.Equal(t, data.Response, responseData)
		})
	}
}

func MockVehicleUpdateRequiredServices(mockCtrl *gomock.Controller, data vehicleUpdateEndpointTestData) services.Bundle {
	vehicleMock := vehicle_service.NewMockVehicleService(mockCtrl)
	vehicleExpect := vehicleMock.EXPECT()
	cacheMock := cache_service.NewMockCacheService(mockCtrl)
	cacheExpect := cacheMock.EXPECT()

	updatedEntry := &models.Vehicle{
		Id: data.Request.Id,
		Make: data.Request.Make,
		Model: data.Request.Model,
		Vin: data.Request.Vin,
	}

	vehicleExpect.Read(updatedEntry.Id).Return(data.MockVehicleReadServiceResponse, nil)
	vehicleExpect.Update(updatedEntry).Return(data.MockVehicleUpdateServiceResponse, nil)
	cacheExpect.Set(string(updatedEntry.Id), updatedEntry)

	serviceBundle := services.Bundle{
		VehicleSvc: vehicleMock,
		CacheSvc: cacheMock,
	}

	return serviceBundle
}
