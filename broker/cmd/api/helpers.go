package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// readJSON tries to read the body of a request and converts it into JSON.
func (s *Service) readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxLength := 1048576 // 1 MB

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxLength))

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(data); err != nil {
		return err
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return errors.New("body must have only a single JSON value")
	}

	return nil
}

// writeJSON tries to write JSON to http.ResponseWriter.
func (s *Service) writeJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil
}

// errorJSON write JSON response struct jsonResponse to http.ResponseWriter.
func (s *Service) errorJSON(w http.ResponseWriter, err error, statuses ...int) error {
	statusCode := http.StatusBadRequest

	if len(statuses) > 0 {
		statusCode = statuses[0]
	}

	payload := jsonResponse{
		Error:   true,
		Message: err.Error(),
	}

	return s.writeJSON(w, statusCode, payload)
}
