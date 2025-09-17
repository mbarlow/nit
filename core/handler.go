package core

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

var db *sql.DB

func InitDB(database *sql.DB) {
	db = database
}

func ensureTable(table string) {
	query := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			id TEXT PRIMARY KEY,
			data JSON NOT NULL,
			created DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated DATETIME DEFAULT CURRENT_TIMESTAMP
		)`, table)
	db.Exec(query)
}

func HandleCRUD(c echo.Context) error {
	table := c.Param("table")
	id := c.Param("id")
	method := c.Request().Method

	// Ensure table exists
	ensureTable(table)

	switch method {
	case "GET":
		if id == "" {
			return getAll(c, table)
		}
		return getOne(c, table, id)
	case "POST":
		return create(c, table)
	case "PUT":
		return update(c, table, id)
	case "DELETE":
		return delete(c, table, id)
	}

	return echo.NewHTTPError(http.StatusMethodNotAllowed)
}
