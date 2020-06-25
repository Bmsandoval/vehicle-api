package api_http

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"github.com/bmsandoval/vehicle-api/internal/endpoints"
	"github.com/bmsandoval/vehicle-api/internal/interchange"
	"github.com/bmsandoval/vehicle-api/pkg/appcontext"
	"github.com/bmsandoval/vehicle-api/pkg/services"
	pb "github.com/bmsandoval/vehicle-api/protos"
	httpTransport "github.com/go-kit/kit/transport/http"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// Read http endpoint wrapper
func MakeVehicleReadHttpHandler(appCtx appcontext.Context, router *mux.Router, servicesBundle services.Bundle) {
	api := router.PathPrefix("/api").Subrouter()

	endpoint := endpoints.MakeVehicleReadEndpoint(appCtx, servicesBundle)
	decoder, _ := MakeVehicleReadHttpRequestDecoder(appCtx)
	encoder, _ := MakeVehicleReadHttpResponseEncoder(appCtx)

	api.Methods("GET").Path("/vehicle").Handler(httpTransport.NewServer(
		endpoint,
		decoder,
		encoder,
	))
}

// Http Request Decoder. Tracks content-type and converts to the endpoint interchange object
func MakeVehicleReadHttpRequestDecoder(appCtx appcontext.Context) (kithttp.DecodeRequestFunc, error) {
	return func(ctx context.Context, r *http.Request) (interface{}, error) {
		var Req interchange.VehicleReadRequest


		contentTypes, ok := r.Header["Content-Type"]
		var contentType string
		if !ok {
			// default to json response
			contentType = "json"
		}
		for _, content := range contentTypes {
			if strings.Contains(content, "json") {
				contentType = "json"
				break
			} else if strings.Contains(content, "xml") {
				contentType = "xml"
				break
			} else if strings.Contains(content, "protobuf") {
				contentType = "protobuf"
				break
			}
		}

		appCtx.Viper.Set("content-type", contentType)

		id, err := strconv.ParseInt(r.FormValue("id"), 10, 64)
		if err != nil {
			return nil, err
		}

		Req.Id = id

		return Req, nil
	}, nil
}

// Http Response Encoder. Takes an endpoint interchange object and converts it to the proper content-type before writing
func MakeVehicleReadHttpResponseEncoder(appCtx appcontext.Context) (kithttp.EncodeResponseFunc, error) {
	return func(ctx context.Context, httpResponse http.ResponseWriter, endpointResponse interface{}) error {
		if endpointResponse != nil {
			response, ok := endpointResponse.(interchange.VehicleReadResponse)
			if !ok {
				return errors.New("unexpected response type returned from endpoint")
			}

			if response.StatusCode > 0 {
				httpResponse.WriteHeader(response.StatusCode)
			}

			// Set Cache Headers
			httpResponse.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
			v := "MISS"
			if response.StatusCode == 304 {
				v = "HIT"
			}
			httpResponse.Header().Set("X-Cache", v)

			contentType, ok := appCtx.Viper.Get("content-type").(string)
			if ! ok {
				return errors.New("non-string content-type found in viper")
			}

			// Convert to proper content-type and return
			if contentType == "protobuf" {
				marshaledResponse, err := proto.Marshal(&pb.VehicleReadResponse{
					Id: response.Vehicle.Id,
					Make:  response.Vehicle.Make,
					Model: response.Vehicle.Model,
					Vin:   response.Vehicle.Vin,
				})
				if err != nil {
					log.Println(err.Error())
				}
				if _, err = httpResponse.Write(marshaledResponse); err != nil {
					return err
				}
			} else {
				var err error
				var marshaledResponse []byte
				if contentType == "json" {
					marshaledResponse, err = json.Marshal(response)
					if err != nil {
						return err
					}
				} else if contentType == "xml" {
					marshaledResponse, err = xml.Marshal(response)
					if err != nil {
						return err
					}
				} else {
					return errors.New("unexpected content type")
				}
				if _, err = httpResponse.Write(marshaledResponse); err != nil {
					return err
				}
			}
		}

		return nil

	}, nil
}
