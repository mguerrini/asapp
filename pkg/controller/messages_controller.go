package controller

import (
	"context"
	"github.com/challenge/pkg/models/errors"
	"github.com/challenge/pkg/modules/logger"
	"net/http"
	"strconv"

	"github.com/challenge/pkg/helpers"
	"github.com/challenge/pkg/models"
)

// SendMessage send a message from one user to another
func (h *Handler) SendMessage(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	reqData := SendMessageRequest{}
	err := helpers.BindJSON(r, &reqData)

	if err != nil {
		http.Error(w, "Invalid data - "+err.Error(), helpers.GetStatusCodeOr(err, http.StatusBadRequest))
		logger.Error("Invalid data for 'SendMessage'", err)
		return
	}

	msg, err := h.createMessage(reqData)
	if err != nil {
		http.Error(w, err.Error(), helpers.GetStatusCodeOr(err, http.StatusBadRequest))
		return
	}

	err = h.msgService.SendMessage(ctx, msg)

	if err != nil {
		http.Error(w, "Error sending message - " + err.Error(), helpers.GetStatusCodeOr(err, http.StatusInternalServerError))
		logger.Error("Error Sending Message.", err)
		return
	}

	res := SendMessageResponse{
		Id:        msg.GetId(),
		Timestamp: msg.Timestamp,
	}

	helpers.RespondJSON(w, res)
}

func (h *Handler) createMessage(reqData SendMessageRequest) (*models.Message, error) {
	msg := models.Message{
		Sender:    reqData.Sender,
		Recipient: reqData.Recipient,
	}

	switch reqData.Content.Type {
	case string(models.ContentType_Video):
		video := &models.VideoData{}
		helpers.Copy(reqData.Content, video)
		msg.Content = video
	case string(models.ContentType_Image):
		img := &models.ImageData{}
		helpers.Copy(reqData.Content, img)
		msg.Content = img
	case string(models.ContentType_Text):
		txt := &models.TextData{}
		helpers.Copy(reqData.Content, txt)
		msg.Content = txt
	default:
		logger.Warn("Invalid data for 'SendMessage'")
		return nil, errors.NewBadRequestMsg("Invalid message content type")
	}

	return &msg, nil
}


// GetMessages get the messages from the logged user to a recipient
func (h *Handler) GetMessages(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	recipientId, err := h.getIntQueryValue(r, "recipient", false)
	if err != nil {
		http.Error(w, err.Error(), helpers.GetStatusCodeOr(err, http.StatusBadRequest))
		return
	}
	startId, err := h.getIntQueryValue(r, "start", false)
	if err != nil {
		http.Error(w, err.Error(), helpers.GetStatusCodeOr(err, http.StatusBadRequest))
		return
	}

	count, err := h.getIntQueryValue(r, "limit", true)
	if err != nil {
		http.Error(w, err.Error(), helpers.GetStatusCodeOr(err, http.StatusBadRequest))
		return
	}

	msgs, err := h.msgService.SearchMessages(ctx, recipientId, startId, count)

	if err != nil {
		http.Error(w, "Error searching messages - " + err.Error(), helpers.GetStatusCodeOr(err, http.StatusInternalServerError))
		return
	}

	helpers.RespondJSON(w, msgs)
}

func (h *Handler) getIntQueryValue(r *http.Request, key string, optional bool) (int, error) {
	query := r.URL.Query()
	intStr := query.Get(key)

	if intStr == "" {
		if !optional {
			logger.Warn("Null " + key + " paramenter.")
			return 0, errors.NewBadRequestMsg("Parameter " + key + " is null")
		} else {
			return 0, nil
		}
	}

	intValue, err := strconv.Atoi(intStr)
	if err != nil {
		logger.Warn("Invalid value for parameter " + key )
		return 0, errors.NewBadRequestMsg("Invalid value for " + key + " parameter")
	}

	return intValue, nil
}
