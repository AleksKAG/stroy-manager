package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Инициализация базы данных
	initDB()
	defer db.Close()

	// ==================== РОУТИНГ ====================

	// Главная страница
	http.HandleFunc("/", dashboardHandler)

	// Проекты
	http.HandleFunc("/projects", projectsHandler)
	http.HandleFunc("/project/", projectDetailHandler)

	// Объекты
	http.HandleFunc("/add-object", addObjectHandler)
	http.HandleFunc("/object/delete/", deleteObjectHandler)

	http.HandleFunc("/add-work", addWorkHandler)
	http.HandleFunc("/work/delete/", deleteWorkHandler)

	fmt.Println("🚀 СтройМенеджер запущен!")
	fmt.Println("👉 http://localhost:8080")

	log.Fatal(http.ListenAndServe(":8080", nil))
}