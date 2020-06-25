package test

import (
	"bytes"
	"context"
	"encoding/gob"
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

type vehicleReadEndpointTestData struct {
	Description string
	Request interchange.VehicleReadRequest
	MockCacheGetServiceResponse *models.Vehicle
	MockVehicleReadServiceResponse *models.Vehicle
	Response interchange.VehicleReadResponse
}
func TestVehicleReadEndpoint(t *testing.T) {
	TestData := []vehicleReadEndpointTestData{
		{
			Description: "cache miss",
			Request: interchange.VehicleReadRequest{
				Id: 1,
			},
			MockCacheGetServiceResponse: nil,
			MockVehicleReadServiceResponse: &models.Vehicle{
				Make: "Make",
				Model: "Model",
				Vin: "Vin",
			},
			Response: interchange.VehicleReadResponse{
				StatusCode: 200,
				Vehicle: models.Vehicle{
					Make:  "Make",
					Model: "Model",
					Vin:   "Vin",
				},
			},
		},
		{
			Description: "cache hit",
			Request: interchange.VehicleReadRequest{
				Id: 1,
			},
			MockCacheGetServiceResponse: &models.Vehicle{
				Make: "Make",
				Model: "Model",
				Vin: "Vin",
			},
			MockVehicleReadServiceResponse: nil,
			Response: interchange.VehicleReadResponse{
				StatusCode: 200,
				Vehicle: models.Vehicle{
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

		bundle := MockVehicleReadRequiredServices(mockCtrl, data)

		t.Run(data.Description, func(t *testing.T) {
			f := endpoints.MakeVehicleReadEndpoint(appCtx, bundle)

			responseData, _ := f(context.Background(), data.Request)

			assert.Equal(t, data.Response, responseData)
		})
	}
}

func MockVehicleReadRequiredServices(mockCtrl *gomock.Controller, data vehicleReadEndpointTestData) services.Bundle {
	vehicleMock := vehicle_service.NewMockVehicleService(mockCtrl)
	vehicleExpect := vehicleMock.EXPECT()
	cacheMock := cache_service.NewMockCacheService(mockCtrl)
	cacheExpect := cacheMock.EXPECT()

	if data.MockCacheGetServiceResponse != nil {
		var network bytes.Buffer        // Stand-in for a network connection
		enc := gob.NewEncoder(&network) // Will write to network.
		// Encode (send) the value.
		_ = enc.Encode(data.MockCacheGetServiceResponse)
		cachable := network.Bytes()

		cacheExpect.Get(string(data.Request.Id)).Return(cachable, nil)
	} else {
		cacheExpect.Get(string(data.Request.Id)).Return(nil, nil)
	}

	if data.MockVehicleReadServiceResponse != nil {
		vehicleExpect.Read(data.Request.Id).Return(data.MockVehicleReadServiceResponse, nil)
		cacheExpect.Set(string(data.MockVehicleReadServiceResponse.Id), data.MockVehicleReadServiceResponse).Return(nil)
	}

	serviceBundle := services.Bundle{
		VehicleSvc: vehicleMock,
		CacheSvc: cacheMock,
	}

	return serviceBundle
}
