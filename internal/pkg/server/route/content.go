package route

import (
	"github.com/gin-gonic/gin"
	"movie-website-backend-core/internal/pkg/database"
	movie_db "movie-website-backend-core/internal/pkg/movie-db"
	"net/http"
	"strconv"
)

func ViewFavorite(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID"})
		return
	}

	var movie movie_db.MovieDB
	if err := database.DB.Model(&movie_db.MovieDB{}).
		Preload("Genres").
		Preload("Cast").
		Preload("Crew").
		Preload("Backdrops").
		Preload("Posters").
		First(&movie, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
		return
	}

	c.JSON(http.StatusOK, movie)
}

func GetMovieByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID"})
		return
	}

	var movie movie_db.MovieDB
	if err := database.DB.Model(&movie_db.MovieDB{}).
		Preload("Genres").
		Preload("Cast").
		Preload("Crew").
		Preload("Backdrops").
		Preload("Posters").
		First(&movie, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
		return
	}

	c.JSON(http.StatusOK, movie)
}
