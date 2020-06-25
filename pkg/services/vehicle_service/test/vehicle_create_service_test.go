package test

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bmsandoval/vehicle-api/pkg/appcontext"
	"github.com/bmsandoval/vehicle-api/pkg/models"
	"github.com/bmsandoval/vehicle-api/pkg/services/vehicle_service"
	"github.com/magiconair/properties/assert"
	"regexp"
	"testing"
)

func TestVehicleCreateService(t *testing.T) {
	var err error
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	vehicleService := vehicle_service.ServiceImplementation{
		AppCtx: appcontext.Context{},
		DB:     db,
	}

	mock.ExpectPrepare(regexp.QuoteMeta(` INSERT INTO vehicles(make, model, vin) VALUES ( $1,$2,$3 ) RETURNING id;`)).
		ExpectQuery().WithArgs("Make", "Model", "Vin").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	// now we execute our method
	result, err := vehicleService.Create(&models.Vehicle{
		Make:  "Make",
		Model: "Model",
		Vin:   "Vin",
	})
	if err != nil {
		t.Errorf(err.Error())
	}

	assert.Equal(t, result, &models.Vehicle{
		Id: 1,
		Make:  "Make",
		Model: "Model",
		Vin:   "Vin",
	})

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
