package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

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

type Task struct {
	ID         int
	ProjectID  int
	Name       string
	DueDate    string
	Status     string
	AssignedTo string
	Progress   int
}

// Инициализация БД + тестовые данные
func initDB() {
	var err error
	dbPath := "./db.sqlite"
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}

	// Таблицы
	_, _ = db.Exec(`CREATE TABLE IF NOT EXISTS projects (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT, description TEXT,
		start_date TEXT, end_date TEXT,
		budget REAL, spent REAL,
		status TEXT, progress INTEGER
	)`)
	_, _ = db.Exec(`CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		project_id INTEGER,
		name TEXT, due_date TEXT,
		status TEXT, assigned_to TEXT,
		progress INTEGER,
		FOREIGN KEY(project_id) REFERENCES projects(id)
	)`)

	// Засеиваем данные, если пусто
	count := 0
	db.QueryRow("SELECT COUNT(*) FROM projects").Scan(&count)
	if count == 0 {
		db.Exec(`INSERT INTO projects (name, description, start_date, end_date, budget, spent, status, progress) 
		VALUES 
		('Проектирование ТЦ "Север"', 'Разработка ПД и РД для торгового центра', '2025-01-15', '2025-06-30', 8500000, 3400000, 'in_progress', 45),
		('Строительство ЖК "Лесной"', 'Возведение 12-этажного жилого комплекса', '2025-03-01', '2026-02-28', 24500000, 8200000, 'in_progress', 35),
		('Реконструкция склада', 'Модернизация логистического комплекса', '2024-11-01', '2025-04-15', 4200000, 4100000, 'completed', 100)`)
		
		db.Exec(`INSERT INTO tasks (project_id, name, due_date, status, assigned_to, progress) 
		VALUES 
		(1, 'Геологические изыскания', '2025-02-10', 'done', 'Иванов И.', 100),
		(1, 'Согласование генплана', '2025-04-05', 'in_progress', 'Петров С.', 60),
		(2, 'Фундамент', '2025-05-20', 'in_progress', 'Сидоров А.', 70)`)
	}
}

// Главная страница — красивый дашборд
func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	var totalProj, active int
	var totalBudget, totalSpent float64
	db.QueryRow("SELECT COUNT(*), SUM(budget), SUM(spent) FROM projects").Scan(&totalProj, &totalBudget, &totalSpent)
	db.QueryRow("SELECT COUNT(*) FROM projects WHERE status != 'completed'").Scan(&active)

	avgProgress := 0
	db.QueryRow("SELECT AVG(progress) FROM projects").Scan(&avgProgress)

	tmpl := template.Must(template.New("dashboard").Parse(dashboardHTML))
	data := map[string]interface{}{
		"TotalProjects": totalProj,
		"ActiveProjects": active,
		"TotalBudget":   fmt.Sprintf("%.0f", totalBudget),
		"TotalSpent":    fmt.Sprintf("%.0f", totalSpent),
		"SpentPercent":  int((totalSpent/totalBudget)*100 + 0.5),
		"AvgProgress":   avgProgress,
	}
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
        @import url('https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600&amp;family=Space+Grotesk:wght@500;600&display=swap');
        body { font-family: 'Inter', sans-serif; }
        .title { font-family: 'Space Grotesk', sans-serif; }
    </style>
</head>
<body class="bg-slate-950 text-slate-100">
    <div class="max-w-7xl mx-auto p-8">
        <div class="flex justify-between items-center mb-10">
            <div>
                <h1 class="text-4xl font-semibold title tracking-tight">СтроМенеджер</h1>
                <p class="text-slate-400">Управление проектированием и строительством</p>
            </div>
            <a href="/projects" class="bg-orange-500 hover:bg-orange-600 px-6 py-3 rounded-2xl font-medium flex items-center gap-2">
                <i class="fas fa-list"></i> Все проекты
            </a>
        </div>

        <!-- Цифры по проекту -->
        <div class="grid grid-cols-2 md:grid-cols-4 gap-6 mb-12">
            <div class="bg-slate-900 rounded-3xl p-6">
                <div class="text-slate-400 text-sm">Всего проектов</div>
                <div class="text-5xl font-semibold mt-2">{{.TotalProjects}}</div>
            </div>
            <div class="bg-slate-900 rounded-3xl p-6">
                <div class="text-slate-400 text-sm">Активных</div>
                <div class="text-5xl font-semibold mt-2 text-orange-400">{{.ActiveProjects}}</div>
            </div>
            <div class="bg-slate-900 rounded-3xl p-6">
                <div class="text-slate-400 text-sm">Общий бюджет</div>
                <div class="text-5xl font-semibold mt-2">{{.TotalBudget}} ₽</div>
            </div>
            <div class="bg-slate-900 rounded-3xl p-6">
                <div class="text-slate-400 text-sm">Потрачено</div>
                <div class="text-5xl font-semibold mt-2 text-emerald-400">{{.TotalSpent}} ₽</div>
                <div class="h-2 bg-slate-700 rounded-full mt-4 overflow-hidden">
                    <div class="h-2 bg-emerald-500 rounded-full" style="width: {{.SpentPercent}}%"></div>
                </div>
                <div class="text-xs text-slate-400 mt-1">{{.SpentPercent}}% от бюджета</div>
            </div>
        </div>

        <h2 class="text-2xl font-medium mb-6">Средний прогресс по всем проектам: <span class="text-orange-400">{{.AvgProgress}}%</span></h2>

        <a href="/projects" class="block bg-white text-slate-900 hover:bg-orange-400 hover:text-white transition-colors text-center py-4 rounded-3xl font-semibold text-lg">
            Перейти к списку проектов →
        </a>
    </div>
</body>
</html>`

// Список проектов
func projectsHandler(w http.ResponseWriter, r *http.Request) {
	rows, _ := db.Query("SELECT id, name, status, progress, budget, spent, end_date FROM projects")
	var projects []Project
	for rows.Next() {
		var p Project
		rows.Scan(&p.ID, &p.Name, &p.Status, &p.Progress, &p.Budget, &p.Spent, &p.EndDate)
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
        <a href="/" class="inline-flex items-center gap-2 text-orange-400 mb-8"><i class="fas fa-arrow-left"></i> На дашборд</a>
        <h1 class="text-4xl font-semibold mb-8">Все проекты</h1>
        
        <div class="grid gap-6">
            {{range .}}
            <div class="bg-slate-900 rounded-3xl p-6 flex items-center gap-6">
                <div class="flex-1">
                    <div class="flex justify-between">
                        <h3 class="text-xl font-medium">{{.Name}}</h3>
                        <span class="px-4 py-1 rounded-2xl text-sm bg-slate-800">{{.Status}}</span>
                    </div>
                    <div class="mt-4 h-3 bg-slate-700 rounded-3xl overflow-hidden">
                        <div class="h-3 bg-orange-500 rounded-3xl" style="width: {{.Progress}}%"></div>
                    </div>
                </div>
                <div class="text-right">
                    <div class="text-3xl font-semibold">{{.Progress}}%</div>
                    <div class="text-sm text-slate-400">{{.Budget}} ₽</div>
                </div>
                <a href="/project/{{.ID}}" class="ml-6 bg-orange-500 hover:bg-orange-600 px-8 py-4 rounded-2xl">Открыть</a>
            </div>
            {{end}}
        </div>
    </div>
</body>
</html>`

// Заглушка для /project/{id} (можно расширить позже)
func projectHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `<h1 class="p-8 text-3xl">Детали проекта (расширим в следующей итерации 😉)</h1>`)
}

func main() {
	initDB()
	defer db.Close()

	http.HandleFunc("/", dashboardHandler)
	http.HandleFunc("/projects", projectsHandler)
	http.HandleFunc("/project/", projectHandler)

	fmt.Println("🚀 СтройМенеджер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}