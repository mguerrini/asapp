package controller

import (
	"context"
	"github.com/challenge/pkg/modules/logger"
	"net/http"

	"github.com/challenge/pkg/helpers"
	"github.com/challenge/pkg/models"
)



// Login authenticates a user and returns a token
func (h *Handler) Login(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	// TODO: User must login and a token must be generated
	cred := models.Login{}
	err := helpers.BindJSON(r,&cred)

	if err != nil {
		http.Error(w, "Login Error - Invalid login data", http.StatusBadRequest)
		logger.Error("Login Error - Invalid login data", err)
		return
	}

	err = h.authServices.ValidateUser(ctx, cred)
	if err != nil {
		http.Error(w, err.Error(), helpers.GetStatusCodeOr(err, http.StatusUnauthorized))
		return
	}

	profile, err := h.userServices.GetUserProfileByUsername(ctx, cred.Username)

	if err != nil {
		http.Error(w, "Login error - " + err.Error(), helpers.GetStatusCodeOr(err, http.StatusInternalServerError))
		logger.Error("Error getting user profile for " + cred.Username, err)
		return
	}

	if profile == nil {
		http.Error(w, "Login error - User not exist", http.StatusNotFound)
		logger.Warn("User not exists: " + cred.Username)
		return
	}

	token, err := h.authServices.GenerateToken(ctx, cred.Username)

	if err != nil {
		http.Error(w, "Error generating Token - " + err.Error(), http.StatusInternalServerError)
		logger.Error("Error generating Token", err)
		return
	}

	resp := LoginResponse{
		Id:    profile.Id,
		Token: token,
	}

	helpers.RespondJSON(w, resp)
}
