package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

type WeatherData struct {
	Name string `json:"name"`
	Base string `json:"base"`
	Main struct {
		Temp float64 `json:"temp"`
		Feels_Like float64 `json:"feels_like"`
	} `json:"main"`
}

var (
	view    = template.Must(template.ParseFiles("./static/index.html"))
	apiKey  = "5f1bcb7a619ad01056258319bcb94057" // Replace with your OpenWeatherMap API key
	city    = "Bangladesh"
	staticDir, _ = filepath.Abs("./static")
)

func main() {
	mux := http.NewServeMux()

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))))
	mux.HandleFunc("/", index)

	server := http.Server{
		Addr:    ":9000",
		Handler: mux,
	}

	log.Println("Your app is Running...")

	log.Fatal(server.ListenAndServe())
}

func index(w http.ResponseWriter, r *http.Request) {
	url := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric", city, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		log.Println("Error fetching weather data:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var data WeatherData

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		log.Println("Error decoding weather data:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	fmt.Printf("Current temperature in '%s' (Feels_Like '%.1f', Base is '%s') is %.1f degree Celsius\n", data.Name, data.Main.Feels_Like, data.Base, data.Main.Temp)

	err = view.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		log.Println("Error executing template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
