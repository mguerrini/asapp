package core


type VideoMessage struct {
	Id        int          `json:"id"`
	Sender    int          `json:"sender"`
	Recipient int          `json:"recipient"`
	Timestamp string       `json:"timestamp"`
	Video     VideoContent `json:"content"`
}

type TextMessage struct {
	Id        int         `json:"id"`
	Sender    int         `json:"sender"`
	Recipient int         `json:"recipient"`
	Timestamp string      `json:"timestamp"`
	Text      TextContent `json:"content"`
}

type ImageMessage struct {
	Id        int          `json:"id"`
	Sender    int          `json:"sender"`
	Recipient int          `json:"recipient"`
	Timestamp string       `json:"timestamp"`
	Image     ImageContent `json:"content"`
}

type VideoContent struct {
	Id     int    `json:"id"`
	Type   string `json:"type"`
	Url    string `json:"url"`
	Source string `json:"source"`
}


type ImageContent struct {
	Id     int    `json:"id"`
	Type   string `json:"type"`
	Url    string `json:"url"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
}

type TextContent struct {
	Id   int    `json:"id"`
	Type string `json:"type"`
	Text string `json:"text"`
}

