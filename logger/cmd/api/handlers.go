package main

import (
	"net/http"

	"logger/data"
)

type JSONPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (s *Service) WriteLog(w http.ResponseWriter, r *http.Request) {
	// read json
	var reqPayload JSONPayload
	_ = s.readJSON(w, r, &reqPayload)

	// insert data
	logEntry := data.LogEntry{
		Name: reqPayload.Name,
		Data: reqPayload.Data,
	}
	if err := s.Models.LogEntry.Insert(logEntry); err != nil {
		_ = s.errorJSON(w, err)
		return
	}

	resp := jsonResponse{
		Message: "logged",
	}

	_ = s.writeJSON(w, http.StatusAccepted, resp)
}
