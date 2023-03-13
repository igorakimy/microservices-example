package main

import (
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

	resPayload := jsonResponse{
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data:    user,
	}
	if err := s.writeJSON(w, http.StatusAccepted, &resPayload); err != nil {
		log.Fatal(err)
		return
	}
}
