package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"
)

// Struktura danych
type SensorData struct {
	SensorID    string  `json:"sensor_id"`
	Timestamp   string  `json:"timestamp"`
	PM25        float64 `json:"pm25"`
	PM10        float64 `json:"pm10"`
	CO2         int     `json:"co2"`
	Temperature float64 `json:"temperature"`
	Pressure    int     `json:"pressure"`
	Humidity    int     `json:"humidity"`
}

func generateData(sensorID string) SensorData {
	// Prosta symulacja danych
	return SensorData{
		SensorID:    sensorID,
		Timestamp:   time.Now().Format(time.RFC3339),
		PM25:        10.0 + rand.Float64()*50.0,
		PM10:        20.0 + rand.Float64()*60.0,
		CO2:         400 + rand.Intn(400),
		Temperature: -5.0 + rand.Float64()*35.0,
		Pressure:    990 + rand.Intn(40),
		Humidity:    30 + rand.Intn(70),
	}
}

func main() {
	// 1. Konfiguracja ze zmiennych środowiskowych (Environment Variables)
	// Domyślnie localhost, jeśli nie podano inaczej
	endpointURL := os.Getenv("TARGET_URL")
	if endpointURL == "" {
		endpointURL = "http://localhost:8000/api/data"
	}

	sensorID := os.Getenv("SENSOR_ID")
	if sensorID == "" {
		sensorID = "test-sensor-01"
	}

	// Pobieramy interwał (np. "15m", "10s"). Domyślnie 15 minut.
	intervalStr := os.Getenv("INTERVAL")
	if intervalStr == "" {
		intervalStr = "15m"
	}
	interval, err := time.ParseDuration(intervalStr)
	if err != nil {
		fmt.Println("Błąd formatu czasu, ustawiam 15 minut")
		interval = 15 * time.Minute
	}

	fmt.Printf("Start symulatora %s. Cel: %s. Interwał: %s\n", sensorID, endpointURL, interval)

	// Inicjalizacja generatora losowości
	rand.Seed(time.Now().UnixNano())

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Pętla wysyłania
	for {
		data := generateData(sensorID)
		jsonData, _ := json.Marshal(data)

		// Logowanie wysyłki
		fmt.Printf("Wysyłam dane z %s... ", data.Timestamp)

		resp, err := http.Post(endpointURL, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			// Wypisujemy błąd, ale nie przerywamy programu (serwer może chwilowo nie działać)
			fmt.Printf("BŁĄD: %v\n", err)
		} else {
			fmt.Printf("OK (Status: %s)\n", resp.Status)
			resp.Body.Close()
		}

		<-ticker.C
	}
}
