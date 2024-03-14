package main

import (
	"github.com/joho/godotenv"
	"movie-website-backend-core/internal/pkg/database"
	"movie-website-backend-core/internal/pkg/server"
)

const (
	apiKey = "0c4e1c680b2f083a493bfababe6ee4f2"
	apiURL = "https://api.themoviedb.org/3"
)

type Movie struct {
	Title    string  `json:"title"`
	Overview string  `json:"overview"`
	Genres   []Genre `json:"genres"`
	Credits  Credits `json:"credits"`
}

type Genre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Credits struct {
	Cast []Cast `json:"cast"`
	Crew []Crew `json:"crew"`
}

type Cast struct {
	Name      string `json:"name"`
	Character string `json:"character"`
}

type Crew struct {
	Name       string `json:"name"`
	Job        string `json:"job"`
	Department string `json:"department"`
}

func main() {
	if err := godotenv.Load("settings.env"); err != nil {
		panic("Error loading .env file")
	}

	database.ConnectDB()

	server.Start()
}

type MovieDB struct {
	ID       uint      `gorm:"primaryKey"`
	Title    string    `gorm:"not null"`
	Overview string    `gorm:"type:text"`
	Genres   []GenreDB `gorm:"many2many:movie_genres;"`
	Cast     []CastDB  `gorm:"many2many:movie_cast;"`
	Crew     []CrewDB  `gorm:"many2many:movie_crew;"`
}

type GenreDB struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"not null"`
}

type CastDB struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"not null"`
	Character string
}

type CrewDB struct {
	ID         uint   `gorm:"primaryKey"`
	Name       string `gorm:"not null"`
	Job        string
	Department string
}

func ConvertToMovieDB(movie Movie) MovieDB {
	var genresDB []GenreDB
	for _, genre := range movie.Genres {
		genresDB = append(genresDB, GenreDB{ID: uint(genre.ID), Name: genre.Name})
	}

	var castDB []CastDB
	for _, cast := range movie.Credits.Cast {
		castDB = append(castDB, CastDB{Name: cast.Name, Character: cast.Character})
	}

	var crewDB []CrewDB
	for _, crew := range movie.Credits.Crew {
		crewDB = append(crewDB, CrewDB{Name: crew.Name, Job: crew.Job, Department: crew.Department})
	}

	movieDB := MovieDB{
		Title:    movie.Title,
		Overview: movie.Overview,
		Genres:   genresDB,
		Cast:     castDB,
		Crew:     crewDB,
	}

	return movieDB
}
