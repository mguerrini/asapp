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