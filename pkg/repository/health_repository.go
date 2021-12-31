package repository

type IHealthRepository interface {
	Ping() error
}
