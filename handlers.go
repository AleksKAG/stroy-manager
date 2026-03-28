package main

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

// ====================== ВСПОМОГАТЕЛЬНАЯ ФУНКЦИЯ ДЛЯ ШАБЛОНОВ ======================
func parseTemplates(files ...string) *template.Template {
	paths := make([]string, len(files))
	for i, file := range files {
		paths[i] = filepath.Join("templates", file)
	}
	return template.Must(template.ParseFiles(paths...))
}

// ====================== ДАШБОРД ======================
func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	data := getDashboardData()
	tmpl := parseTemplates("layout.html", "dashboard.html")
	tmpl.ExecuteTemplate(w, "layout", data)
}

// ====================== СПИСОК ПРОЕКТОВ ======================
func projectsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		name := r.FormValue("name")
		description := r.FormValue("description")
		startDate := r.FormValue("start_date")
		endDate := r.FormValue("end_date")
		budget, _ := strconv.ParseFloat(r.FormValue("budget"), 64)

		_, _ = db.Exec(`INSERT INTO projects (name, description, start_date, end_date, budget) 
			VALUES (?, ?, ?, ?, ?)`, name, description, startDate, endDate, budget)

		http.Redirect(w, r, "/projects", http.StatusSeeOther)
		return
	}

	projects := getAllProjects()
	tmpl := parseTemplates("layout.html", "projects.html")
	tmpl.ExecuteTemplate(w, "layout", projects)
}

// ====================== ДЕТАЛИ ПРОЕКТА ======================
func projectDetailHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/project/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный ID проекта", http.StatusBadRequest)
		return
	}

	var project Project
	err = db.QueryRow(`SELECT id, name, description, start_date, end_date, budget, spent, status, progress 
		FROM projects WHERE id = ?`, id).Scan(
		&project.ID, &project.Name, &project.Description,
		&project.StartDate, &project.EndDate,
		&project.Budget, &project.Spent, &project.Status, &project.Progress)

	if err != nil {
		http.Error(w, "Проект не найден", http.StatusNotFound)
		return
	}

	objects := getObjectsByProject(id)
	tasks := getTasksByProject(id)

	data := map[string]interface{}{
		"Project": project,
		"Objects": objects,
		"Tasks":   tasks,
	}

	tmpl := parseTemplates("layout.html", "project-detail.html")
	tmpl.ExecuteTemplate(w, "layout", data)
}

// ====================== CRUD ОБЪЕКТОВ ======================
func addObjectHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	projectID, _ := strconv.Atoi(r.FormValue("project_id"))
	name := r.FormValue("name")
	objType := r.FormValue("type")
	area, _ := strconv.ParseFloat(r.FormValue("area"), 64)
	budget, _ := strconv.ParseFloat(r.FormValue("budget"), 64)
	floors, _ := strconv.Atoi(r.FormValue("floors"))
	material := r.FormValue("material")

	_, _ = db.Exec(`INSERT INTO objects (project_id, name, type, area, budget, floors, material, status) 
		VALUES (?, ?, ?, ?, ?, ?, ?, 'in_progress')`,
		projectID, name, objType, area, budget, floors, material)

	http.Redirect(w, r, fmt.Sprintf("/project/%d", projectID), http.StatusSeeOther)
}

func editObjectHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	id, _ := strconv.Atoi(r.FormValue("id"))
	name := r.FormValue("name")
	objType := r.FormValue("type")
	area, _ := strconv.ParseFloat(r.FormValue("area"), 64)
	budget, _ := strconv.ParseFloat(r.FormValue("budget"), 64)
	floors, _ := strconv.Atoi(r.FormValue("floors"))
	material := r.FormValue("material")
	status := r.FormValue("status")

	_, _ = db.Exec(`UPDATE objects SET name=?, type=?, area=?, budget=?, floors=?, material=?, status=? WHERE id=?`,
		name, objType, area, budget, floors, material, status, id)

	var projectID int
	db.QueryRow("SELECT project_id FROM objects WHERE id = ?", id).Scan(&projectID)

	http.Redirect(w, r, fmt.Sprintf("/project/%d", projectID), http.StatusSeeOther)
}

func deleteObjectHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/object/delete/")
	id, _ := strconv.Atoi(idStr)

	var projectID int
	db.QueryRow("SELECT project_id FROM objects WHERE id = ?", id).Scan(&projectID)

	_, _ = db.Exec("DELETE FROM objects WHERE id = ?", id)
	http.Redirect(w, r, fmt.Sprintf("/project/%d", projectID), http.StatusSeeOther)
}

// ====================== CRUD ЗАДАЧ ======================
func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	projectID, _ := strconv.Atoi(r.FormValue("project_id"))
	name := r.FormValue("name")
	startDate := r.FormValue("start_date")
	endDate := r.FormValue("end_date")
	assignedTo := r.FormValue("assigned_to")
	estimated, _ := strconv.ParseFloat(r.FormValue("estimated"), 64)

	_, _ = db.Exec(`INSERT INTO tasks (project_id, name, start_date, end_date, assigned_to, estimated, status) 
		VALUES (?, ?, ?, ?, ?, ?, 'in_progress')`,
		projectID, name, startDate, endDate, assignedTo, estimated)

	http.Redirect(w, r, fmt.Sprintf("/project/%d", projectID), http.StatusSeeOther)
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/task/delete/")
	id, _ := strconv.Atoi(idStr)

	var projectID int
	db.QueryRow("SELECT project_id FROM tasks WHERE id = ?", id).Scan(&projectID)

	_, _ = db.Exec("DELETE FROM tasks WHERE id = ?", id)
	http.Redirect(w, r, fmt.Sprintf("/project/%d", projectID), http.StatusSeeOther)
}

// ====================== ВСПОМОГАТЕЛЬНЫЕ ФУНКЦИИ (ДОБАВЛЕНЫ) ======================

func getDashboardData() map[string]interface{} {
	var total, active int
	var totalBudget, totalSpent float64
	var avgProgress int

	db.QueryRow("SELECT COUNT(*), COALESCE(SUM(budget), 0), COALESCE(SUM(spent), 0) FROM projects").
		Scan(&total, &totalBudget, &totalSpent)
	db.QueryRow("SELECT COUNT(*) FROM projects WHERE status != 'completed'").Scan(&active)
	db.QueryRow("SELECT COALESCE(AVG(progress), 0) FROM projects").Scan(&avgProgress)

	spentPercent := 0
	if totalBudget > 0 {
		spentPercent = int((totalSpent / totalBudget) * 100)
	}

	return map[string]interface{}{
		"Total":        total,
		"Active":       active,
		"Budget":       fmt.Sprintf("%.0f", totalBudget),
		"Spent":        fmt.Sprintf("%.0f", totalSpent),
		"SpentPercent": spentPercent,
		"AvgProgress":  avgProgress,
	}
}

func getAllProjects() []Project {
	rows, _ := db.Query(`SELECT id, name, description, start_date, end_date, budget, spent, status, progress 
		FROM projects ORDER BY id DESC`)
	defer rows.Close()

	var projects []Project
	for rows.Next() {
		var p Project
		rows.Scan(&p.ID, &p.Name, &p.Description, &p.StartDate, &p.EndDate,
			&p.Budget, &p.Spent, &p.Status, &p.Progress)
		projects = append(projects, p)
	}
	return projects
}

func getObjectsByProject(projectID int) []Object {
	rows, _ := db.Query(`SELECT id, project_id, name, type, area, budget, spent, progress, floors, material, status 
		FROM objects WHERE project_id = ?`, projectID)
	defer rows.Close()

	var objects []Object
	for rows.Next() {
		var o Object
		rows.Scan(&o.ID, &o.ProjectID, &o.Name, &o.Type, &o.Area, &o.Budget,
			&o.Spent, &o.Progress, &o.Floors, &o.Material, &o.Status)
		objects = append(objects, o)
	}
	return objects
}

func getTasksByProject(projectID int) []Task {
	rows, _ := db.Query(`SELECT id, project_id, name, start_date, end_date, assigned_to, 
		estimated, spent, progress, status FROM tasks WHERE project_id = ?`, projectID)
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		rows.Scan(&t.ID, &t.ProjectID, &t.Name, &t.StartDate, &t.EndDate,
			&t.AssignedTo, &t.Estimated, &t.Spent, &t.Progress, &t.Status)
		tasks = append(tasks, t)
	}
	return tasks
}