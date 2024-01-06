package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, data interface{}) error {
	maxBytes := 1048576 // max one megabyte in request body
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	if err != nil {
		return err
	}

	// we only allow one entry in the json file
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only have a single JSON value")
	}

	return nil
}

// writeJSON writes aribtrary data out as JSON
func (app *application) writeJSON(w http.ResponseWriter, status int, data interface{}, headers ...http.Header) error {
	out, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for k, v := range headers[0] {
			w.Header()[k] = v
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(out)

	return nil
}

// badRequest sends a JSON response with status http.StatusBadRequest, describing the error
// func (app *application) badRequest(w http.ResponseWriter, r *http.Request, err error) error {
// 	var payload struct {
// 		Error   bool   `json:"error"`
// 		Message string `json:"message"`
// 	}

// 	payload.Error = true
// 	payload.Message = err.Error()

// 	out, err := json.MarshalIndent(payload, "", "\t")
// 	if err != nil {
// 		return err
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusBadRequest)
// 	w.Write(out)
// 	return nil
// }

// func (app *application) invalidCredentials(w http.ResponseWriter) error {
// 	var payload struct {
// 		Error   bool   `json:"error"`
// 		Message string `json:"message"`
// 	}

// 	payload.Error = true
// 	payload.Message = "invalid authentication credentials"

// 	err := app.writeJSON(w, http.StatusUnauthorized, payload)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// ham lay token trong header
// func (app *application) getTokenOfHeader(r *http.Request) (string, error) {
// 	authorizationHeader := r.Header.Get("Authorization")
// 	if authorizationHeader == "" {
// 		return "", errors.New("no authorization header received")
// 	}

// 	headerParts := strings.Split(authorizationHeader, " ")
// 	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
// 		return "", errors.New("no authorization header received")
// 	}

// 	token := headerParts[1]
// 	if len(token) != 26 {
// 		return "", errors.New("authentication token wrong size")
// 	}

// 	return token, nil
// }
