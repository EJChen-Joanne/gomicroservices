package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type JsonPayload struct {
	Err  bool   `json:"error"`
	Msg  string `json:"message"`
	Data any    `json:"data,omitempty"` //golang could use type any over v1.18
}

func (appli *Config) readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1048576 //one megabytes

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	if err != nil {
		return nil
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body has only a single JSON value")
	}

	return nil
}

func (appli *Config) writeJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	//if any headers were included as the final parameters to the function
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

func (appli *Config) errorJSON(w http.ResponseWriter, err error, status ...int) error {
	//write error messages as JSON file
	statusCode := http.StatusBadRequest //default status if it's not specified in the call to the func

	if len(status) > 0 {
		//something is specified
		statusCode = status[0]
	}

	var payload JsonPayload
	payload.Err = true
	payload.Msg = err.Error()

	return appli.writeJSON(w, statusCode, payload)
}
