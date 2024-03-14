package route

import (
	"github.com/gin-gonic/gin"
	"movie-website-backend-core/internal/pkg/database"
	movie_db "movie-website-backend-core/internal/pkg/movie-db"
	"net/http"
	"strconv"
)

func Add(c *gin.Context) {

	id := c.Query("movieId")
	movieID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  false,
			"message": "invalid movieId",
		})
		return
	}
	moviedb := movie_db.Process(movieID)

	database.DB.Save(&moviedb)
	c.JSON(http.StatusOK, moviedb)
}
