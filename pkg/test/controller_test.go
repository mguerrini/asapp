package test

import (
	"github.com/challenge/pkg/helpers"
	"github.com/challenge/pkg/modules/config"
	"github.com/challenge/pkg/test/core"
	"net/http"
	"testing"
)

func Test_CreateUserAndLogin(t *testing.T) {
	config.JoinSingleton("./pkg/test/configs", "dbtest.yml")

	username := "user1"
	pass := "111111"

	w, srv := core.StartTestServer()
	defer core.StopTestServer(w, srv)

	status, body := core.CreateUserRequest(username, pass)

	if status != 200 {
		t.Error(body)
		t.Fail()
		return
	}

	//the body has the id of de user created
	j := helpers.NewJsonHelper()
	createUserId, err := j.GetIntFromJson(body, "id")

	if err != nil || createUserId == nil || *createUserId <= 0 {
		t.Error("Invalid user created")
	}

	//LOGIN
	status, body = core.LoginRequest(username, pass)

	if status != 200 {
		t.Error(body)
		t.Fail()
		return
	}

	//Get the id and token
	loginUserId, _ := j.GetIntFromJson(body, "id")
	token, _ := j.GetStringFromJson(body, "token")

	if loginUserId == nil || *loginUserId <= 0 || *loginUserId != *createUserId {
		t.Error("Invalid login")
	}

	if len(*token) == 0 {
		t.Error("Invalid token")
	}
}

func Test_InvalidLogin(t *testing.T) {
	config.JoinSingleton("./pkg/test/configs", "dbtest.yml")

	username := "user2"
	pass := "222222"

	w, srv := core.StartTestServer()
	defer core.StopTestServer(w, srv)

	//LOGIN
	status, body := core.LoginRequest(username, pass)

	if status != http.StatusNotFound {
		t.Error("The user not exists")
		return
	}

	//create user
	_, body = core.CreateUserRequest(username, pass)

	//LOGIN With Invalid Pass
	status, body = core.LoginRequest(username, "989898989")

	if status == http.StatusOK {
		t.Error(body)
		t.Fail()
		return
	}

	//LOGIN With Invalid Username
	status, body = core.LoginRequest("rrrr", pass)

	if status == http.StatusOK {
		t.Error(body)
		t.Fail()
		return
	}

	//LOGIN OK
	status, body = core.LoginRequest(username, pass)

	if status != http.StatusOK {
		t.Error(body)
		t.Fail()
		return
	}
}

func Test_SendVideoMessage(t *testing.T) {
	config.JoinSingleton("./pkg/test/configs", "dbtest.yml")

	username1 := "user11"
	pass1 := "333333"

	username2 := "user12"
	pass2 := "444444"

	w, srv := core.StartTestServer()
	defer core.StopTestServer(w, srv)

	//create 2 users
	_, body1 := core.CreateUserRequest(username1, pass1)
	_, body2 := core.CreateUserRequest(username2, pass2)

	j := helpers.NewJsonHelper()
	userId1, _ := j.GetIntFromJson(body1, "id")
	userId2, _ := j.GetIntFromJson(body2, "id")

	//Login user1
	_, loginBody1 := core.LoginRequest(username1, pass1)
	token1, _ := j.GetStringFromJson(loginBody1, "token")

	//Login user2
	_, loginBody2 := core.LoginRequest(username2, pass2)
	token2, _ := j.GetStringFromJson(loginBody2, "token")

	//send message with valid token
	status, videoBody := core.SendVideoMessageRequest(*userId1, *userId2, "http://videos.com/video1", "vimeo", *token1)

	if status != http.StatusOK {
		t.Error("Error creating video")
		return
	}

	videoMessageId, _ := j.GetIntFromJson(videoBody, "id")

	//read message user 2
	count := 1
	status, bodyMessage := core.GetMessagesRequest(*userId2, *videoMessageId, &count, token2)

	if status != http.StatusOK {
		t.Error("Error reading video message")
		return
	}

	getMessageId, _ := j.GetIntFromJson(bodyMessage, "[0].id")
	getVideoType, _ := j.GetStringFromJson(bodyMessage, "[0].content.type")

	if *getMessageId != *videoMessageId {
		t.Error("Error getting video message")
		return
	}

	if *getVideoType != "video" {
		t.Error("Error getting video message")
		return
	}
}



