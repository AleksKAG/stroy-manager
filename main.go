package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	initDB()           // из db.go
	defer db.Close()

	// ==================== РОУТЫ ====================
	http.HandleFunc("/", dashboardHandler)
	http.HandleFunc("/projects", projectsHandler)
	http.HandleFunc("/project/", projectDetailHandler)

	// CRUD для объектов
	http.HandleFunc("/add-object", addObjectHandler)
	http.HandleFunc("/object/edit", editObjectHandler)    
	http.HandleFunc("/object/delete/", deleteObjectHandler)

	// CRUD для задач
	http.HandleFunc("/add-task", addTaskHandler)
	http.HandleFunc("/task/delete/", deleteTaskHandler)

	fmt.Println("🚀 СтроМенеджер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}