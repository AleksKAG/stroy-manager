package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	initDB()
	defer db.Close()

	http.HandleFunc("/", dashboardHandler)
	http.HandleFunc("/projects", projectsHandler)
	http.HandleFunc("/project/", projectDetailHandler)
	http.HandleFunc("/add-object", addObjectHandler)
	http.HandleFunc("/add-task", addTaskHandler)
	http.HandleFunc("/object/delete/", deleteObjectHandler)
	http.HandleFunc("/task/delete/", deleteTaskHandler)

	fmt.Println("🚀 СтроМенеджер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}