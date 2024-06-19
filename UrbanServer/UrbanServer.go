package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"math/rand"
	"time"
	"fmt"
    "strconv"
)

// Create struct for GET queries
type GetData struct {
	Latitud	    float64 `json:"lat"`
	Longitud    float64 `json:"lng"`
	Timestamp   string  `json:"timestamp"`
	Value       float64 `json:"value"`
}

// Create struct for data generation for testing purposes
type DebugData struct {
	Latitud     float64 `json:"lat"`
	Longitud    float64 `json:"lng"`
	Timestamp   string  `json:"timestamp"`
	Temperature float64 `json:"temp"`
	Humidity    float64 `json:"humidity"`
	AirQuality  float64 `json:"air"`
	Noise       float64 `json:"noise"`
}

func dataGenerator() []DebugData {
	var data []DebugData

	rand.Seed(time.Now().UnixNano())

	// Valencia range
	minLatitud := 39.40
	maxLatitud := 39.50
	minLongitud := -0.45
	maxLongitud := -0.32

	for i := 0; i < 1000; i++ {
		latitud := rand.Float64()*(maxLatitud-minLatitud) + minLatitud
		longitud := rand.Float64()*(maxLongitud-minLongitud) + minLongitud

		// Generate random date from two months ago to now
		timestamp := time.Now().Add(-time.Duration(rand.Intn(60*24*60*60)) * time.Second).Format("2006-01-02 15:04:05")

		// Generate random temperature, humidity, air quality and noise values
		temp := rand.Float64() * 40.0 // Between 0 and 40 Celsius
		humidity := rand.Float64() * 100.0    // Between 0 and 100%
		airQuality := rand.Float64() * 100.0 // Between 0 and 100
		noise := rand.Float64() * 100.0      // Between 0 and 100 db

		// Create data
		dato := DebugData{
			Latitud:     latitud,
			Longitud:    longitud,
			Timestamp:   timestamp,
			Temperature: temp,
			Humidity:    humidity,
			AirQuality:  airQuality,
			Noise:       noise,
		}
		data = append(data, dato)
	}

	return datos
}

func main() {
	// Initialize database
	db, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Creating table Data
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS Data (" +
		"id INTEGER PRIMARY KEY AUTOINCREMENT," +
		"Latitud FLOAT DEFAULT 0," +
		"Longitud FLOAT DEFAULT 0," +
		"Timestamp TIMESTAMP DEFAULT (strftime('%Y-%m-%d %H:%M:%S', 'now', 'localtime'))," +
		"Temperature FLOAT DEFAULT 0," +
		"Humidity FLOAT DEFAULT 0," +
		"AirQuality FLOAT DEFAULT 0," +
		"Noise FLOAT DEFAULT 0)")
	if err != nil {
		log.Fatal("Could not create table Data: ", err)
	}

	// Preparing CRUD statements (only Create and Read)
	createData, err := db.Prepare("INSERT INTO Data(Latitud, Longitud, Temperature, Humidity, AirQuality, Noise) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal("Could not prepare Data Creation: ", err)
	}
	createDebugData, err := db.Prepare("INSERT INTO Data(Latitud, Longitud, Timestamp, Temperature, Humidity, AirQuality, Noise) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal("Could not prepare debug Data Creation: ", err)
	}
	readTempData, err := db.Prepare("SELECT Latitud, Longitud, Timestamp, Temperature FROM Data WHERE Timestamp BETWEEN ? AND ?")
	if err != nil {
		log.Fatal("Could not prepare temp data Read: ", err)
	}
	readHumidityData, err := db.Prepare("SELECT Latitud, Longitud, Timestamp, Humidity FROM Data WHERE Timestamp BETWEEN ? AND ?")
	if err != nil {
		log.Fatal("Could not prepare humidity data Read: ", err)
	}
	readAirData, err := db.Prepare("SELECT Latitud, Longitud, Timestamp, AirQuality FROM Data WHERE Timestamp BETWEEN ? AND ?")
	if err != nil {
		log.Fatal("Could not prepare air quality data Read: ", err)
	}
	readNoiseData, err := db.Prepare("SELECT Latitud, Longitud, Timestamp, Noise FROM Data WHERE Timestamp BETWEEN ? AND ?")
	if err != nil {
		log.Fatal("Could not prepare noise data Read: ", err)
	}

	// Test Data
	data := dataGenerator();
	for _, d := range data {
		_, err := createDebugData.Exec(d.Latitud, d.Longitud, d.Timestamp, d.Temperature, d.Humidity, d.AirQuality, d.Noise)
		if err != nil {
			log.Print("Could not insert Debug Data(", d, err)
		}
	}
	
	// Create and configure Gin instance and its routes
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Next()
	})

	// Define route to handle POST queries
	r.POST("/data", func(c *gin.Context) {
		_, err := createData.Exec(c.Query("Longitud"), c.Query("lat"), c.Query("lng"), c.Query("humidity"), c.Query("air"), c.Query("noise"))
		if err != nil {
			log.Print("Could not create data: ", err)
			c.JSON(http.StatusBadRequest, "error: Could not create data")
		}
		c.JSON(http.StatusOK, gin.H{"mensaje": "Data stored succesfully"})
	})

	// Define route to handle GET queries
	r.GET("/data", func(c *gin.Context) {
		dataType := c.Query("data_type")
		date := c.Query("date")
		hourStr := c.Query("hour")

		hour, err := strconv.Atoi(hourStr)
	
		startHour := fmt.Sprintf("%02d:00:00", hour-1)
		endHour := fmt.Sprintf("%02d:00:00", hour)

		startTimestamp := date + " " + startHour
		endTimestamp := date + " " + endHour
		
		// Reach data from database depending on data type requested
		var rows *sql.Rows
		if (dataType == "temp") { rows, err = readTempData.Query(startTimestamp, endTimestamp) }
		if (dataType == "humidity") { rows, err = readHumidityData.Query(startTimestamp, endTimestamp) }
		if (dataType == "air") { rows, err = readAirData.Query(startTimestamp, endTimestamp) }
		if (dataType == "noise") { rows, err = readNoiseData.Query(startTimestamp, endTimestamp) }

		if err != nil {
			c.JSON(http.StatusBadRequest, "error:"+err.Error())
			return
		}
		defer rows.Close()

		// Structure data in JSON format
		var data []GetData
		for rows.Next() {
			var Longitud float64
			var Latitud float64
			var Timestamp string
			var Value float64

			err = rows.Scan(&Longitud, &Latitud, &Timestamp, &Value)
			if err != nil {
				c.JSON(http.StatusBadRequest, "error:"+err.Error())
				return
			}
			d := GetData{Longitud, Latitud, Timestamp, Value}
			data = append(data, d)
		}

		// Send JSON data
		c.IndentedJSON(http.StatusOK, data)
	})

	r.Run(":8080")
}
