package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

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

// ====================== ИНИЦИАЛИЗАЦИЯ БД ======================
func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "./db.sqlite")
	if err != nil {
		log.Fatal(err)
	}

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

	// Добавляем тестовые данные один раз
	var count int
	db.QueryRow("SELECT COUNT(*) FROM projects").Scan(&count)
	if count == 0 {
		db.Exec(`INSERT INTO projects 
			(name, description, start_date, end_date, budget, spent, status, progress) VALUES 
			('Проектирование ТЦ "Север"', 'Разработка проектной и рабочей документации', '2025-01-15', '2025-06-30', 8500000, 3400000, 'in_progress', 45),
			('Строительство ЖК "Лесной"', 'Возведение жилого комплекса на 12 этажей', '2025-03-01', '2026-02-28', 24500000, 8200000, 'in_progress', 35),
			('Реконструкция склада', 'Полная модернизация логистического комплекса', '2024-11-01', '2025-04-15', 4200000, 4100000, 'completed', 100)`)
	}
}

// ====================== ДАШБОРД ======================
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
<body class="bg-slate-950 text-slate-100">
    <div class="max-w-7xl mx-auto p-8">
        <div class="flex justify-between items-center mb-10">
            <div>
                <h1 class="text-5xl font-semibold title">СтроМенеджер</h1>
                <p class="text-slate-400">Управление проектированием и строительством</p>
            </div>
            <a href="/projects" class="bg-orange-500 hover:bg-orange-600 px-6 py-3 rounded-2xl font-medium flex items-center gap-2">
                <i class="fas fa-list"></i> Все проекты
            </a>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-4 gap-6">
            <div class="bg-slate-900 rounded-3xl p-8">
                <div class="text-slate-400">Всего проектов</div>
                <div class="text-6xl font-bold mt-3">{{.Total}}</div>
            </div>
            <div class="bg-slate-900 rounded-3xl p-8">
                <div class="text-slate-400">Активных</div>
                <div class="text-6xl font-bold mt-3 text-orange-400">{{.Active}}</div>
            </div>
            <div class="bg-slate-900 rounded-3xl p-8">
                <div class="text-slate-400">Общий бюджет</div>
                <div class="text-6xl font-bold mt-3">{{.Budget}} ₽</div>
            </div>
            <div class="bg-slate-900 rounded-3xl p-8">
                <div class="text-slate-400">Потрачено</div>
                <div class="text-6xl font-bold mt-3 text-emerald-400">{{.Spent}} ₽</div>
                <div class="h-3 bg-slate-700 rounded-full mt-4 overflow-hidden">
                    <div class="h-3 bg-emerald-500 rounded-full" style="width: {{.SpentPercent}}%"></div>
                </div>
            </div>
        </div>

        <div class="mt-12 text-center text-3xl font-medium">
            Средний прогресс по проектам: <span class="text-orange-400">{{.AvgProgress}}%</span>
        </div>
    </div>
</body>
</html>`

// ====================== СПИСОК + CRUD ======================
func projectsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// Добавление нового проекта
		name := r.FormValue("name")
		desc := r.FormValue("description")
		start := r.FormValue("start_date")
		end := r.FormValue("end_date")
		budget, _ := strconv.ParseFloat(r.FormValue("budget"), 64)

		_, err := db.Exec(`INSERT INTO projects (name, description, start_date, end_date, budget) 
			VALUES (?, ?, ?, ?, ?)`, name, desc, start, end, budget)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		http.Redirect(w, r, "/projects", http.StatusSeeOther)
		return
	}

	// GET — показываем список
	rows, _ := db.Query("SELECT id, name, description, start_date, end_date, budget, spent, status, progress FROM projects")
	var projects []Project
	for rows.Next() {
		var p Project
		rows.Scan(&p.ID, &p.Name, &p.Description, &p.StartDate, &p.EndDate, &p.Budget, &p.Spent, &p.Status, &p.Progress)
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
    <title>Проекты — СтроМенеджер</title>
    <script src="https://cdn.tailwindcss.com"></script>
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.5.1/css/all.min.css">
</head>
<body class="bg-slate-950 text-slate-100">
    <div class="max-w-7xl mx-auto p-8">
        <a href="/" class="text-orange-400 hover:text-orange-300 mb-8 inline-flex items-center gap-2">
            ← На дашборд
        </a>

        <div class="flex justify-between items-center mb-8">
            <h1 class="text-4xl font-semibold">Все проекты</h1>
            <button onclick="document.getElementById('addModal').classList.remove('hidden')" 
                    class="bg-orange-500 hover:bg-orange-600 px-6 py-3 rounded-2xl flex items-center gap-2">
                <i class="fas fa-plus"></i> Новый проект
            </button>
        </div>

        <!-- Список проектов -->
        <div class="space-y-6">
            {{range .}}
            <div class="bg-slate-900 rounded-3xl p-6 flex flex-col md:flex-row md:items-center gap-6">
                <div class="flex-1">
                    <h3 class="text-2xl font-medium">{{.Name}}</h3>
                    <p class="text-slate-400 text-sm mt-1">{{.Description}}</p>
                    <div class="mt-4 h-3 bg-slate-700 rounded-full overflow-hidden">
                        <div class="h-3 bg-orange-500 rounded-full" style="width: {{.Progress}}%"></div>
                    </div>
                </div>
                <div class="text-right md:text-center">
                    <div class="text-4xl font-bold text-orange-400">{{.Progress}}%</div>
                    <div class="text-xs text-slate-400">прогресс</div>
                </div>
                <div class="flex gap-3">
                    <a href="/project/edit/{{.ID}}" class="bg-slate-700 hover:bg-slate-600 px-5 py-3 rounded-2xl">
                        <i class="fas fa-edit"></i>
                    </a>
                    <form action="/project/delete/{{.ID}}" method="POST" onsubmit="return confirm('Удалить проект?')">
                        <button type="submit" class="bg-red-600 hover:bg-red-700 px-5 py-3 rounded-2xl">
                            <i class="fas fa-trash"></i>
                        </button>
                    </form>
                    <a href="/project/{{.ID}}" class="bg-orange-500 hover:bg-orange-600 px-8 py-3 rounded-2xl font-medium">Открыть</a>
                </div>
            </div>
            {{end}}
        </div>

        <!-- Модальное окно добавления проекта -->
        <div id="addModal" class="hidden fixed inset-0 bg-black/70 flex items-center justify-center">
            <div class="bg-slate-900 rounded-3xl p-8 w-full max-w-lg">
                <h2 class="text-2xl font-semibold mb-6">Новый проект</h2>
                <form method="POST" class="space-y-4">
                    <input type="text" name="name" placeholder="Название проекта" required class="w-full bg-slate-800 rounded-2xl px-5 py-4">
                    <textarea name="description" placeholder="Описание" rows="3" class="w-full bg-slate-800 rounded-2xl px-5 py-4"></textarea>
                    <div class="grid grid-cols-2 gap-4">
                        <input type="date" name="start_date" class="bg-slate-800 rounded-2xl px-5 py-4">
                        <input type="date" name="end_date" class="bg-slate-800 rounded-2xl px-5 py-4">
                    </div>
                    <input type="number" name="budget" placeholder="Бюджет (₽)" step="1000" class="w-full bg-slate-800 rounded-2xl px-5 py-4">
                    <div class="flex gap-4">
                        <button type="submit" class="flex-1 bg-orange-500 hover:bg-orange-600 py-4 rounded-2xl font-medium">Создать проект</button>
                        <button type="button" onclick="document.getElementById('addModal').classList.add('hidden')" 
                                class="flex-1 bg-slate-700 hover:bg-slate-600 py-4 rounded-2xl">Отмена</button>
                    </div>
                </form>
            </div>
        </div>
    </div>
</body>
</html>`

// Удаление проекта
func deleteProjectHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/project/delete/"):]
	id, _ := strconv.Atoi(idStr)

	_, _ = db.Exec("DELETE FROM projects WHERE id = ?", id)
	http.Redirect(w, r, "/projects", http.StatusSeeOther)
}

// Заглушка для просмотра и редактирования (пока простая)
func projectHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `<div style="padding:80px; font-size:28px; text-align:center">
		<h1>Страница проекта</h1>
		<p>Здесь позже будут задачи, бюджет по статьям, график и т.д.</p>
		<a href="/projects" style="color:#f97316">← Вернуться к проектам</a>
	</div>`)
}

func main() {
	initDB()
	defer db.Close()

	http.HandleFunc("/", dashboardHandler)
	http.HandleFunc("/projects", projectsHandler)
	http.HandleFunc("/project/delete/", deleteProjectHandler)
	http.HandleFunc("/project/", projectHandler) // для /project/1 и /project/edit/1

	fmt.Println("🚀 СтроМенеджер с CRUD запущен на :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}