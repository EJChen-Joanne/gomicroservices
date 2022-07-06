package main

import (
	"fmt"
	"log"
	"net/http"
)

func (appli *Config) SendMail(w http.ResponseWriter, r *http.Request) {
	type mailMsg struct {
		From    string `json:"from"`
		To      string `json:"to"`
		Subject string `json:"subject"`
		Message string `json:"message"`
	}

	var requestPayload mailMsg

	err := appli.readJSON(w, r, &requestPayload)
	if err != nil {
		log.Println(err)
		appli.errorJSON(w, err)
		return
	}

	msg := Message{
		From:    requestPayload.From,
		To:      requestPayload.To,
		Subject: requestPayload.Subject,
		Data:    requestPayload.Message,
	}

	err = appli.Mailer.SendSMTPMessage(msg)
	if err != nil {
		log.Println(err)
		appli.errorJSON(w, err)
		return
	}

	payload := JsonPayload{
		Err: false,
		Msg: fmt.Sprintf("sent to %s", requestPayload.To),
	}

	appli.writeJSON(w, http.StatusAccepted, payload)
}
