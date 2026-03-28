package main

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

type Object struct {
	ID          int
	ProjectID   int
	Name        string
	Type        string
	Area        float64
	Budget      float64
	Spent       float64
	Progress    int
	Floors      int     // новая характеристика
	Material    string  // новая характеристика
	Status      string
}

type Task struct {
	ID         int
	ProjectID  int
	Name       string
	StartDate  string
	EndDate    string
	AssignedTo string
	Estimated  float64
	Spent      float64
	Progress   int
	Status     string
}