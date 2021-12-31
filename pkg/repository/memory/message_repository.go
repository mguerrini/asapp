package memory

import (
	"context"
	"github.com/challenge/pkg/helpers"
	"github.com/challenge/pkg/models"
	"sync"
	"sync/atomic"
	"time"
)

var 	messages []models.Message

func init() {
	messages = make([]models.Message, 0)
}

type memoryMessageRepository struct {
	sync     sync.Mutex
	idSeq    int32
}

func NewMemoryMessageRepository () *memoryMessageRepository {
	return &memoryMessageRepository{
	}
}


func (m *memoryMessageRepository) CreateMessage(ctx context.Context, msg *models.Message) error {
	msgId := atomic.AddInt32(&m.idSeq, 1)

	msg.SetId(int(msgId))
	msg.Timestamp = time.Now().UTC().Format("2006-01-02T15:04:05Z")

	copy := models.Message{}
	switch msg.Content.Type() {
	case models.ContentType_Image:
		copy.Content = &models.ImageData{}

	case models.ContentType_Video:
		copy.Content = &models.VideoData{}

	case models.ContentType_Text:
		copy.Content = &models.TextData{}
	}

	helpers.Copy(msg, &copy)

	m.sync.Lock()
	defer m.sync.Unlock()

	messages = append(messages, copy)
	return nil
}

func (m *memoryMessageRepository) SearchMessages(ctx context.Context, recipientId int, startId int, count int) ([]models.Message, error) {
	m.sync.Lock()
	defer m.sync.Unlock()

	output := make([]models.Message, 0)

	for _, m := range messages {
		if m.Recipient == recipientId &&  m.GetId() >= startId {
			output = append(output, m)

			if len (output) >= count {
				break
			}
		}
	}

	return output, nil
}