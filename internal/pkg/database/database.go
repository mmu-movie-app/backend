package database

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"movie-website-backend-core/internal/pkg/config"
	movie_db "movie-website-backend-core/internal/pkg/movie-db"
	"strconv"
)

// DB gorm connector
var DB *gorm.DB

// ConnectDB connect to db
func ConnectDB() {
	var err error
	p := config.Config("DB_PORT")
	port, err := strconv.ParseUint(p, 10, 32)
	if err != nil {
		panic(err)
	}
	println(fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Config("DB_HOST"),
		port,
		config.Config("DB_USER"),
		config.Config("DB_PASSWORD"),
		config.Config("DB_NAME"),
	))
	DB, err = gorm.Open("postgres", fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Config("DB_HOST"),
		port,
		config.Config("DB_USER"),
		config.Config("DB_PASSWORD"),
		config.Config("DB_NAME"),
	))
	if err != nil {
		panic("failed to connect database")
	}

	fmt.Println("Connection Opened to Database")
	DB.AutoMigrate(&movie_db.MovieDB{})
	DB.AutoMigrate(&movie_db.CrewDB{})
	DB.AutoMigrate(&movie_db.GenreDB{})
	DB.AutoMigrate(&movie_db.CastDB{})
	DB.AutoMigrate(&movie_db.ImageDB{})
	fmt.Println("Database Migrated")
}
