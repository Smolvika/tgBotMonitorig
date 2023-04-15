package repository

import "github.com/jmoiron/sqlx"

type Cost interface {
	AddNewUser(chatId int, currency string) error
	AllChatIdCostDB(currency string) ([]userCost, error)
	DeleteUserChatIdCostDB(chatId int) error
	AllUserChatIdCostDB() ([]int64, error)
}

type ChangeCost interface {
	AddUserChangeCostDB(chatId int, Currency string, cost float64) error
	AllChatIdChangeCostDB(currency string, cost float64) ([]userChangeCost, error)
	DeleteUserChatIdChangeCostDB(chatId int) error
	AllUserChatIdChangeCostDB() ([]int64, error)
}

type Repository struct {
	Cost
	ChangeCost
}

func New(db *sqlx.DB) *Repository {
	return &Repository{
		Cost:       NewCostPostgres(db),
		ChangeCost: NewChangeCostPostgres(db),
	}
}
