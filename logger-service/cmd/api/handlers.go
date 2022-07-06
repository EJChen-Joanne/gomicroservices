package main

import (
	"logger-service/data"
	"net/http"
)

type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (appli *Config) WriteLog(w http.ResponseWriter, r *http.Request) {
	//read json into varialble
	var requestPayload Payload

	_ = appli.readJSON(w, r, &requestPayload)

	//insert data
	event := data.LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}

	err := appli.Models.LogEntry.Insert(event)
	if err != nil {
		appli.errorJSON(w, err)
		return
	}

	response := JsonPayload{
		Err: false,
		Msg: "Logged",
	}

	appli.writeJSON(w, http.StatusAccepted, response)
}
