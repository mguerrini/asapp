package services

import (
	"context"
	"github.com/challenge/pkg/models"
	"github.com/challenge/pkg/models/errors"
	"github.com/challenge/pkg/repository"
)

type MessageServices struct {
	messageRepository repository.IMessageRepository
	userServices      *UserServices
}

func NewMessageServices (sessionName string) *MessageServices {
	rep := repository.RepositoryFactory().CreateMessageRepository(sessionName)

	return &MessageServices{
		messageRepository: rep,
		userServices: NewUserServices(sessionName),
	}
}

func (s *MessageServices) SendMessage(ctx context.Context, msg *models.Message) error {
	err := validateMessage(msg)
	if err != nil {
		return err
	}

	//validate recipient and sender exists
	senderProfile, errSender := s.userServices.GetUserProfileById(ctx, msg.Sender)

	if errSender != nil {
		return errSender
	} else if senderProfile == nil {
		return errors.NewNotFoundMsg("The sender not exist.")
	}

	recipientProfile, errRec := s.userServices.GetUserProfileById(ctx, msg.Recipient)

	if errRec != nil {
		return errRec
	} else if recipientProfile == nil {
		return errors.NewNotFoundMsg("The recipient not exist.")
	}

	err = s.messageRepository.CreateMessage(ctx, msg)

	if err != nil {
		return err
	}

	return nil
}

// The searches could be implemented in another component in charge of carrying out different types of searches, returning DTOs and evolving separately
func (s *MessageServices) SearchMessages(ctx context.Context, recipientId int, startId int, count int) ([]models.Message, error) {
	if count <= 0 {
		count = 100
	}

	recipientProfile, errRec := s.userServices.GetUserProfileById(ctx, recipientId)

	if errRec != nil {
		return nil, errRec
	} else if recipientProfile == nil {
		return nil, errors.NewNotFoundMsg("The recipient not exist.")
	}

	return s.messageRepository.SearchMessages(ctx, recipientId, startId, count)
}


func validateMessage(msg *models.Message) error {
	if msg == nil {
		return errors.NewBadRequestMsg("Message argument nil")
	}

	if msg.Recipient <= 0 {
		return errors.NewBadRequestMsg("Invalid recipient argument")
	}

	if msg.Sender <= 0 {
		return errors.NewBadRequestMsg("Invalid sender argument")
	}

	switch msg.Content.Type() {
	case models.ContentType_Image:
		return validateImage(msg.Content)
	case models.ContentType_Video:
		return validateVideo(msg.Content)
	}

	return nil
}

func validateVideo(content models.IMessageContent) error {
	video := content.(*models.VideoData)

	if video == nil {
		return errors.NewBadRequestMsg("Video argument nil")
	}

	if video.Url == "" {
		return errors.NewBadRequestMsg("Video url empty")
	}

	if video.Source != models.VideoSource_Vimeo && video.Source != models.VideoSource_Youtube {
		return errors.NewBadRequestMsg("Invalid video source")
	}

	return nil
}

func validateImage(content models.IMessageContent) error {
	img := content.(*models.ImageData)

	if img == nil {
		return errors.NewBadRequestMsg("Image argument nil")
	}

	if img.Url == "" {
		return errors.NewBadRequestMsg("Image url empty")
	}

	if img.Width <= 0 || img.Height <= 0 {
		return errors.NewBadRequestMsg("Invalid image size")
	}

	return nil
}
