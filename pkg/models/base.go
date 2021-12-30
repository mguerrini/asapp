package models

type IIdentificable interface {
	GetId()  int
	SetId(int)
}

type identificable struct {
	Id int `json:"id"`
}

func  (i *identificable) GetId() int {
	return i.Id
}

func  (i *identificable) SetId(anId int) {
	i.Id = anId
}

