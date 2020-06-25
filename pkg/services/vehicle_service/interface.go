package vehicle_service

import (
	"database/sql"
	"github.com/bmsandoval/vehicle-api/pkg/appcontext"
	"github.com/bmsandoval/vehicle-api/pkg/models"
	_ "github.com/lib/pq"
)

type ServiceImplementation struct {
	AppCtx appcontext.Context
	DB *sql.DB
}

type VehicleService interface {
	Create(*models.Vehicle) (*models.Vehicle, error)
	Read(int64) (*models.Vehicle, error)
	Search(limit int, offset int, vMake string, vModel string) ([]*models.Vehicle, error)
	Update(*models.Vehicle) (int64, error)
	Delete(id int64) (int64, error)
}
