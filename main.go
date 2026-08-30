package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Event struct {
	Message string `json:"message"`
	Number  int    `json:"number"`
}

func (e Event) String() string {
	b, _ := json.MarshalIndent(e, "", "  ")
	return string(b)
}

// Function that returns nothing
func healthCheck(w http.ResponseWriter, r *http.Request) { // w is where the response is written, r the incoming request
	event := Event{Message: "Hello there", Number: 43}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	//w.Write([]byte(`{"status":"ok"}`)) // Write deals in bytes, string is converted to bytes here
	fmt.Println(event)
	if event.Message == "" {
		http.Error(w, "message required", http.StatusBadRequest)
		return
	}
	err := json.NewEncoder(w).Encode(event)
	if err != nil {
		return
	}
}

func main() {
	mux := http.NewServeMux()                   // registry of routes
	mux.HandleFunc("GET /healthz", healthCheck) // register the route, map path to function
	log.Println("Beacon listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux)) // ListenAndServe exposes the port
}
