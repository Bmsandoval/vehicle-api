package test

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bmsandoval/vehicle-api/pkg/appcontext"
	"github.com/bmsandoval/vehicle-api/pkg/services/vehicle_service"
	"github.com/magiconair/properties/assert"
	"regexp"
	"testing"
)

func TestVehicleDeleteService(t *testing.T) {
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

	mock.ExpectPrepare(regexp.QuoteMeta(`DELETE FROM vehicles WHERE (id = $1);`)).
		ExpectExec().WithArgs(1).WillReturnResult(sqlmock.NewResult(1,1))

	// now we execute our method
	result, err := vehicleService.Delete(1)
	if err != nil {
		t.Errorf(err.Error())
	}

	assert.Equal(t, result, int64(1))

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
