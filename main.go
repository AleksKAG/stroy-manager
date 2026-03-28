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

	// Главная страница (дашборд)
	http.HandleFunc("/", dashboardHandler)

	// Список всех проектов + создание нового
	http.HandleFunc("/projects", projectsHandler)

	// Детальная страница проекта
	http.HandleFunc("/project/", projectDetailHandler)

	// CRUD для объектов
	http.HandleFunc("/add-object", addObjectHandler)
	http.HandleFunc("/object/edit", editObjectHandler)
	http.HandleFunc("/object/delete/", deleteObjectHandler)

	// CRUD для задач
	http.HandleFunc("/add-task", addTaskHandler)
	http.HandleFunc("/task/delete/", deleteTaskHandler)

	fmt.Println("🚀 СтроМенеджер успешно запущен!")
	fmt.Println("   Открой в браузере → http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}