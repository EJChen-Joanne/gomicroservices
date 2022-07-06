package main

import (
	"broker/events"
	"broker/logs"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/rpc"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayload  `json:"log,omitempty"`
	Mail   MailPayload `json:"mail,omitempty"`
}

type MailPayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (appli *Config) brokerHandler(w http.ResponseWriter, r *http.Request) {
	//receive JSON payload
	payload := new(JsonPayload)
	payload.Err = false
	payload.Msg = "Hit the broker"

	_ = appli.writeJSON(w, http.StatusOK, payload) //substitute the following codes
	/*
		//write these data out with json format
		out, _ := json.MarshalIndent(payload, "", "\t")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		w.Write(out)
	*/
}

func (appli *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	err := appli.readJSON(w, r, &requestPayload)
	if err != nil {
		appli.errorJSON(w, err)
		return
	}

	switch requestPayload.Action {
	case "auth":
		appli.authenticate(w, requestPayload.Auth)
	case "log":
		//appli.logItem(w, requestPayload.Log)
		//appli.logEventViaRabbitMQ(w, requestPayload.Log)
		appli.logEventViaRPC(w, requestPayload.Log)
	case "mail":
		appli.sendMail(w, requestPayload.Mail)
	default:
		appli.errorJSON(w, errors.New("unknown action"))

	}
}

func (appli *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	//create json files which sent to the auth-service
	jsonData, _ := json.MarshalIndent(a, "", "\t")

	//call the service
	request, err := http.NewRequest("POST", "http://auth-service/auth", bytes.NewBuffer(jsonData))
	if err != nil {
		appli.errorJSON(w, err)
		return
	}

	clients := &http.Client{}
	response, err := clients.Do(request)
	if err != nil {
		appli.errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	//make sure we receive the correct status code
	if response.StatusCode == http.StatusUnauthorized {
		appli.errorJSON(w, errors.New("invalid credentials"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		appli.errorJSON(w, errors.New("error: calling auth service"))
		return
	}

	//create a variable response.Body read into
	var jsonFromService JsonPayload

	//decode the json from the auth service
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		appli.errorJSON(w, err)
		return
	}

	if jsonFromService.Err {
		appli.errorJSON(w, err, http.StatusUnauthorized)
		return
	}

	var payload JsonPayload
	payload.Err = false
	payload.Msg = "Authentication checks!"
	payload.Data = jsonFromService.Data

	appli.writeJSON(w, http.StatusAccepted, payload)
}

func (appli *Config) logItem(w http.ResponseWriter, l LogPayload) {
	jsonData, _ := json.MarshalIndent(l, "", "\t")

	logServiceURL := "http://logger-service/log"

	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println(err) //debug
		appli.errorJSON(w, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		log.Println(err) //debug
		appli.errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		log.Println(err) //debug
		appli.errorJSON(w, err)
		return
	}

	var payload JsonPayload
	payload.Err = false
	payload.Msg = "Logged!"

	appli.writeJSON(w, http.StatusAccepted, payload)
}

func (appli *Config) sendMail(w http.ResponseWriter, mp MailPayload) {
	jsonData, _ := json.MarshalIndent(mp, "", "\t")

	//call the mail service
	mailServiceURL := "http://mail-service/send"

	//post to mail service
	request, err := http.NewRequest("POST", mailServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		appli.errorJSON(w, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		appli.errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		appli.errorJSON(w, errors.New("error calling mail service"))
		return
	}

	payload := JsonPayload{
		Err: false,
		Msg: fmt.Sprintf("Sent email to %s", mp.To),
	}

	appli.writeJSON(w, http.StatusAccepted, payload)
}

//handle logItem through emitting event into rabbitmq
func (appli *Config) logEventViaRabbitMQ(w http.ResponseWriter, l LogPayload) {
	err := appli.pushToQueue(l.Name, l.Data)
	if err != nil {
		appli.errorJSON(w, err)
		return
	}

	var payload JsonPayload
	payload.Err = false
	payload.Msg = "logged via RabbitMQ"

	appli.writeJSON(w, http.StatusAccepted, payload)
}

func (appli *Config) pushToQueue(name, msg string) error {
	emitter, err := events.NewEventEmitter(appli.Rabbit)
	if err != nil {
		return err
	}

	payload := LogPayload{
		Name: name,
		Data: msg,
	}

	j, _ := json.MarshalIndent(&payload, "", "\t")
	err = emitter.Push(string(j), "log.INFO")
	if err != nil {
		return err
	}

	return nil
}

type RPCPayload struct {
	Name string
	Data string
}

//handle log event via RPC server
func (appli *Config) logEventViaRPC(w http.ResponseWriter, l LogPayload) {
	client, err := rpc.Dial("tcp", "logger-service:5001")
	if err != nil {
		appli.errorJSON(w, err)
		return
	}

	//create a payload with type matching the remote RPC server
	rpcpayload := RPCPayload{
		Name: l.Name,
		Data: l.Data,
	}

	var res string
	err = client.Call("RPCServer.LogInfo", rpcpayload, &res)
	if err != nil {
		appli.errorJSON(w, err)
		return
	}

	payload := JsonPayload{
		Err: false,
		Msg: res,
	}

	appli.writeJSON(w, http.StatusAccepted, payload)

}

//grpc client
func (appli *Config) LogViaGRPC(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	err := appli.readJSON(w, r, &requestPayload)
	if err != nil {
		appli.errorJSON(w, err)
		return
	}

	conn, err := grpc.Dial("logger-service:50001", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		appli.errorJSON(w, err)
		return
	}
	defer conn.Close()

	c := logs.NewLogServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err = c.WriteLog(ctx, &logs.LogRequest{
		LogEntry: &logs.Log{
			Name: requestPayload.Log.Name,
			Data: requestPayload.Log.Data,
		},
	})
	if err != nil {
		appli.errorJSON(w, err)
		return
	}

	var payload JsonPayload
	payload.Err = false
	payload.Msg = "logged in!"

	appli.writeJSON(w, http.StatusAccepted, payload)

}
