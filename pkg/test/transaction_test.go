package test

import (
	"github.com/challenge/pkg/modules/config"
	"github.com/challenge/pkg/test/core"
	"testing"
)

func Test_CreateUserThatFailsAfterSaveAndTryLogin(t *testing.T) {
	config.JoinSingleton("./pkg/test/configs", "dbtest.yml")

	username := "user100"
	pass := "100100"

	w, srv := core.StartTestServerThatFailsOnCreateUser()
	defer core.StopTestServer(w, srv)

	status, body := core.CreateUserRequest(username, pass)

	if status == 200 {
		t.Error(body)
		t.Fail()
		return
	}

	//LOGIN
	status, body = core.LoginRequest(username, pass)

	if status == 200 {
		t.Error(body)
		t.Fail()
		return
	}
}

