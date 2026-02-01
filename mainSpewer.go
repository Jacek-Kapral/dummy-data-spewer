package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type SensorData struct {
	SensorID    string  `json:"sensor_id"`
	Timestamp   string  `json:"timestamp"`
	PM25        int     `json:"pm25"`
	PM10        int     `json:"pm10"`
	CO2         int     `json:"co2"`
	Temperature float64 `json:"temperature"`
	Pressure    int     `json:"pressure"`
	Humidity    int     `json:"humidity"`
}

func roundToTwo(val float64) float64 {
	return math.Round(val*100) / 100
}

func generateData(sensorID string) SensorData {
	now := time.Now()
	formattedTime := fmt.Sprintf("Date: %s Time: %s", now.Format("2006-01-02"), now.Format("15:04:05"))
	return SensorData{
		SensorID:    sensorID,
		Timestamp:   formattedTime,
		PM25:        int(10.0 + rand.Float64()*50.0),
		PM10:        int(20.0 + rand.Float64()*60.0),
		CO2:         400 + rand.Intn(400),
		Temperature: roundToTwo(-5.0 + rand.Float64()*35.0),
		Pressure:    990 + rand.Intn(40),
		Humidity:    30 + rand.Intn(70),
	}
}

func main() {
	endpointURL := os.Getenv("TARGET_URL")
	if endpointURL == "" {
		endpointURL = "http://localhost:8000/api/data"
	}

	sensorID := os.Getenv("SENSOR_ID")
	if sensorID == "" {
		sensorID = "test-sensor-01"
	}

	interval := 30 * time.Second

	fmt.Printf("Simulation start for: %s\nTarget: %s\nInterval: %s\n", sensorID, endpointURL, interval)
	fmt.Println("--------------------------------------------------")

	rand.Seed(time.Now().UnixNano())

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		data := generateData(sensorID)
		jsonData, _ := json.MarshalIndent(data, "", "  ")

		fmt.Printf("[%s] Transmission attempt... ", data.Timestamp)

		resp, err := client.Post(endpointURL, "application/json", bytes.NewBuffer(jsonData))

		if err != nil {
			// If server is unreachable
			fmt.Printf("\n!!! Connection Error: %v\n", err)
			fmt.Println(">>> Dumping data to console:")
			fmt.Println(string(jsonData))
			fmt.Println("--------------------------------------------------")
		} else {
			// If server responds
			fmt.Printf("OK (Status: %s)\n", resp.Status)
			resp.Body.Close()
		}

		<-ticker.C
	}
}
