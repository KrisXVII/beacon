package alert

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

type Frame struct {
	File string `json:"file"`
	Line int    `json:"line"`
	Func string `json:"func"`
}

// Define the struct which encodes the shape of the object this API works with and requires from Flask

type Event struct {
	PosthogEventID string    `json:"posthog_event_id" validate:"required"`
	Service        string    `json:"service"`
	ErrorType      string    `json:"error_type" validate:"required"`
	Method         string    `json:"method"`
	Route          string    `json:"route" validate:"required"`
	Environment    string    `json:"environment" validate:"required"`
	Message        string    `json:"message" validate:"required"`
	Frames         []Frame   `json:"frames"`
	OccurredAt     time.Time `json:"occurred_at"`
}

var validate = newValidator()

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
