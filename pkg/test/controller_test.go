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
		t.Error("The user exists")
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

func Test_SendAndGetVideoMessage(t *testing.T) {
	config.JoinSingleton("./pkg/test/configs", "dbtest.yml")

	username1 := "user3"
	pass1 := "333333"

	username2 := "user4"
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

func Test_SendAndGetImageMessage(t *testing.T) {
	config.JoinSingleton("./pkg/test/configs", "dbtest.yml")

	username1 := "user5"
	pass1 := "555555"

	username2 := "user6"
	pass2 := "666666"

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
	status, imageBody := core.SendImageMessageRequest(*userId1, *userId2, 100, 200, "http://videos.com/video1", *token1)

	if status != http.StatusOK {
		t.Error("Error creating image")
		return
	}

	imageMessageId, _ := j.GetIntFromJson(imageBody, "id")

	//read message user 2
	count := 1
	status, bodyMessage := core.GetMessagesRequest(*userId2, *imageMessageId, &count, token2)

	if status != http.StatusOK {
		t.Error("Error reading image message")
		return
	}

	getMessageId, _ := j.GetIntFromJson(bodyMessage, "[0].id")
	getImageType, _ := j.GetStringFromJson(bodyMessage, "[0].content.type")

	if *getMessageId != *imageMessageId {
		t.Error("Error getting image message")
		return
	}

	if *getImageType != "image" {
		t.Error("Error getting image message")
		return
	}
}

func Test_SendAndGetTextMessage(t *testing.T) {
	config.JoinSingleton("./pkg/test/configs", "dbtest.yml")

	username1 := "user7"
	pass1 := "777777"

	username2 := "user8"
	pass2 := "888888"

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
	status, textBody := core.SendTextMessageRequest(*userId1, *userId2, "Este es un mensaje de texto.", *token1)

	if status != http.StatusOK {
		t.Error("Error creating text message")
		return
	}

	textMessageId, _ := j.GetIntFromJson(textBody, "id")

	//read message user 2
	count := 1
	status, bodyMessage := core.GetMessagesRequest(*userId2, *textMessageId, &count, token2)

	if status != http.StatusOK {
		t.Error("Error reading text message")
		return
	}

	getMessageId, _ := j.GetIntFromJson(bodyMessage, "[0].id")
	getTextType, _ := j.GetStringFromJson(bodyMessage, "[0].content.type")

	if *getMessageId != *textMessageId {
		t.Error("Error getting text message")
		return
	}

	if *getTextType != "text" {
		t.Error("Error getting text message")
		return
	}
}

func Test_InvalidToken(t *testing.T) {
	config.JoinSingleton("./pkg/test/configs", "dbtest.yml")

	username1 := "user9"
	pass1 := "999999"

	w, srv := core.StartTestServer()
	defer core.StopTestServer(w, srv)

	//create 2 users
	_, body1 := core.CreateUserRequest(username1, pass1)

	j := helpers.NewJsonHelper()
	userId1, _ := j.GetIntFromJson(body1, "id")

	//Login user1
	_, loginBody1 := core.LoginRequest(username1, pass1)
	token1, _ := j.GetStringFromJson(loginBody1, "token")
	token := *token1

	lastChar := token[len(token)-1]
	if lastChar != 'A' {
		lastChar = 'A'
	} else {
		lastChar = 'B'
	}

	token = token[:len(token)-1]
	token = token + string(lastChar)

	//send message with valid token
	status, _ := core.SendVideoMessageRequest(*userId1, 2, "http://videos.com/video1", "vimeo", token)

	if status == http.StatusOK {
		t.Error("Error validating token video")
		return
	}
}

