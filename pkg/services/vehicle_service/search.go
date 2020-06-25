package vehicle_service

import (
	"fmt"
	"github.com/bmsandoval/vehicle-api/pkg/models"
	_ "github.com/lib/pq"
)

func (s ServiceImplementation) Search(limit int, offset int, vMake string, vModel string) ([]*models.Vehicle, error) {
	query := `SELECT id,make,model,vin FROM "vehicles" WHERE (make = $1) AND (model = $2) LIMIT $3 OFFSET $4`

	statement, err := s.DB.Prepare(query)
	if err != nil {
		return nil, err
	}
	rows, err := statement.Query(vMake, vModel, limit, offset)
	if err != nil {
		return nil, err
	}

	var vehicles []*models.Vehicle
	for rows.Next() {
		vehicle := &models.Vehicle{}
		err = rows.Scan(&vehicle.Id, &vehicle.Make, &vehicle.Model, &vehicle.Vin)
		if err != nil {
			fmt.Println("Failed to scan row", err)
			return nil, err
		}
		vehicles = append(vehicles, vehicle)
	}

	return vehicles, nil
}
