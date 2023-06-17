package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.ForwardedByClientIP = true
	r.SetTrustedProxies([]string{"127.0.0.1"})

	v1 := r.Group("/api/v1")
	{
		v1.GET("plant", getPlants)
		v1.GET("plant/:id", getPlantByID)
		v1.POST("plant", addPlant)
		v1.PUT("plant/:id", updatePlant)
		v1.OPTIONS("plant", options)
	}

	r.Run() // listen and serve on 0.0.0.0:8080
}

// getPlants returns a short version of all plants
func getPlants(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "getPlants called"})
}

// getPlantByID returns details of a specific plant based on the id
func getPlantByID(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"message": "getPlantByID " + id + " called"})
}

// addPlant creates a new plant
func addPlant(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "addPlant called"})
}

// updatePlant updates a specific plant based on the id
func updatePlant(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "updatePlant called"})
}

// options answers requests regarding permitted communications
func options(c *gin.Context) {
	ourOptions := "HTTP/1.1 200 OK\n" +
		"Allow: GET, POST, PUT, OPTIONS\n" +
		"Access-Control-Allow-Origin: http://locahost:8080\n" +
		"Access-Control-Allow-Methods: GET, POST, PUT, OPTIONS\n" +
		"Access-Control-Allow-Headers: Content-Type\n"

	c.String(http.StatusOK, ourOptions)
}
