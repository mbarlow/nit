package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
)

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("sqlite3", "./data.db?_foreign_keys=on")
	if err != nil {
		logrus.Fatal(err)
	}
	defer db.Close()

	// Enable JSON1 extension
	db.Exec("PRAGMA foreign_keys = ON")

	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())

	// Single route handles all CRUD operations
	e.Any("/:table", handleCRUD)
	e.Any("/:table/:id", handleCRUD)

	logrus.Info("Server starting on :8080")
	logrus.Fatal(e.Start(":8080"))
}

func handleCRUD(c echo.Context) error {
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

func ensureTable(table string) {
	query := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			data JSON NOT NULL,
			created DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated DATETIME DEFAULT CURRENT_TIMESTAMP
		)`, table)
	db.Exec(query)
}

func getAll(c echo.Context, table string) error {
	query := fmt.Sprintf("SELECT id, data, created, updated FROM %s", table)
	rows, err := db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var id int
		var data, created, updated string
		rows.Scan(&id, &data, &created, &updated)

		var jsonData map[string]interface{}
		json.Unmarshal([]byte(data), &jsonData)

		result := map[string]interface{}{
			"id":      id,
			"data":    jsonData,
			"created": created,
			"updated": updated,
		}
		results = append(results, result)
	}

	return c.JSON(http.StatusOK, results)
}

func getOne(c echo.Context, table, id string) error {
	query := fmt.Sprintf("SELECT id, data, created, updated FROM %s WHERE id = ?", table)
	var idVal int
	var data, created, updated string

	err := db.QueryRow(query, id).Scan(&idVal, &data, &created, &updated)
	if err == sql.ErrNoRows {
		return echo.NewHTTPError(http.StatusNotFound)
	}
	if err != nil {
		return err
	}

	var jsonData map[string]interface{}
	json.Unmarshal([]byte(data), &jsonData)

	result := map[string]interface{}{
		"id":      idVal,
		"data":    jsonData,
		"created": created,
		"updated": updated,
	}

	return c.JSON(http.StatusOK, result)
}

func create(c echo.Context, table string) error {
	var data map[string]interface{}
	if err := c.Bind(&data); err != nil {
		return err
	}

	jsonData, _ := json.Marshal(data)
	query := fmt.Sprintf("INSERT INTO %s (data) VALUES (json(?))", table)

	result, err := db.Exec(query, string(jsonData))
	if err != nil {
		return err
	}

	id, _ := result.LastInsertId()
	return c.JSON(http.StatusCreated, map[string]interface{}{"id": id})
}

func update(c echo.Context, table, id string) error {
	var data map[string]interface{}
	if err := c.Bind(&data); err != nil {
		return err
	}

	jsonData, _ := json.Marshal(data)
	query := fmt.Sprintf("UPDATE %s SET data = json(?), updated = CURRENT_TIMESTAMP WHERE id = ?", table)

	result, err := db.Exec(query, string(jsonData), id)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "updated"})
}

func delete(c echo.Context, table, id string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = ?", table)
	result, err := db.Exec(query, id)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "deleted"})
}
