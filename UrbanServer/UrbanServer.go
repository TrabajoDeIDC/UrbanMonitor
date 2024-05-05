package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
)

// Create struct
type Data struct {
	Timestamp   string  `json:"timestamp"`
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
	AirQuality  float64 `json:"air_quality"`
	Noise       float64 `json:"noise"`
}

func main() {
	// Initialize database
	db, err := sql.Open("sqlite3", "test.db")

	// Handling database opening error
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close() // Here we are saying the program to close DB at the end of this code

	// Creating table Data
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS Data (" +
		"id INTEGER PRIMARY KEY AUTOINCREMENT," +
		"Timestamp STRING DEFAULT CURRENT_TIMESTAMP," +
		"Temperature FLOAT DEFAULT 0," +
		"Humidity FLOAT DEFAULT 0," +
		"AirQuality FLOAT DEFAULT 0," +
		"Noise FLOAT DEFAULT 0)")
	if err != nil {
		log.Fatal("Could not create table Data: ", err)
	}

	// Preparing CRUD statements (only Create and Read)
	createData, err := db.Prepare("INSERT INTO Data(Timestamp, Temperature, Humidity, AirQuality, Noise) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal("Could not prepare Data Creation: ", err)
	}
	readData, err := db.Prepare("SELECT Timestamp, Temperature, Humidity, AirQuality, Noise FROM Data")
	if err != nil {
		log.Fatal("Could not prepare Data Read: ", err)
	}

	/*
		// Test Data
		data := []Data{{"10:24:11", 23, 12, 4, 56},
			{"20:56:01", 12, 12, 12, 12},
			{"15:22:22", 68.1212, 13.2, 4, 44.44}}

		for _, d := range data {
			_, err := createData.Exec(d.Timestamp, d.Temperature, d.Humidity, d.AirQuality, d.Noise)
			if err != nil {
				log.Print("Could not insert Data(", d, err)
			}
		}
	*/

	// Create Gin instance
	r := gin.Default()

	// Define route to handle POST queries
	r.POST("/datos", func(c *gin.Context) {
		_, err := createData.Exec(c.Query("Timestamp"), c.Query("Humidity"), c.Query("AirQuality"), c.Query("Noise"))
		if err != nil {
			log.Print("Could not create data: ", err)
			c.JSON(http.StatusBadRequest, "error: Could not create data")
		}
		c.JSON(http.StatusOK, gin.H{"mensaje": "Data stored succesfully"})
	})

	// Define route to handle GET queries
	r.GET("/datos", func(c *gin.Context) {
		rows, err := readData.Query()
		if err != nil {
			c.JSON(http.StatusBadRequest, "error:"+err.Error())
		} else {
			for rows.Next() {
				var Timestamp string
				var Temperature float64
				var Humidity float64
				var AirQuality float64
				var Noise float64
				err = rows.Scan(&Timestamp, &Temperature, &Humidity, &AirQuality, &Noise)
				if err != nil {
					c.JSON(http.StatusBadRequest, "error:"+err.Error())
				}
				d := Data{Timestamp, Temperature, Humidity, AirQuality, Noise}
				c.IndentedJSON(http.StatusOK, d)
			}
		}
	})

	// Run server on localhost, port 8080
	r.Run(":8080")
}
