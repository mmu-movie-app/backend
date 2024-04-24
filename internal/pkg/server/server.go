package server

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"movie-website-backend-core/internal/pkg/server/route"
)

func Start() {
	r := gin.Default()

	// Configure CORS middleware
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}

	r.Use(cors.New(config))

	r.GET("/process", route.Add)
	r.GET("/view/:id", route.ViewFavorite)
	r.POST("/search", route.Search)

	r.Run()
}
