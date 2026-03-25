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
	ID       int
	Name     string
	Status   string
	Progress int
	Budget   float64
	Spent    float64
}

// Инициализация БД
func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "./db.sqlite")
	if err != nil {
		log.Fatal(err)
	}

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

	// Тестовые данные
	var count int
	db.QueryRow("SELECT COUNT(*) FROM projects").Scan(&count)
	if count == 0 {
		db.Exec(`INSERT INTO projects (name, description, start_date, end_date, budget, spent, status, progress) VALUES 
			('Проектирование ТЦ "Север"', 'Разработка ПД и РД', '2025-01-15', '2025-06-30', 8500000, 3400000, 'in_progress', 45),
			('Строительство ЖК "Лесной"', 'Жилой комплекс 12 этажей', '2025-03-01', '2026-02-28', 24500000, 8200000, 'in_progress', 35),
			('Реконструкция склада', 'Модернизация склада', '2024-11-01', '2025-04-15', 4200000, 4100000, 'completed', 100)`)
	}
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	var total, active int
	var budget, spent float64
	var avg int

	db.QueryRow("SELECT COUNT(*), SUM(budget), SUM(spent) FROM projects").Scan(&total, &budget, &spent)
	db.QueryRow("SELECT COUNT(*) FROM projects WHERE status != 'completed'").Scan(&active)
	db.QueryRow("SELECT COALESCE(AVG(progress), 0) FROM projects").Scan(&avg)

	spentPercent := 0
	if budget > 0 {
		spentPercent = int(spent / budget * 100)
	}

	data := map[string]interface{}{
		"Total":        total,
		"Active":       active,
		"Budget":       fmt.Sprintf("%.0f", budget),
		"Spent":        fmt.Sprintf("%.0f", spent),
		"SpentPercent": spentPercent,
		"AvgProgress":  avg,
	}

	tmpl := template.Must(template.New("dash").Parse(dashboardHTML))
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
<body class="bg-slate-950 text-slate-100">
    <div class="max-w-7xl mx-auto p-8">
        <h1 class="text-5xl font-semibold title mb-2">СтроМенеджер</h1>
        <p class="text-slate-400 mb-10">Управление проектированием и строительством</p>

        <div class="grid grid-cols-1 md:grid-cols-4 gap-6">
            <div class="bg-slate-900 rounded-3xl p-8">
                <div class="text-slate-400">Всего проектов</div>
                <div class="text-6xl font-bold mt-2">{{.Total}}</div>
            </div>
            <div class="bg-slate-900 rounded-3xl p-8">
                <div class="text-slate-400">Активных</div>
                <div class="text-6xl font-bold mt-2 text-orange-400">{{.Active}}</div>
            </div>
            <div class="bg-slate-900 rounded-3xl p-8">
                <div class="text-slate-400">Бюджет</div>
                <div class="text-6xl font-bold mt-2">{{.Budget}} ₽</div>
            </div>
            <div class="bg-slate-900 rounded-3xl p-8">
                <div class="text-slate-400">Потрачено</div>
                <div class="text-6xl font-bold mt-2 text-emerald-400">{{.Spent}} ₽</div>
                <div class="h-2 bg-slate-700 rounded-full mt-4">
                    <div class="h-2 bg-emerald-500 rounded-full" style="width: {{.SpentPercent}}%"></div>
                </div>
            </div>
        </div>

        <div class="mt-12 text-center text-3xl">
            Средний прогресс: <span class="text-orange-400 font-semibold">{{.AvgProgress}}%</span>
        </div>
    </div>
</body>
</html>`

func projectsHandler(w http.ResponseWriter, r *http.Request) {
	rows, _ := db.Query("SELECT id, name, status, progress, budget, spent FROM projects")
	var projects []Project
	for rows.Next() {
		var p Project
		rows.Scan(&p.ID, &p.Name, &p.Status, &p.Progress, &p.Budget, &p.Spent)
		projects = append(projects, p)
	}
	rows.Close()

	tmpl := template.Must(template.New("projects").Parse(projectsHTML))
	tmpl.Execute(w, projects)
}

const projectsHTML = `<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Проекты</title>
    <script src="https://cdn.tailwindcss.com"></script>
</head>
<body class="bg-slate-950 text-slate-100 p-8">
    <a href="/" class="text-orange-400 mb-8 inline-block">← На дашборд</a>
    <h1 class="text-4xl font-semibold mb-8">Проекты</h1>
    <div class="space-y-6">
        {{range .}}
        <div class="bg-slate-900 rounded-3xl p-6 flex justify-between items-center">
            <div>
                <h3 class="text-2xl">{{.Name}}</h3>
                <div class="h-2 bg-slate-700 rounded-full mt-3 w-96">
                    <div class="h-2 bg-orange-500 rounded-full" style="width: {{.Progress}}%"></div>
                </div>
            </div>
            <div class="text-right">
                <div class="text-5xl font-bold text-orange-400">{{.Progress}}%</div>
                <div class="text-sm text-slate-400">{{.Budget}} ₽</div>
            </div>
            <a href="/project/{{.ID}}" class="bg-orange-500 px-8 py-4 rounded-2xl hover:bg-orange-600">Открыть</a>
        </div>
        {{end}}
    </div>
</body>
</html>`

func projectHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `<h1 style="padding:100px;font-size:30px">Детали проекта (будет позже)</h1>`)
}

func main() {
	initDB()
	defer db.Close()

	http.HandleFunc("/", dashboardHandler)
	http.HandleFunc("/projects", projectsHandler)
	http.HandleFunc("/project/", projectHandler)

	fmt.Println("🚀 СтроМенеджер запущен на :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}