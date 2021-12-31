package controller

import (
	"context"
	"net/http"

	"github.com/challenge/pkg/helpers"
	"github.com/challenge/pkg/models"
)

// Check returns the health of the service and DB
func (h *Handler) Check(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	err := h.healthRepository.Ping()

	if err == nil {
		helpers.RespondJSON(w, models.Health{Status: "ok"})
	} else {
		helpers.RespondJSON(w, models.Health{Status: "fail"})
	}
}
