package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang/gddo/httputil/header"
	"io"
	. "made.by.jst10/outfit7/hancock/cmd/custom_errors"
	"net/http"
	"strings"
)

func decodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) *CustomError {
	if r.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
		if value != "application/json" {
			return GetNotSupportedMediaTypeInRequest(value)
		}
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1048576000)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			return GetJsonDecodingError(err, msg)

		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := fmt.Sprintf("Request body contains badly-formed JSON")
			return GetJsonDecodingError(err, msg)

		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			return GetJsonDecodingError(err, msg)

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			return GetJsonDecodingError(err, msg)

		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			return GetJsonDecodingError(err, msg)

		case err.Error() == "http: request body too large":
			msg := "Request body must not be larger than 1MB"
			return GetJsonDecodingError(err, msg)

		default:
			return GetJsonDecodingError(err, "Unknown error")
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		msg := "Request body must only contain a single JSON object"
		return GetJsonDecodingError(err, msg)
	}
	return nil
}
