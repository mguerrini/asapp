package models

type VideoSourceType string
type ContentType string

const (
	VideoSource_Vimeo VideoSourceType = "vimeo"
	VideoSource_Youtube VideoSourceType = "youtube"

	ContentType_Video ContentType = "video"
	ContentType_Image ContentType = "image"
	ContentType_Text  ContentType = "text"
)

type IMessageContent interface {
	Type() ContentType
}

type Message struct {
	identificable
	Sender    int             `json:"sender"`
	Recipient int             `json:"recipient"`
	Timestamp string          `json:"timestamp"`
	Content   IMessageContent `json:"content"`
}

//VIDEO

type VideoData struct {
	DataType ContentType     `json:"type"`
	Url      string          `json:"url"`
	Source   VideoSourceType `json:"source"`
}

func (v *VideoData) Type() ContentType {
	return ContentType_Video
}

// IMAGE

type ImageData struct {
	DataType ContentType `json:"type"`
	Url      string      `json:"url"`
	Height   int         `json:"height"`
	Width    int         `json:"width"`
}

func (i *ImageData) Type() ContentType {
	return ContentType_Image
}

// TEXT

type TextData struct {
	DataType ContentType `json:"type"`
	Text     string      `json:"text"`
}

func (i *TextData) Type() ContentType {
	return ContentType_Text
}
