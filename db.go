package main

import (
	"database/sql"
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

	// Таблица проектов
	_, _ = db.Exec(`CREATE TABLE IF NOT EXISTS projects (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT,
		start_date TEXT,
		end_date TEXT,
		budget REAL DEFAULT 0,
		spent REAL DEFAULT 0,
		status TEXT DEFAULT 'in_progress',
		progress INTEGER DEFAULT 0
	)`)

	// Таблица объектов с расширенными характеристиками
	_, _ = db.Exec(`CREATE TABLE IF NOT EXISTS objects (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		project_id INTEGER,
		name TEXT NOT NULL,
		type TEXT,
		area REAL DEFAULT 0,
		budget REAL DEFAULT 0,
		spent REAL DEFAULT 0,
		progress INTEGER DEFAULT 0,
		floors INTEGER DEFAULT 0,
		material TEXT,
		status TEXT DEFAULT 'in_progress',
		FOREIGN KEY(project_id) REFERENCES projects(id)
	)`)

	// Таблица задач
	_, _ = db.Exec(`CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		project_id INTEGER,
		name TEXT NOT NULL,
		start_date TEXT,
		end_date TEXT,
		assigned_to TEXT,
		estimated REAL DEFAULT 0,
		spent REAL DEFAULT 0,
		progress INTEGER DEFAULT 0,
		status TEXT DEFAULT 'in_progress',
		FOREIGN KEY(project_id) REFERENCES projects(id)
	)`)

	// Добавляем тестовые данные, если база пустая
	var count int
	db.QueryRow("SELECT COUNT(*) FROM projects").Scan(&count)
	if count == 0 {
		db.Exec(`INSERT INTO projects (name, description, start_date, end_date, budget, spent, status, progress) VALUES 
			('Проектирование ТЦ "Север"', 'Разработка ПД и РД', '2025-01-15', '2025-06-30', 8500000, 3400000, 'in_progress', 45),
			('Строительство ЖК "Лесной"', 'Жилой комплекс 12 этажей', '2025-03-01', '2026-02-28', 24500000, 8200000, 'in_progress', 35)`)

		db.Exec(`INSERT INTO objects (project_id, name, type, area, budget, spent, progress, floors, material, status) VALUES 
			(1, 'Корпус А', 'Торговое здание', 12500, 3200000, 1400000, 45, 4, 'Железобетон', 'in_progress'),
			(1, 'Парковка', 'Наружные работы', 4500, 800000, 650000, 80, 0, 'Асфальт', 'completed'),
			(2, 'Фундамент', 'Нулевой цикл', 18000, 4500000, 2100000, 50, 0, 'Железобетон', 'in_progress')`)
	}
}