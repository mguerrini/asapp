package sql

import (
	"context"
	"github.com/challenge/pkg/helpers"
	"github.com/challenge/pkg/models"
	"github.com/challenge/pkg/models/errors"
	"github.com/challenge/pkg/modules/storage"
)

type sqlMessageRepository struct {
	dbCnn storage.IDBConnection
}

func NewSqlMessageRepository (sessionName string) (*sqlMessageRepository, error) {
	//get connection string
	dbConn, err := storage.DBManagerSingleton().CreateConnection(sessionName)

	if err != nil {
		return nil, err
	}
	return &sqlMessageRepository{
		dbCnn: dbConn,
	}, nil
}


func (s *sqlMessageRepository) CreateMessage(ctx context.Context, msg *models.Message) error {
	query :=
	`INSERT INTO Messages (
		SenderId,
		RecipientId,
		Timestamp
	)
	VALUES (
		@senderId,
		@recipientId,
		datetime('now')
	);`

	res, err := s.dbCnn.Exec(ctx, query, msg.Sender, msg.Recipient)

	if err != nil {
		return err
	}

	id, err := res.LastInsertId()

	if err != nil {
		return err
	}

	msg.Id = int (id)
	err = s.createImageContent(ctx, msg)
	if err != nil {
		return err
	}
	err = s.createTextContent(ctx, msg)
	if err != nil {
		return err
	}
	err = s.createVideoContent(ctx, msg)
	if err != nil {
		return err
	}

	//set timestamp
	newMsg, _ := s.SearchMessages(ctx, msg.Recipient, msg.Id, 1)

	if newMsg != nil && len(newMsg) == 1{
		msg.Timestamp = newMsg[0].Timestamp
	}

	return nil
}

func (s *sqlMessageRepository) createVideoContent(ctx context.Context, msg *models.Message) error {
	vid, ok := msg.Content.(*models.VideoData)

	if !ok {
		return nil
	}

	query :=
	`INSERT INTO VideoContent (
		MessageId,
		Url,
		Source
	)
		VALUES (
		@MessageId,
		@Url,
		@Source
	);`

	res, err := s.dbCnn.Exec(ctx, query, msg.Id, vid.Url, vid.Source)

	if err != nil {
		return err
	}

	id, err := res.LastInsertId()

	if err != nil {
		return err
	}

	vid.Id = int (id)
	return nil
}

func (s *sqlMessageRepository) createTextContent(ctx context.Context, msg *models.Message) error {
	txt, ok := msg.Content.(*models.Text)

	if !ok {
		return nil
	}

	query :=
	`INSERT INTO TextContent (
		MessageId,
		Text
	)
	VALUES (
		@MessageId,
		@Text
	);`

	res, err := s.dbCnn.Exec(ctx, query, msg.Id, txt.Text)

	if err != nil {
		return err
	}

	id, err := res.LastInsertId()

	if err != nil {
		return err
	}

	txt.Id = int (id)
	return nil
}

func (s *sqlMessageRepository) createImageContent(ctx context.Context, msg *models.Message) error {
	img, ok := msg.Content.(*models.ImageData)

	if !ok {
		return nil
	}

	query :=
	`INSERT INTO ImageContent (
		MessageId,
		Url,
		Height,
		Width
	)
	VALUES (
		@MessageId,
		@Url,
		@Height,
		@Width
	);`

	res, err := s.dbCnn.Exec(ctx, query, msg.Id, img.Url, img.Height, img.Width)

	if err != nil {
		return err
	}

	id, err := res.LastInsertId()

	if err != nil {
		return err
	}

	img.Id = int (id)
	return nil
}



func (s *sqlMessageRepository) SearchMessages(ctx context.Context, recipientId int, startId int, count int) ([]models.Message, error) {

	query :=
		`SELECT m.*,
		CASE
			WHEN txt.Id IS NOT NULL THEN 'text'
			WHEN vd.Id IS NOT NULL THEN 'video'
			WHEN img.Id IS NOT NULL THEN 'image'
		END as ContentType,
		txt.Id AS text_Id,
		txt.Text AS text_Text,
		vd.Id AS video_Id,
		vd.Url AS video_Url,
		vd.Source AS video_Source,
		img.Id AS image_Id,
		img.Url AS image_Url,
		img.Height AS image_Height,
		img.Width AS image_Width

	FROM Messages m
	LEFT JOIN TextContent txt ON m.Id = txt.MessageId
	LEFT JOIN VideoContent vd ON m.Id = vd.MessageId
	LEFT JOIN ImageContent img ON m.Id = img.MessageId

	WHERE @startId <= m.Id AND m.Id < @finishId AND @recipient = m.RecipientId`

	finishId := startId + count
	rows, err := s.dbCnn.Query(ctx, query, startId, finishId, recipientId)

	defer func() {
		if rows != nil {
			rows.Close();
		}
	}()

	if err != nil {
		return nil, err
	}

	output := make([]models.Message, 0)

	for rows.Next() {
		row, err := helpers.GetMapFromReader(rows)
		if err != nil {
			return nil, err
		}

		msg, err := s.buildMessage(row)
		output = append(output, *msg)
	}

	return output, nil
}

func (s *sqlMessageRepository) buildMessage(row map[string]interface{}) (*models.Message, error) {
	msg := &models.Message{
		Sender:    helpers.GetInt("SenderId", row),
		Recipient: helpers.GetInt("RecipientId", row),
		Timestamp: helpers.GetString("Timestamp", row),
		Content:   nil,
	}

	msg.Id = helpers.GetInt("Id", row)
	contentType := helpers.GetString("ContentType", row)

	if contentType == "text" {
		txt := &models.Text{Text: helpers.GetString("text_Text", row)}
		txt.Id = helpers.GetInt("text_Id", row)
		msg.Content = txt
	} else if contentType == "video" {
		vid := &models.VideoData{
			Url:    helpers.GetString("video_Url", row),
			Source: models.VideoSourceType(helpers.GetString("video_Source", row)),
		}
		vid.Id = helpers.GetInt("video_Id", row)
		msg.Content = vid
	} else if contentType == "image" {
		img := &models.ImageData{
			Url:    helpers.GetString("image_Url", row),
			Height: helpers.GetInt("image_Height", row),
			Width:  helpers.GetInt("image_Width", row),
		}
		img.Id = helpers.GetInt("image_Id", row)
		msg.Content = img
	} else {
		return nil, errors.NewInternalServerErrorMsg("Invalid content type")
	}

	return msg, nil
}



