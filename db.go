package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

// Инициализация базы данных
func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "./db.sqlite")
	if err != nil {
		log.Fatal("Ошибка открытия базы данных:", err)
	}

	// Создание таблицы проектов
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS projects (
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
	if err != nil {
		log.Println("Ошибка создания таблицы projects:", err)
	}

	// Создание таблицы объектов (с расширенными характеристиками)
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS objects (
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
	if err != nil {
		log.Println("Ошибка создания таблицы objects:", err)
	}

	// Создание таблицы задач (график работ)
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS tasks (
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
	if err != nil {
		log.Println("Ошибка создания таблицы tasks:", err)
	}

	// Добавление тестовых данных (только если таблица projects пустая)
	var count int
	db.QueryRow("SELECT COUNT(*) FROM projects").Scan(&count)

	if count == 0 {
		fmt.Println("Добавляем тестовые данные...")

		// Тестовые проекты
		db.Exec(`INSERT INTO projects (name, description, start_date, end_date, budget, spent, status, progress) VALUES 
			('Проектирование ТЦ "Север"', 'Разработка проектной и рабочей документации для торгового центра', '2025-01-15', '2025-06-30', 8500000, 3400000, 'in_progress', 45),
			('Строительство ЖК "Лесной"', 'Возведение 12-этажного жилого комплекса', '2025-03-01', '2026-02-28', 24500000, 8200000, 'in_progress', 35)`)

		// Тестовые объекты
		db.Exec(`INSERT INTO objects (project_id, name, type, area, budget, spent, progress, floors, material, status) VALUES 
			(1, 'Корпус А', 'Торговое здание', 12500, 3200000, 1400000, 45, 4, 'Железобетон', 'in_progress'),
			(1, 'Парковка', 'Наружные работы', 4500, 800000, 650000, 80, 0, 'Асфальтобетон', 'completed'),
			(2, 'Фундамент', 'Нулевой цикл', 18000, 4500000, 2100000, 50, 0, 'Железобетон', 'in_progress')`)

		// Тестовые задачи
		db.Exec(`INSERT INTO tasks (project_id, name, start_date, end_date, assigned_to, estimated, spent, progress, status) VALUES 
			(1, 'Геологические изыскания', '2025-02-01', '2025-02-20', 'Иванов И.И.', 450000, 450000, 100, 'done'),
			(1, 'Согласование ПД', '2025-03-01', '2025-04-15', 'Петров С.В.', 1200000, 680000, 55, 'in_progress'),
			(2, 'Земляные работы', '2025-04-01', '2025-05-15', 'Сидоров А.А.', 2800000, 1100000, 40, 'in_progress')`)
	}

	fmt.Println("✅ База данных SQLite успешно инициализирована")
}