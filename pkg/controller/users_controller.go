package controller

import (
	"context"
	"net/http"

	"github.com/challenge/pkg/helpers"
	"github.com/challenge/pkg/models"
)

// CreateUser creates a new user
func (h Handler) CreateUser(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	// TODO: Create a New User
	helpers.RespondJSON(w, models.User{})
}
