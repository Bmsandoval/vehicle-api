package vehicle_service

import (
	"github.com/bmsandoval/vehicle-api/pkg/models"
	_ "github.com/lib/pq"
)

func (s ServiceImplementation) Update(vehicle *models.Vehicle) (int64, error) {
	statement, err := s.DB.Prepare(`
UPDATE vehicles
SET make = $1, model = $2, vin = $3
WHERE id = $4;`)
	if err != nil {
		return 0, err
	}

	result, err := statement.Exec(vehicle.Make, vehicle.Model, vehicle.Vin, vehicle.Id)
	if err != nil {
		return 0, err
	}

	count, err := result.RowsAffected()

	return count, err
}
