package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

// ====================== ДАШБОРД ======================
func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	data := getDashboardData()

	tmpl := template.Must(template.New("dashboard").Parse(dashboardHTML))
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// ====================== СПИСОК ПРОЕКТОВ ======================
func projectsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Добавление нового проекта
		name := r.FormValue("name")
		description := r.FormValue("description")
		startDate := r.FormValue("start_date")
		endDate := r.FormValue("end_date")
		budget, _ := strconv.ParseFloat(r.FormValue("budget"), 64)

		_, err := db.Exec(`
			INSERT INTO projects (name, description, start_date, end_date, budget)
			VALUES (?, ?, ?, ?, ?)`, name, description, startDate, endDate, budget)
		if err != nil {
			http.Error(w, "Ошибка при создании проекта", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/projects", http.StatusSeeOther)
		return
	}

	// GET — показ списка
	projects := getAllProjects()

	tmpl := template.Must(template.New("projects").Parse(projectsHTML))
	tmpl.Execute(w, projects)
}

// ====================== ДЕТАЛЬНАЯ СТРАНИЦА ПРОЕКТА ======================
func projectDetailHandler(w http.ResponseWriter, r *http.Request) {
	// Извлекаем ID из URL: /project/5
	path := strings.TrimPrefix(r.URL.Path, "/project/")
	id, err := strconv.Atoi(path)
	if err != nil {
		http.Error(w, "Неверный ID проекта", http.StatusBadRequest)
		return
	}

	var project Project
	err = db.QueryRow(`
		SELECT id, name, description, start_date, end_date, budget, spent, status, progress 
		FROM projects WHERE id = ?`, id).Scan(
		&project.ID, &project.Name, &project.Description,
		&project.StartDate, &project.EndDate, &project.Budget,
		&project.Spent, &project.Status, &project.Progress)

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

	tmpl := template.Must(template.New("projectDetail").Parse(projectDetailHTML))
	tmpl.Execute(w, data)
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

	_, err := db.Exec(`
		INSERT INTO objects (project_id, name, type, area, budget, spent, progress)
		VALUES (?, ?, ?, ?, ?, 0, 0)`, projectID, name, objType, area, budget)

	if err != nil {
		http.Error(w, "Ошибка добавления объекта", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/project/%d", projectID), http.StatusSeeOther)
}

func deleteObjectHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/object/delete/")
	id, _ := strconv.Atoi(idStr)

	var projectID int
	db.QueryRow("SELECT project_id FROM objects WHERE id = ?", id).Scan(&projectID)

	db.Exec("DELETE FROM objects WHERE id = ?", id)

	http.Redirect(w, r, fmt.Sprintf("/project/%d", projectID), http.StatusSeeOther)
}

// ====================== CRUD ЗАДАЧ (ГРАФИК РАБОТ) ======================
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

	_, err := db.Exec(`
		INSERT INTO tasks (project_id, name, start_date, end_date, assigned_to, estimated, spent, progress, status)
		VALUES (?, ?, ?, ?, ?, ?, 0, 0, 'in_progress')`,
		projectID, name, startDate, endDate, assignedTo, estimated)

	if err != nil {
		http.Error(w, "Ошибка добавления задачи", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/project/%d", projectID), http.StatusSeeOther)
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/task/delete/")
	id, _ := strconv.Atoi(idStr)

	var projectID int
	db.QueryRow("SELECT project_id FROM tasks WHERE id = ?", id).Scan(&projectID)

	db.Exec("DELETE FROM tasks WHERE id = ?", id)

	http.Redirect(w, r, fmt.Sprintf("/project/%d", projectID), http.StatusSeeOther)
}

// ====================== ВСПОМОГАТЕЛЬНЫЕ ФУНКЦИИ ======================

func getAllProjects() []Project {
	rows, err := db.Query(`
		SELECT id, name, description, start_date, end_date, budget, spent, status, progress 
		FROM projects`)
	if err != nil {
		return nil
	}
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

func getDashboardData() map[string]interface{} {
	var total, active int
	var totalBudget, totalSpent float64
	var avgProgress int

	db.QueryRow("SELECT COUNT(*), SUM(budget), SUM(spent) FROM projects").Scan(&total, &totalBudget, &totalSpent)
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

func getObjectsByProject(projectID int) []Object {
	rows, _ := db.Query(`
		SELECT id, project_id, name, type, area, budget, spent, progress 
		FROM objects WHERE project_id = ?`, projectID)
	defer rows.Close()

	var objects []Object
	for rows.Next() {
		var o Object
		rows.Scan(&o.ID, &o.ProjectID, &o.Name, &o.Type, &o.Area, &o.Budget, &o.Spent, &o.Progress)
		objects = append(objects, o)
	}
	return objects
}

func getTasksByProject(projectID int) []Task {
	rows, _ := db.Query(`
		SELECT id, project_id, name, start_date, end_date, assigned_to, estimated, spent, progress, status 
		FROM tasks WHERE project_id = ?`, projectID)
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