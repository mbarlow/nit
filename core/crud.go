package core

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func getAll(c echo.Context, table string) error {
	// Parse pagination parameters
	limitStr := c.QueryParam("limit")
	offsetStr := c.QueryParam("offset")

	limit := 10 // default limit
	offset := 0 // default offset

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
			if limit > 100 { // max limit to prevent abuse
				limit = 100
			}
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	// Get total count
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s", table)
	var totalItems int
	err := db.QueryRow(countQuery).Scan(&totalItems)
	if err != nil {
		return err
	}

	// Get paginated results
	query := fmt.Sprintf("SELECT id, data, created, updated FROM %s ORDER BY created DESC LIMIT ? OFFSET ?", table)
	rows, err := db.Query(query, limit, offset)
	if err != nil {
		return err
	}
	defer rows.Close()

	items := []map[string]interface{}{}
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
		items = append(items, result)
	}

	// Build paginated response
	response := map[string]interface{}{
		"items":       items,
		"total_items": totalItems,
		"limit":       limit,
		"offset":      offset,
		"has_more":    offset+limit < totalItems,
	}

	return c.JSON(http.StatusOK, response)
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
