package core

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func getAll(c echo.Context, table string) error {
	query := fmt.Sprintf("SELECT id, data, created, updated FROM %s", table)
	rows, err := db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var id, data, created, updated string
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
	var idVal, data, created, updated string

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

	id := uuid.New().String()
	jsonData, _ := json.Marshal(data)
	query := fmt.Sprintf("INSERT INTO %s (id, data) VALUES (?, json(?))", table)

	_, err := db.Exec(query, id, string(jsonData))
	if err != nil {
		return err
	}

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
