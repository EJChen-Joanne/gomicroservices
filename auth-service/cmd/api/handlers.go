package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func (appli *Config) Auth(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := appli.readJSON(w, r, &requestPayload)
	if err != nil {
		appli.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	//validate the user against the DB
	user, err := appli.Models.User.GetByEmail(requestPayload.Email)
	if err != nil {
		appli.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		appli.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	//log authentication
	err = appli.logRequest("authentication", fmt.Sprintf("%s logged in", user.Email))
	if err != nil {
		appli.errorJSON(w, err)
		return
	}

	payload := JsonPayload{
		Err:  false,
		Msg:  fmt.Sprintf("Logged in user %s", user.Email),
		Data: user,
	}

	appli.writeJSON(w, http.StatusAccepted, payload)
}

func (appli *Config) logRequest(name, data string) error {
	var entry struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}

	entry.Name = name
	entry.Data = data

	jsonData, _ := json.MarshalIndent(entry, "", "\t")
	logServiceURL := "http://logger-service/log"

	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	client := &http.Client{}

	_, err = client.Do(request)
	if err != nil {
		return err
	}

	return nil
}
