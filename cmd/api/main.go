package main

import (
	"devzgabriel/goexpert-lab-deploy/internal/services/cep"
	"devzgabriel/goexpert-lab-deploy/internal/services/weather"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type SuccessResponse struct {
	TempC float32 `json:"temp_c"`
	TempF float32 `json:"temp_f"`
	TempK float32 `json:"temp_k"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func main() {

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	weatherSecretKey := os.Getenv("WEATHER_SECRET_KEY")

	if weatherSecretKey == "" {
		fmt.Println("WEATHER_SECRET_KEY is not set in the environment variables")
		weatherSecretKey = "f0b7cd19ad1841aeb40192648250906"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Hello, World! This is the Go Expert Lab Deploy API!"))
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/execute", func(w http.ResponseWriter, r *http.Request) {

		cepParam := r.URL.Query().Get("cep")
		if cepParam == "" || len(cepParam) != 8 {
			http.Error(w, "CEP is required", http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json")
			response := ErrorResponse{
				Message: "Mensagem: invalid zipcode",
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(response)
			return
		}

		cepResponse, err := cep.GetCepFromViaCep(cepParam)
		if err != nil {
			http.Error(w, "Error fetching CEP data", http.StatusInternalServerError)
			w.Header().Set("Content-Type", "application/json")
			response := ErrorResponse{
				Message: "Mensagem: can not find zipcode",
			}
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(response)
			return
		}
		if cepResponse.Cep == "" {
			http.Error(w, "Invalid CEP", http.StatusNotFound)
			w.Header().Set("Content-Type", "application/json")
			response := ErrorResponse{
				Message: "Mensagem: can not find zipcode",
			}
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(response)
			return
		}

		weatherResponse, err := weather.GetWeatherData(cepResponse.City, weatherSecretKey)
		if err != nil {
			http.Error(w, "Error fetching weather data", http.StatusInternalServerError)
			w.Header().Set("Content-Type", "application/json")
			response := ErrorResponse{
				Message: "Mensagem: can not find weather data",
			}
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(response)
			return
		}

		response := SuccessResponse{
			TempC: weatherResponse.Current.TempC,
			TempF: weatherResponse.Current.TempF,
			TempK: weatherResponse.Current.TempC + 273,
		}

		fmt.Println("Request processed successfully for CEP:", cepParam)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	})

	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}
