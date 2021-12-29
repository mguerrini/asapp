package controller

import (
	"context"
	"github.com/challenge/pkg/modules/logger"
	"net/http"

	"github.com/challenge/pkg/helpers"
	"github.com/challenge/pkg/models"
)

// CreateUser creates a new user
func (h Handler) CreateUser(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	reqData := CreateUserRequest{}
	err :=  helpers.BindJSON(r, &reqData)

	if err != nil {
		http.Error(w, "Invalid data", http.StatusBadRequest)
		logger.Error("Invalid data for 'CreateUser'", err)
		return
	}

	user := models.User{
		Id:       0,
		Username: reqData.Username,
		Password: reqData.Password,
	}

	newUser, err := h.userServices.CreateUser(ctx, user)

	if err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		logger.Error("Error on 'CreateUser'", err)
		return
	}

	resData := CreateUserResponse{Id: newUser.Id}

	helpers.RespondJSON(w, resData)
}
