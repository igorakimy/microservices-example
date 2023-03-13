package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

const (
	authURL = "http://auth/authenticate"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s *Service) Broker(w http.ResponseWriter, _ *http.Request) {
	payload := jsonResponse{
		Message: "Hit the broker",
	}

	if err := s.writeJSON(w, http.StatusAccepted, payload); err != nil {
		log.Println(err)
	}
}

func (s *Service) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var reqPayload RequestPayload

	if err := s.readJSON(w, r, &reqPayload); err != nil {
		_ = s.errorJSON(w, err)
		return
	}

	switch reqPayload.Action {
	case "auth":
		s.authenticate(w, reqPayload.Auth)
	default:
		_ = s.errorJSON(w, errors.New("unknown action"))
	}
}

func (s *Service) authenticate(w http.ResponseWriter, ap AuthPayload) {
	// create some json we'll send to the auth microservice
	jsonData, err := json.MarshalIndent(ap, "", "\t")
	if err != nil {
		log.Println(err)
		return
	}
	// call the service
	request, err := http.NewRequest(
		"POST",
		authURL,
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		_ = s.errorJSON(w, err)
		return
	}

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		_ = s.errorJSON(w, err)
		return
	}
	defer resp.Body.Close()

	// make sure we get back the correct status code
	if resp.StatusCode == http.StatusUnauthorized {
		_ = s.errorJSON(w, errors.New("invalid credentials"))
		return
	} else if resp.StatusCode != http.StatusAccepted {
		_ = s.errorJSON(w, errors.New("error calling auth service"))
		return
	}

	var jsonFromService jsonResponse
	// decode the json from the auth service
	err = json.NewDecoder(resp.Body).Decode(&jsonFromService)
	if err != nil {
		_ = s.errorJSON(w, err)
		return
	}

	if jsonFromService.Error {
		_ = s.errorJSON(w, err, http.StatusUnauthorized)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "Authenticated",
		Data:    jsonFromService.Data,
	}

	_ = s.writeJSON(w, http.StatusAccepted, payload)
}
