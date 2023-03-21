package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

type authRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s *Service) Authenticate(w http.ResponseWriter, r *http.Request) {
	reqPayload := authRequest{}

	err := s.readJSON(w, r, &reqPayload)
	if err != nil {
		log.Println(err)
		_ = s.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// validate the user exists in the database
	user, err := s.Models.User.GetByEmail(reqPayload.Email)
	if err != nil {
		log.Println(err)
		_ = s.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	// validate user's password
	passwordIsValid, err := user.PasswordMatches(reqPayload.Password)
	if err != nil || !passwordIsValid {
		log.Println(err)
		_ = s.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	// log authentication
	if err := s.logRequest(
		"authentication",
		fmt.Sprintf("%s logged in", user.Email),
	); err != nil {
		log.Printf("Error on log request: %v\n", err)
		_ = s.errorJSON(w, err)
		return
	}

	resPayload := jsonResponse{
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data:    user,
	}
	if err := s.writeJSON(w, http.StatusAccepted, resPayload); err != nil {
		log.Println(err)
		_ = s.errorJSON(w, err)
		return
	}
}

func (s *Service) logRequest(name, data string) error {
	var entry struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}

	entry.Name = name
	entry.Data = data

	jsonData, err := json.MarshalIndent(entry, "", "\t")
	if err != nil {
		return err
	}
	logServiceURL := "http://logger/log"

	request, err := http.NewRequest(
		http.MethodPost,
		logServiceURL,
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		log.Printf("Error on make logger request: %v\n", err)
		return err
	}

	client := &http.Client{}
	_, err = client.Do(request)
	if err != nil {
		log.Printf("Error on logger request: %v\n", err)
		return err
	}

	return nil
}
