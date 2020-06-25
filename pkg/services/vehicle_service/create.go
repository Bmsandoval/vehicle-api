package vehicle_service

import (
	"github.com/bmsandoval/vehicle-api/pkg/models"
	_ "github.com/lib/pq"
)

func (s ServiceImplementation) Create(vehicle *models.Vehicle) (*models.Vehicle, error) {
	statement, err := s.DB.Prepare(`INSERT INTO vehicles(make, model, vin) VALUES
			( $1,$2,$3 ) RETURNING id;`)
	if err != nil {
		return nil, err
	}

	rows, err := statement.Query(vehicle.Make, vehicle.Model, vehicle.Vin)
	if err != nil {
		return nil, err
	}

	var id int64
	for rows.Next() {
		rows.Scan(&id)
	}

	vehicle.Id = id

	return vehicle, nil
}
