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
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// Delete http endpoint wrapper
func MakeVehicleDeleteHttpHandler(appCtx appcontext.Context, router *mux.Router, servicesBundle services.Bundle) {
	api := router.PathPrefix("/api").Subrouter()

	endpoint := endpoints.MakeVehicleDeleteEndpoint(appCtx, servicesBundle)
	decoder, _ := MakeVehicleDeleteHttpRequestDecoder(appCtx)
	encoder, _ := MakeVehicleDeleteHttpResponseEncoder(appCtx)

	api.Methods("DELETE").Path("/vehicle").Handler(httpTransport.NewServer(
		endpoint,
		decoder,
		encoder,
	))
}

// Http Request Decoder. Tracks content-type and converts to the endpoint interchange object
func MakeVehicleDeleteHttpRequestDecoder(appCtx appcontext.Context) (kithttp.DecodeRequestFunc, error) {
	return func(ctx context.Context, r *http.Request) (interface{}, error) {
		var Req interchange.VehicleDeleteRequest

		contentTypes, ok := r.Header["Content-Type"]
		var contentType string
		if !ok {
			// default to json response
			contentType = "json"
		}

		for _, content := range contentTypes {
			if strings.Contains(content, "json") {
				contentType = "json"
				decode := json.NewDecoder(r.Body)
				err := decode.Decode(&Req)
				if err != nil {
					return nil, err
				}
				break
			} else if strings.Contains(content, "xml") {
				contentType = "xml"
				decode := xml.NewDecoder(r.Body)
				err := decode.Decode(&Req)
				if err != nil {
					return nil, err
				}
				break
			} else if strings.Contains(content, "protobuf") {
				contentType = "protobuf"
				var msg pb.VehicleDeleteRequest
				body, err := ioutil.ReadAll(r.Body)
				if err != nil {
					return nil, err
				}
				if err := proto.Unmarshal(body, &msg); err != nil {
					return nil, err
				}
				Req.Id = msg.Id
				break
			}
		}

		// store content type
		appCtx.Viper.Set("content-type", contentType)

		return Req, nil
	}, nil
}

// Http Response Encoder. Takes an endpoint interchange object and converts it to the proper content-type before writing
func MakeVehicleDeleteHttpResponseEncoder(appCtx appcontext.Context) (kithttp.EncodeResponseFunc, error) {
	return func(ctx context.Context, httpResponse http.ResponseWriter, endpointResponse interface{}) error {
		if endpointResponse != nil {
			response, ok := endpointResponse.(interchange.VehicleDeleteResponse)
			if !ok {
				return errors.New("unexpected response type returned from endpoint")
			}

			if response.StatusCode > 0 {
				httpResponse.WriteHeader(response.StatusCode)
			}

			contentType, ok := appCtx.Viper.Get("content-type").(string)
			if ! ok {
				return errors.New("non-string content-type found in viper")
			}

			// Convert to proper content-type and return
			if contentType == "protobuf" {
				marshaledResponse, err := proto.Marshal(&pb.VehicleDeleteResponse{
					CountDeleted: response.CountDeleted,
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
