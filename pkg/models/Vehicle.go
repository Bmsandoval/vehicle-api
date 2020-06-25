package models

type Vehicle struct {
	Id int64 `json:"id" xml:"id"`
	Make string `json:"make" xml:"make"`
	Model string `json:"model" xml:"model"`
	Vin string `json:"vin" xml:"vin"`
}
