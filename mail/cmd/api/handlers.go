package main

import (
	"log"
	"net/http"
	"os"
)

type mailMessage struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

func (s *Service) SendMail(w http.ResponseWriter, r *http.Request) {
	var reqPayload mailMessage
	if err := s.readJSON(w, r, &reqPayload); err != nil {
		_ = s.errorJSON(w, err)
		return
	}

	var fromEmail string
	if reqPayload.From != "" {
		fromEmail = reqPayload.From
	} else {
		fromEmail = os.Getenv("MAIL_FROM_ADDRESS")
	}
	msg := Message{
		From:    fromEmail,
		To:      reqPayload.To,
		Subject: reqPayload.Subject,
		Data:    reqPayload.Message,
	}

	if err := s.Mailer.SendSMTPMessage(msg); err != nil {
		log.Printf("Failed sending mail: %v\n", err)
		_ = s.errorJSON(w, err)
		return
	}

	log.Println("Mail successfully sent!")

	_ = s.writeJSON(w, http.StatusAccepted, jsonResponse{
		Message: "mail sent to " + reqPayload.To,
	})
}
