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
	ID        int
	ProjectID int
	Name      string
	Type      string
	Area      float64
	Floors    int
	Material  string
	Budget    float64
	Spent     float64
	Status    string
}

type Stage struct {
	ID       int
	ObjectID int
	Type     string // design / construction
	Name     string

	Budget   float64
	Spent    float64
	Progress int
}

type Work struct {
	ID      int
	StageID int

	Name string

	StartDate string
	EndDate   string

	AssignedTo string

	Estimated float64
	Spent     float64

	Progress int
	Status   string
}