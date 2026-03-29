package main

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

// ====================== TEMPLATE ======================
func parseTemplates(files ...string) *template.Template {
	paths := make([]string, len(files))
	for i, file := range files {
		paths[i] = filepath.Join("templates", file)
	}
	return template.Must(template.ParseFiles(paths...))
}

// ====================== DASHBOARD ======================
func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := parseTemplates("layout.html", "dashboard.html")
	tmpl.ExecuteTemplate(w, "layout", nil)
}

// ====================== PROJECTS ======================
func projectsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		name := r.FormValue("name")
		description := r.FormValue("description")

		db.Exec(`INSERT INTO projects (name, description) VALUES (?, ?)`,
			name, description)

		http.Redirect(w, r, "/projects", http.StatusSeeOther)
		return
	}

	rows, _ := db.Query(`SELECT id, name, description FROM projects`)
	defer rows.Close()

	var projects []Project
	for rows.Next() {
		var p Project
		rows.Scan(&p.ID, &p.Name, &p.Description)
		projects = append(projects, p)
	}

	tmpl := parseTemplates("layout.html", "projects.html")
	tmpl.ExecuteTemplate(w, "layout", projects)
}

// ====================== PROJECT DETAIL ======================
func projectDetailHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/project/")
	projectID, _ := strconv.Atoi(idStr)

	var project Project
	db.QueryRow(`SELECT id, name FROM projects WHERE id=?`, projectID).
		Scan(&project.ID, &project.Name)

	objects := getObjectsByProject(projectID)

	type StageView struct {
		Stage Stage
		Works []Work
	}

	type ObjectView struct {
		Object Object
		Stages []StageView
	}

	var result []ObjectView

	for _, obj := range objects {
		stages := getStages(obj.ID)

		var stageViews []StageView

		for _, st := range stages {
			works := getWorks(st.ID)

			stageViews = append(stageViews, StageView{
				Stage: st,
				Works: works,
			})
		}

		result = append(result, ObjectView{
			Object: obj,
			Stages: stageViews,
		})
	}

	data := map[string]interface{}{
		"Project": project,
		"Objects": result,
	}

	tmpl := parseTemplates("layout.html", "project-detail.html")
	tmpl.ExecuteTemplate(w, "layout", data)
}

// ====================== OBJECTS ======================
func addObjectHandler(w http.ResponseWriter, r *http.Request) {
	projectID, _ := strconv.Atoi(r.FormValue("project_id"))
	name := r.FormValue("name")

	res, _ := db.Exec(`INSERT INTO objects (project_id, name) VALUES (?, ?)`,
		projectID, name)

	objectID, _ := res.LastInsertId()

	// создаем стадии
	db.Exec(`INSERT INTO stages (object_id, type, name) VALUES (?, 'design', 'Проектирование')`, objectID)
	db.Exec(`INSERT INTO stages (object_id, type, name) VALUES (?, 'construction', 'Строительство')`, objectID)

	http.Redirect(w, r, fmt.Sprintf("/project/%d", projectID), http.StatusSeeOther)
}

func deleteObjectHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/object/delete/")
	id, _ := strconv.Atoi(idStr)

	var projectID int
	db.QueryRow("SELECT project_id FROM objects WHERE id=?", id).Scan(&projectID)

	db.Exec("DELETE FROM objects WHERE id=?", id)

	http.Redirect(w, r, fmt.Sprintf("/project/%d", projectID), http.StatusSeeOther)
}

// ====================== STAGES ======================
func getStages(objectID int) []Stage {
	rows, _ := db.Query(`SELECT id, object_id, type, name, budget, spent, progress 
		FROM stages WHERE object_id=?`, objectID)

	var stages []Stage
	for rows.Next() {
		var s Stage
		rows.Scan(&s.ID, &s.ObjectID, &s.Type, &s.Name, &s.Budget, &s.Spent, &s.Progress)
		stages = append(stages, s)
	}
	return stages
}

// ====================== WORKS ======================
func addWorkHandler(w http.ResponseWriter, r *http.Request) {
	stageID, _ := strconv.Atoi(r.FormValue("stage_id"))
	name := r.FormValue("name")

	db.Exec(`INSERT INTO works (stage_id, name, status) VALUES (?, ?, 'in_progress')`,
		stageID, name)

	http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
}

func deleteWorkHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/work/delete/")
	id, _ := strconv.Atoi(idStr)

	db.Exec("DELETE FROM works WHERE id=?", id)

	http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
}

func getWorks(stageID int) []Work {
	rows, _ := db.Query(`SELECT id, stage_id, name, start_date, end_date, assigned_to, estimated, spent, progress, status 
		FROM works WHERE stage_id=?`, stageID)

	var works []Work
	for rows.Next() {
		var w Work
		rows.Scan(&w.ID, &w.StageID, &w.Name, &w.StartDate, &w.EndDate,
			&w.AssignedTo, &w.Estimated, &w.Spent, &w.Progress, &w.Status)
		works = append(works, w)
	}
	return works
}

// ====================== HELPERS ======================
func getObjectsByProject(projectID int) []Object {
	rows, _ := db.Query(`SELECT id, project_id, name FROM objects WHERE project_id=?`, projectID)

	var objects []Object
	for rows.Next() {
		var o Object
		rows.Scan(&o.ID, &o.ProjectID, &o.Name)
		objects = append(objects, o)
	}
	return objects
}