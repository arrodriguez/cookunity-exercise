package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Define routes
	r.POST("/traces", handleTraces)
	r.GET("/statistics", handleStatistics)

	// Run the server
	r.Run(":8080")
}
