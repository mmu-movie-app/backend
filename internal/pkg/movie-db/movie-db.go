package movie_db

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	apiKey = "0c4e1c680b2f083a493bfababe6ee4f2"
	apiURL = "https://api.themoviedb.org/3"
)

type Movie struct {
	Title       string  `json:"title"`
	Overview    string  `json:"overview"`
	Genres      []Genre `json:"genres"`
	Credits     Credits `json:"credits"`
	Images      Images  `json:"images"`
	Budget      int     `json:"budget"`
	Revenue     int     `json:"revenue"`
	ReleaseDate string  `json:"release_date"`
	VoteAverage float64 `json:"vote_average"`
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
	Name        string `json:"name"`
	Character   string `json:"character"`
	ProfilePath string `json:"profile_path"`
}

type Crew struct {
	Name        string `json:"name"`
	Job         string `json:"job"`
	Department  string `json:"department"`
	ProfilePath string `json:"profile_path"`
}

type Images struct {
	Backdrops []Image `json:"backdrops"`
	Posters   []Image `json:"posters"`
}

type Image struct {
	FilePath string `json:"file_path"`
}

func Process(movieID int) MovieDB {
	// Make API request to fetch movie details
	movieURL := fmt.Sprintf("%s/movie/%d?api_key=%s&append_to_response=credits,images", apiURL, movieID, apiKey)
	response, err := http.Get(movieURL)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Parse the JSON response
	var movie Movie
	err = json.Unmarshal(body, &movie)
	if err != nil {
		log.Fatal(err)
	}

	moviedb := ConvertToMovieDB(movie)
	return moviedb
}

type MovieDB struct {
	ID          uint      `gorm:"primaryKey" json:"-"`
	Title       string    `gorm:"not null" json:"title"`
	Overview    string    `gorm:"type:text" json:"overview"`
	Genres      []GenreDB `gorm:"many2many:movie_genres;" json:"genres"`
	Cast        []CastDB  `gorm:"many2many:movie_cast;" json:"cast"`
	Crew        []CrewDB  `gorm:"many2many:movie_crew;" json:"crew"`
	Backdrops   []ImageDB `gorm:"many2many:movie_backdrops;" json:"backdrops"`
	Posters     []ImageDB `gorm:"many2many:movie_posters;" json:"posters"`
	Budget      int       `gorm:"not null" json:"budget"`
	Revenue     int       `gorm:"not null" json:"revenue"`
	ReleaseYear int       `gorm:"not null" json:"release_year"`
	Rating      float64   `gorm:"not null" json:"rating"`
}

type GenreDB struct {
	ID   uint   `gorm:"primaryKey" json:"-"`
	Name string `gorm:"not null" json:"name"`
}

type CastDB struct {
	ID          uint   `gorm:"primaryKey" json:"-"`
	Name        string `gorm:"not null" json:"name"`
	Character   string `json:"character"`
	ProfilePath string `json:"profile_path"`
}

type CrewDB struct {
	ID          uint   `gorm:"primaryKey" json:"-"`
	Name        string `gorm:"not null" json:"name"`
	Job         string `json:"job"`
	Department  string `json:"department"`
	ProfilePath string `json:"profile_path"`
}

type ImageDB struct {
	ID       uint   `gorm:"primaryKey" json:"-"`
	FilePath string `gorm:"not null" json:"file_path"`
}

func ConvertToMovieDB(movie Movie) MovieDB {
	var genresDB []GenreDB
	for _, genre := range movie.Genres {
		genresDB = append(genresDB, GenreDB{ID: uint(genre.ID), Name: genre.Name})
	}

	var castDB []CastDB
	for _, cast := range movie.Credits.Cast {
		castDB = append(castDB, CastDB{Name: cast.Name, Character: cast.Character, ProfilePath: cast.ProfilePath})
	}

	var crewDB []CrewDB
	for _, crew := range movie.Credits.Crew {
		crewDB = append(crewDB, CrewDB{Name: crew.Name, Job: crew.Job, Department: crew.Department, ProfilePath: crew.ProfilePath})
	}

	var backdropsDB []ImageDB
	for _, backdrop := range movie.Images.Backdrops {
		backdropsDB = append(backdropsDB, ImageDB{FilePath: backdrop.FilePath})
	}

	var postersDB []ImageDB
	for _, poster := range movie.Images.Posters {
		postersDB = append(postersDB, ImageDB{FilePath: poster.FilePath})
	}

	releaseYear := 0
	if len(movie.ReleaseDate) >= 4 {
		fmt.Sscanf(movie.ReleaseDate[:4], "%d", &releaseYear)
	}
	rating := movie.VoteAverage / 2

	movieDB := MovieDB{
		Title:       movie.Title,
		Overview:    movie.Overview,
		Genres:      genresDB,
		Cast:        castDB,
		Crew:        crewDB,
		Backdrops:   backdropsDB,
		Posters:     postersDB,
		Budget:      movie.Budget,
		Revenue:     movie.Revenue,
		ReleaseYear: releaseYear,
		Rating:      rating,
	}

	return movieDB
}
