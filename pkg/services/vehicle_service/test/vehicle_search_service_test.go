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

func TestVehicleSearchService(t *testing.T) {
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

	mock.ExpectPrepare(regexp.QuoteMeta(`SELECT id,make,model,vin FROM "vehicles" WHERE (make = $1) AND (model = $2) LIMIT $3 OFFSET $4`)).
		ExpectQuery().WithArgs("Make", "Model", 2, 2).
		WillReturnRows(sqlmock.NewRows([]string{"id", "make", "model", "vin"}).
			AddRow(3, "Make", "Model", "Vin").
			AddRow(4, "Make", "Model", "Vin"))

	// now we execute our method
	result, err := vehicleService.Search(2, 2, "Make", "Model")
	if err != nil {
		t.Errorf(err.Error())
	}

	assert.Equal(t, result, []*models.Vehicle{
		{
			Id: 3,
			Make:  "Make",
			Model: "Model",
			Vin:   "Vin",
		},
		{
			Id: 4,
			Make:  "Make",
			Model: "Model",
			Vin:   "Vin",
		},
	})

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
