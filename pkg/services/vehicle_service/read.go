package vehicle_service

import (
	"fmt"
	"github.com/bmsandoval/vehicle-api/pkg/models"
	_ "github.com/lib/pq"
)

func (s ServiceImplementation) Read(id int64) (*models.Vehicle, error) {
	query, err := s.DB.Prepare(`SELECT id,make,model,vin FROM "vehicles" WHERE (id = $1)`)
	if err != nil {
		return nil, err
	}
	rows, err := query.Query(id)
	if err != nil {
		return nil, err
	}

	var vehicle *models.Vehicle
	for rows.Next() {
		vehicle = &models.Vehicle{}
		err = rows.Scan(&vehicle.Id, &vehicle.Make, &vehicle.Model, &vehicle.Vin)
		if err != nil {
			fmt.Println("Failed to scan row", err)
			return nil, err
		}
	}

	return vehicle, nil
}
