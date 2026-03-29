package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "./db.sqlite")
	if err != nil {
		log.Fatal(err)
	}

	// PROJECTS
	db.Exec(`CREATE TABLE IF NOT EXISTS projects (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		description TEXT,
		start_date TEXT,
		end_date TEXT,
		budget REAL,
		spent REAL,
		status TEXT,
		progress INTEGER
	)`)

	// OBJECTS
	db.Exec(`CREATE TABLE IF NOT EXISTS objects (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		project_id INTEGER,
		name TEXT,
		type TEXT,
		area REAL,
		floors INTEGER,
		material TEXT,
		budget REAL,
		spent REAL,
		status TEXT
	)`)

	// STAGES (НОВОЕ)
	db.Exec(`CREATE TABLE IF NOT EXISTS stages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		object_id INTEGER,
		type TEXT,
		name TEXT,
		budget REAL,
		spent REAL,
		progress INTEGER
	)`)

	// WORKS (НОВОЕ)
	db.Exec(`CREATE TABLE IF NOT EXISTS works (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		stage_id INTEGER,
		name TEXT,
		start_date TEXT,
		end_date TEXT,
		assigned_to TEXT,
		estimated REAL,
		spent REAL,
		progress INTEGER,
		status TEXT
	)`)

	fmt.Println("✅ DB готова")
}