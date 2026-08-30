package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = newValidator()

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
	err := validate.Struct(e)
	if err == nil {
		return nil
	}
	var verrs validator.ValidationErrors
	if errors.As(err, &verrs) && len(verrs) > 0 {
		f := verrs[0]
		switch f.Tag() {
		case "required":
			return fmt.Errorf("field %q is required", f.Field())
		case "gte":
			return fmt.Errorf("field %q must be at least %s", f.Field(), f.Param())
		default:
			return fmt.Errorf("field %q is invalid", f.Field())
		}
	}
	return errors.New("invalid payload")
}

func newValidator() *validator.Validate {
	v := validator.New()
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	return v
}

// healthCheck reports that the service is running.
func healthCheck(w http.ResponseWriter, r *http.Request) { // w is where the response is written, r the incoming request
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

	//fmt.Println(event)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // must precede the body, or Encode sends 200
	if err := json.NewEncoder(w).Encode(event); err != nil {
		log.Printf("createEvent: encoding response: %v", err)
	}
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
