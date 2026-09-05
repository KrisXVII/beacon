package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

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
