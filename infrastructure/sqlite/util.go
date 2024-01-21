package sqlite

import (
	"database/sql"
)

const GoquDialect = "sqlite3_with_returning"

func scanRows[ModelType any](rows *sql.Rows) ([]ModelType, error) {
	defer rows.Close()

	var models []ModelType

	for rows.Next() {
		var model ModelType
		err := rows.Scan(&model)
		if err != nil {
			return nil, nil
		}

		models = append(models, model)
	}

	return models, nil
}
