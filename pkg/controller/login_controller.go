package controller

import (
	"context"
	"github.com/challenge/pkg/modules/logger"
	"net/http"

	"github.com/challenge/pkg/helpers"
	"github.com/challenge/pkg/models"
)



// Login authenticates a user and returns a token
func (h Handler) Login(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	// TODO: User must login and a token must be generated
	cred := models.Login{}
	err := helpers.BindJSON(r,&cred)

	if err != nil {
		http.Error(w, "Invalid login data", http.StatusBadRequest)
		logger.Error("Invalid login data", err)
		return
	}

	err = h.authServices.ValidateUser(ctx, cred)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
	}

	profile, err := h.userServices.GetUserProfile(cred.Username)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Error("Error getting user profile for " + cred.Username, err)
		return
	}

	token := h.authServices.GenerateToken(ctx, cred.Username)

	//Id = 0. I dont know if I must response with the user id or just 0
	resp := LoginResponse{
		Id:    profile.Id,
		Token: token,
	}

	helpers.RespondJSON(w, resp)
}
