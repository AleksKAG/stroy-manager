package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

type Project struct {
	ID          int
	Name        string
	Description string
	StartDate   string
	EndDate     string
	Budget      float64
	Spent       float64
	Status      string
	Progress    int
}

// Инициализация базы данных + тестовые данные
func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "./db.sqlite")
	if err != nil {
		log.Fatal(err)
	}

	// Создаём таблицы
	_, _ = db.Exec(`CREATE TABLE IF NOT EXISTS projects (
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

	_, _ = db.Exec(`CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		project_id INTEGER,
		name TEXT, 
		due_date TEXT,
		status TEXT, 
		assigned_to TEXT,
		progress INTEGER,
		FOREIGN KEY(project_id) REFERENCES projects(id)
	)`)

	// Добавляем тестовые данные, если таблица пустая
	var count int
	db.QueryRow("SELECT COUNT(*) FROM projects").Scan(&count)
	if count == 0 {
		db.Exec(`INSERT INTO projects (name, description, start_date, end_date, budget, spent, status, progress) VALUES 
			('Проектирование ТЦ "Север"', 'Разработка ПД и РД для торгового центра', '2025-01-15', '2025-06-30', 8500000, 3400000, 'in_progress', 45),
			('Строительство ЖК "Лесной"', 'Возведение 12-этажного жилого комплекса', '2025-03-01', '2026-02-28', 24500000, 8200000, 'in_progress', 35),
			('Реконструкция склада', 'Модернизация логистического комплекса', '2024-11-01', '2025-04-15', 4200000, 4100000, 'completed', 100)`)
	}
}

// ====================== ДАШБОРД ======================
func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	var totalProj, active int
	var totalBudget, totalSpent float64
	var avgProgress int

	db.QueryRow("SELECT COUNT(*), SUM(budget), SUM(spent) FROM projects").Scan(&totalProj, &totalBudget, &totalSpent)
	db.QueryRow("SELECT COUNT(*) FROM projects WHERE status != 'completed'").Scan(&active)
	db.QueryRow("SELECT COALESCE(AVG(progress), 0) FROM projects").Scan(&avgProgress)

	spentPercent := 0
	if totalBudget > 0 {
		spentPercent = int((totalSpent / totalBudget) * 100)
	}

	data := map[string]interface{}{
		"TotalProjects":  totalProj,
		"ActiveProjects": active,
		"TotalBudget":    fmt.Sprintf("%.0f", totalBudget),
		"TotalSpent":     fmt.Sprintf("%.0f", totalSpent),
		"SpentPercent":   spentPercent,
		"AvgProgress":    avgProgress,
	}

	tmpl := template.Must(template.New("dashboard").Parse(dashboardHTML))
	tmpl.Execute(w, data)
}

const dashboardHTML = `<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>СтроМенеджер</title>
    <script src="https://cdn.tailwindcss.com"></script>
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.5.1/css/all.min.css">
    <style>
        @import url('https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600&family=Space+Grotesk:wght@500;600&display=swap');
        body { font-family: 'Inter', sans-serif; }
        .title { font-family: 'Space Grotesk', sans-serif; }
    </style>
</head>
<body class="bg-slate-950 text-slate-100 min-h-screen">
    <div class="max-w-7xl mx-auto p-8">
        <div class="flex justify-between items-center mb-10">
            <div>
                <h1 class="text-5xl font-semibold title tracking-tight text-white">СтроМенеджер</h1>
                <p class="text-slate-400 text-lg">Управление строительными проектами</p>
            </div>
            <a href="/projects" 
               class="bg-orange-500 hover:bg-orange-600 px-8 py-4 rounded-2xl font-medium flex items-center gap-3 transition-colors">
                <i class="fas fa-list"></i> Все проекты
            </a>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-12">
            <div class="bg-slate-900 rounded-3xl p-8">
                <div class="text-slate-400">Всего проектов</div>
                <div class="text-6xl font-semibold mt-3">{{.TotalProjects}}</div>
            </div>
            <div class="bg-slate-900 rounded-3xl p-8">
                <div class="text-slate-400">Активных сейчас</div>
                <div class="text-6xl font-semibold mt-3 text-orange-400">{{.ActiveProjects}}</div>
            </div>
            <div class="bg-slate-900 rounded-3xl p-8">
                <div class="text-slate-400">Общий бюджет</div>
                <div class="text-6xl font-semibold mt-3">{{.TotalBudget}} ₽</div>
            </div>
            <div class="bg-slate-900 rounded-3xl p-8">
                <div class="text-slate-400">Потрачено</div>
                <div class="text-6xl font-semibold mt-3 text-emerald-400">{{.TotalSpent}} ₽</div>
                <div class="mt-4 h-2.5 bg-slate-700 rounded-full overflow-hidden">
                    <div class="h-2.5 bg-emerald-500 rounded-full" style="width: {{.SpentPercent}}%"></div>
                </div>
                <div class="text-sm text-slate-400 mt-2">{{.SpentPercent}}% от бюджета</div>
            </div>
        </div>

        <div class="text-center text-2xl font-medium">
            Средний прогресс по всем проектам: 
            <span class="text-orange-400 font-semibold">{{.AvgProgress}}%</span>
        </div>
    </div>
</body>
</html>`

// ====================== СПИСОК ПРОЕКТОВ ======================
func projectsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, name, status, progress, budget, spent FROM projects")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	var projects []Project
	for rows.Next() {
		var p Project
		rows.Scan(&p.ID, &p.Name, &p.Status, &p.Progress, &p.Budget, &p.Spent)
		projects = append(projects, p)
	}

	tmpl := template.Must(template.New("projects").Parse(projectsHTML))
	tmpl.Execute(w, projects)
}

const projectsHTML = `<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Проекты — СтроМенеджер</title>
    <script src="https://cdn.tailwindcss.com"></script>
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.5.1/css/all.min.css">
</head>
<body class="bg-slate-950 text-slate-100">
    <div class="max-w-7xl mx-auto p-8">
        <a href="/" class="inline-flex items-center gap-2 text-orange-400 mb-8 hover:text-orange-300">
            ← На главную
        </a>
        <h1 class="text-4xl font-semibold mb-10">Все проекты</h1>
        
        <div class="grid gap-6">
            {{range .}}
            <div class="bg-slate-900 rounded-3xl p-7 flex flex-col md:flex-row md:items-center gap-6">
                <div class="flex-1">
                    <h3 class="text-2xl font-medium">{{.Name}}</h3>
                    <div class="mt-4 h-3 bg-slate-700 rounded-full overflow-hidden">
                        <div class="h-3 bg-orange-500 rounded-full transition-all" 
                             style="width: {{.Progress}}%"></div>
                    </div>
                </div>
                <div class="text-right">
                    <div class="text-4xl font-semibold text-orange-400">{{.Progress}}%</div>
                    <div class="text-slate-400 text-sm">{{.Budget}} ₽</div>
                </div>
                <a href="/project/{{.ID}}" 
                   class="md:ml-6 bg-orange-500 hover:bg-orange-600 px-10 py-4 rounded-2xl font-medium transition-colors">
                    Открыть проект
                </a>
            </div>
            {{end}}
        </div>
    </div>
</body>
</html>`

// Заглушка для детальной страницы проекта
func projectHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `<div style="padding:50px;font-family:sans-serif">
		<h1>Детали проекта</h1>
		<p>Пока здесь заглушка. Хочешь — добавим задачи, бюджет по статьям и график?</p>
		<a href="/projects" style="color:#f97316">← Вернуться к проектам</a>
	</div>`)
}

func main() {
	initDB()
	defer db.Close()

	http.HandleFunc("/", dashboardHandler)
	http.HandleFunc("/projects", projectsHandler)
	http.HandleFunc("/project/", projectHandler)

	fmt.Println("🚀 СтроМенеджер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}