package interchange

import "github.com/bmsandoval/vehicle-api/pkg/models"

/// VEHICLE CREATE
type VehicleCreateRequest struct {
	Make string `json:"make" xml:"make"`
	Model string `json:"model" xml:"model"`
	Vin string `json:"vin" xml:"vin"`
}
type VehicleCreateResponse struct {
	StatusCode int `json:"-" xml:"-"`
	models.Vehicle
}

/// VEHICLE READ
type VehicleReadRequest struct {
	Id int64
}
type VehicleReadResponse struct {
	StatusCode int `json:"-" xml:"-"`
	models.Vehicle
}

/// VEHICLE SEARCH
type VehicleSearchRequest struct {
	Make string `json:"make" xml:"make"`
	Model string `json:"model" xml:"model"`
}

/// VEHICLE UPDATE
type VehicleUpdateRequest struct {
	models.Vehicle
}
type VehicleUpdateResponse struct {
	StatusCode int `json:"-" xml:"-"`
	CountUpdated int64
}

/// VEHICLE DELETE
type VehicleDeleteRequest struct {
	Id int64 `json:"id" xml:"id"`
}
type VehicleDeleteResponse struct {
	StatusCode int `json:"-" xml:"-"`
	CountDeleted int64
}
