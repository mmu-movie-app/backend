package route

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Content(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"status":  false,
		"message": "invalid movieId",
	})
}
