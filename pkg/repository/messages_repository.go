package repository

import (
	"context"
	"github.com/challenge/pkg/models"
)

type IMessageRepository interface {
	CreateMessage(ctx context.Context, msg *models.Message) error

	// The searches could be implemented in another component in charge of carrying out different types of searches, returning DTOs and evolving separately
	SearchMessages(ctx context.Context, recipientId int, startId int, count int) ([]models.Message, error)
}




