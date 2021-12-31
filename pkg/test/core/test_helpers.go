package core

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/challenge/pkg/models"
	"github.com/challenge/pkg/modules/logger"
	"github.com/challenge/pkg/server"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

func StopTestServer(w *sync.WaitGroup, srv *http.Server) {
	if r := recover(); r != nil {
		logger.Error("Error on finished test" ,r.(error))
	}

	srv.Shutdown(context.TODO())
	w.Wait()
}

func StartTestServer () (*sync.WaitGroup, *http.Server) {
	httpServerExitDone := &sync.WaitGroup{}
	httpServerExitDone.Add(1)

	srv := server.StartHttpServer(httpServerExitDone)

	//wait 1 second
	time.Sleep(time.Duration(1) * time.Second)

	return httpServerExitDone, srv
}


func CreateUserRequest(username string, pass string) (int, string) {
	user := models.User{
		Id:       0,
		Username: username,
		Password: pass,
	}

	return doPost("http://localhost:8080/users", user, nil)
}

func LoginRequest(username string, pass string) (int, string) {
	user := models.Login{
		Username: username,
		Password: pass,
	}

	return doPost("http://localhost:8080/login", user, nil)
}

func SendVideoMessageRequest(sender, recipient int, url, source string, token string) (int, string) {
	msg := VideoMessage{
		Sender:    sender,
		Recipient: recipient,
		Video:     VideoContent{
			Type:   "video",
			Url:    url,
			Source: source,
		},
	}

	return doPost("http://localhost:8080/messages", msg, &token)
}

func SendTextMessageRequest(sender, recipient int, text string, token string) (int, string) {
	msg := TextMessage{
		Sender:    sender,
		Recipient: recipient,
		Text:      TextContent{
			Type: "text",
			Text: text,
		},
	}

	return doPost("http://localhost:8080/messages", msg, &token)
}

func SendImageMessageRequest(sender, recipient int, w, h int, url string, token string) (int, string) {
	msg := ImageMessage{
		Sender:    sender,
		Recipient: recipient,
		Image:     ImageContent{
			Type:   "image",
			Url:    url,
			Height: h,
			Width:  w,
		},
	}

	return doPost("http://localhost:8080/messages", msg, &token)
}

func doPost(url string, bodyObj interface{}, token *string) (int, string) {
	body, _ := json.Marshal(bodyObj)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")
	if token != nil {
		req.Header.Set("Authorization", "Bearer " + *token)
	}

	client := &http.Client{}
	resp, err := client.Do(req)

	//resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	status := resp.StatusCode
	respBody, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		//Failed to read response.
		panic(err)
	}

	jsonStr := string(respBody)

	return status, jsonStr
}

func GetMessagesRequest(recipient int, start int, count *int, token *string) (int, string) {
	query := ""
	if count == nil {
		query = fmt.Sprintf("http://localhost:8080/messages?recipient=%d&start=%d&limit=%d", recipient, start, count)
	} else {
		query = fmt.Sprintf("http://localhost:8080/messages?recipient=%d&start=%d&limit=", recipient, start)
	}

	req, err := http.NewRequest("GET", query, bytes.NewBuffer(nil))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")
	if token != nil {
		req.Header.Set("Authorization", "Bearer " + *token)
	}

	client := &http.Client{}
	resp, err := client.Do(req)

	//resp, err := http.Get(query)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	status := resp.StatusCode
	respBody, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		//Failed to read response.
		panic(err)
	}

	jsonStr := string(respBody)

	return status, jsonStr
}

