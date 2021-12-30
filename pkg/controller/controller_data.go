package controller

type LoginResponse struct {
	Id int `json:"id"`
	Token string `json:"token"`
}


type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CreateUserResponse struct {
	Id int `json:"id"`
}

type SendMessageRequest struct {
	Sender    int                       `json:"sender"`
	Recipient int                       `json:"recipient"`
	Content   SendMessageContentRequest `json:"content"`
}

//Another options is to implement UnmarshalJSON and deserialize the different contents
//For this case is not necessary
type SendMessageContentRequest struct {
	Type string `json:"type"`

	//Image and Video
	Url string `json:"url"`
	
	//Image
	Height int `json:"height"`
	Width int `json:"width"`

	//Text
	Text string `json:"text"`

	//Video
	Source string `json:"source"`
}

type SendMessageResponse struct {
	Id        int    `json:"id"`
	Timestamp string `json:"timestamp"`
}