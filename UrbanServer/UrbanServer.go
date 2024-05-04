package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	// Create struct
	type Data struct {
		Temperature float64 `json:"temperature"`
		Humidity    float64 `json:"humidity"`
		AirQuality  float64 `json:"air_quality"`
		Noise       float64 `json:"noise"`
	}

	// Test Data'
	data := []Data{{23, 12, 4, 56},
		{12, 12, 12, 12},
		{68.1212, 13.2, 4, 44.44}}

	// Crea una instancia de Gin
	r := gin.Default()

	// Define la ruta para manejar las solicitudes POST
	r.POST("/datos", func(c *gin.Context) {
		// Aquí guarda los datos en la base de datos
		// ...
		c.JSON(http.StatusOK, gin.H{"mensaje": "Datos guardados correctamente"})
	})

	// Define la ruta para manejar las solicitudes GET
	r.GET("/datos", func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, data)
	})

	// Ejecuta el servidor en el puerto 8080
	r.Run(":8080")
}
