package main

import (
	"database/sql"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"

	"github.com/mbarlow/nit/core"
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

	// Initialize the core package with the database
	core.InitDB(db)

	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())

	// Single route handles all CRUD operations
	e.Any("/:table", core.HandleCRUD)
	e.Any("/:table/:id", core.HandleCRUD)

	logrus.Info("Server starting on :8080")
	logrus.Fatal(e.Start(":8080"))
}
