package test

import (
	"context"
	"github.com/bmsandoval/vehicle-api/internal/interchange"
	"github.com/bmsandoval/vehicle-api/internal/transport/api_http"
	"github.com/bmsandoval/vehicle-api/pkg/appcontext"
	"github.com/bmsandoval/vehicle-api/pkg/models"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type vehicleReadHttpRequestDecoderTestData struct {
	Description string
	EncodedRequest http.Request
	DecodedRequest interface{}
	Error bool
}
func TestVehicleReadHttpRequestDecoder(t *testing.T) {
	TestData := []vehicleReadHttpRequestDecoderTestData{
		{
			Description: "happy path",
			EncodedRequest: http.Request{
				Method:           "GET",
				Header:           map[string][]string{
					"Content-Type": {"json"},
				},
				Form:             url.Values{
					"id": {"1"},
				},
			},
			DecodedRequest: interchange.VehicleReadRequest{
				Id: 1,
			},
			Error: false,
		},
	}

	for _, data := range TestData {
		appCtx := appcontext.Context{
			Viper: viper.New(),
			GoContext: context.Background(),
		}
		t.Run(data.Description, func(t *testing.T) {
			f, _ := api_http.MakeVehicleReadHttpRequestDecoder(appCtx)
			responseTestData, err := f(appCtx.GoContext, &data.EncodedRequest)

			assert.Equal(t, data.Error, err != nil)

			assert.Equal(t, responseTestData, data.DecodedRequest)
		})
	}
}

/// ***************** RESPONSE ENCODER
type vehicleReadHttpResponseEncoderTestData struct {
	Description string
	DecodedResponse interface{}
	EncodedResponseContentType string
	EncodedResponseBody []byte
	Error bool
}
func TestVehicleReadHttpResponseEncoder(t *testing.T) {
	TestData := []vehicleReadHttpResponseEncoderTestData{
		{
			Description: "json happy path",
			DecodedResponse: interchange.VehicleReadResponse{
				Vehicle: models.Vehicle{
					Id: 1,
					Make:  "Make",
					Model: "Model",
					Vin:   "Vin",
				},
			},
			EncodedResponseContentType: "json",
			EncodedResponseBody:        []byte(`{"id":1,"make":"Make","model":"Model","vin":"Vin"}`),
			Error:                      false,
		},
		{
				Description: "xml happy path",
				DecodedResponse: interchange.VehicleReadResponse{
					Vehicle: models.Vehicle{
						Id: 1,
						Make:  "Make",
						Model: "Model",
						Vin:   "Vin",
					},
				},
				EncodedResponseContentType: "xml",
				EncodedResponseBody: []byte(`<VehicleReadResponse><id>1</id><make>Make</make><model>Model</model><vin>Vin</vin></VehicleReadResponse>`),
				Error: false,
		},
	}

	for _, data := range TestData {
		appCtx := appcontext.Context{
			Viper: viper.New(),
			GoContext: context.Background(),
		}
		appCtx.Viper.Set("content-type", data.EncodedResponseContentType)
		t.Run(data.Description, func(t *testing.T) {
			responseRecorder := httptest.NewRecorder()
			f, _ := api_http.MakeVehicleReadHttpResponseEncoder(appCtx)
			err := f(appCtx.GoContext, responseRecorder, data.DecodedResponse)

			assert.Equal(t, data.Error, err != nil)

			actualResponseBody, _ := ioutil.ReadAll(responseRecorder.Body)

			assert.Equal(t, data.EncodedResponseBody,
				actualResponseBody)
		})
	}
}
