package vehicle_service

import (
	_ "github.com/lib/pq"
)

func (s ServiceImplementation) Delete(id int64) (int64, error) {
	statement, err := s.DB.Prepare(`DELETE FROM vehicles WHERE (id = $1);`)
	if err != nil {
		return 0, err
	}

	result, err := statement.Exec(id)
	if err != nil {
		return 0, err
	}

	count, err := result.RowsAffected()
	return count, err
}
