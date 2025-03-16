package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	r := gin.Default()

	// Endpoint returning JSON
	r.GET("/b", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello from Service B"})
	})

	// Expose Prometheus metrics at /metrics
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	r.Run(":8000") // listen on port 8000
}
