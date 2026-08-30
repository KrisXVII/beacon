package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Event struct {
	Message string `json:"message" validate:"required"`
	Number  int    `json:"number"  validate:"gte=0"`
	//Level   string `json:"level"   validate:"required,oneof=info warn error"`
}

func (e Event) String() string {
	b, _ := json.MarshalIndent(e, "", "  ")
	return string(b)
}

// Validate reports whether the event carries the fields Beacon requires.
func (e Event) Validate() error {
	if e.Message == "" {
		return errors.New(`field "message" is required`)
	}
	if e.Number < 0 {
		return errors.New(`field "number" must be zero or greater`)
	}
	return nil
}

// Function that returns nothing
func healthCheck(w http.ResponseWriter, r *http.Request) { // w is where the response is written, r the incoming request
	//event := Event{Message: "Hello there", Number: 43} // creating a data structure

	//err := w.Write([]byte(`{"status":"ok"}`)) // Write deals in bytes, string is converted to bytes here
	//fmt.Println(event)
	//if event.Message == "" {
	//	http.Error(w, "message required", http.StatusBadRequest)
	//	return
	//}
	//err := json.NewEncoder(w).Encode(event)
	//if err != nil {
	//	return
	//}
	//w.Header().Set("Content-Type", "application/json")
	//w.WriteHeader(http.StatusOK)

	w.Header().Set("Content-Type", "application/json")
	// No WriteHeader: 200 is what we want, and the first write sends it.
	if err := json.NewEncoder(w).Encode(map[string]string{"status": "ok"}); err != nil {
		log.Printf("healthz: encoding response: %v", err)
	}
}

func createEvent(w http.ResponseWriter, r *http.Request) {
	var event Event

	if err := decodeJSON(r, &event); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := event.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//if event.Message == "" { // handwritten field check, explicit validation
	//	http.Error(w, "message is required", http.StatusBadRequest)
	//	return
	//}

	//// declare and test in one line, err exists in the if only
	//if err := json.NewDecoder(r.Body).Decode(&e); err != nil { // &e passes a copy of the event
	//	http.Error(w, "invalid JSON", http.StatusBadRequest)
	//	return
	//}

	fmt.Println(event)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // must precede the body, or Encode sends 200
	if err := json.NewEncoder(w).Encode(event); err != nil {
		log.Printf("createEvent: encoding response: %v", err)
	}

	//// Start encoding, write the event into the response and implicitly write a 200 OK status
	//if err := json.NewEncoder(w).Encode(event); err != nil {
	//	http.Error(w, "Event creation failed", http.StatusInternalServerError)
	//	log.Printf("encoding response: %v", err)
	//}
}

// decodeJSON reads a single JSON object from the request body into dst,
// translating decoder failures into messages safe to return to the client.
func decodeJSON(r *http.Request, dst any) error {

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		return decodeError(err)
	}

	return nil
}

func decodeError(err error) error {
	var typeErr *json.UnmarshalTypeError
	var syntaxErr *json.SyntaxError
	//var maxErr *http.MaxBytesError

	switch {
	case errors.As(err, &typeErr):
		return fmt.Errorf("field %q must be a %s", typeErr.Field, typeErr.Type)
	case errors.As(err, &syntaxErr):
		return fmt.Errorf("malformed JSON at position %d", syntaxErr.Offset)
	case errors.Is(err, io.EOF):
		return errors.New("request body is empty")
	default:
		return errors.New("invalid JSON")
	}
}

func main() {
	mux := http.NewServeMux()                   // registry of routes
	mux.HandleFunc("GET /healthz", healthCheck) // register the route, map path to function
	mux.HandleFunc("POST /events", createEvent)

	log.Println("Beacon listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux)) // ListenAndServe exposes the port
}
