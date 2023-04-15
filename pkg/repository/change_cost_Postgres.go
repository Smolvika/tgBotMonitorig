package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type userChangeCost struct {
	UsersId  int     `db:"user_id"`
	ChatId   int64   `db:"chat_id"`
	Сurrency string  `db:"currency"`
	Cost     float64 `db:"cost_currency"`
}

type ChangeCostPostgres struct {
	db *sqlx.DB
}

func NewChangeCostPostgres(db *sqlx.DB) *ChangeCostPostgres {
	return &ChangeCostPostgres{db: db}
}

func (p *ChangeCostPostgres) AddUserChangeCostDB(chatId int, Currency string, cost float64) error {
	_, err := p.db.Exec("INSERT INTO change_currency (chat_id, currency, cost_currency) VALUES ($1,$2,$3)", chatId, Currency, cost)
	return err
}

func (p *ChangeCostPostgres) AllChatIdChangeCostDB(currency string, cost float64) ([]userChangeCost, error) {
	var user []userChangeCost
	query := fmt.Sprint("SELECT chat_id FROM change_currency WHERE currency = $1 AND cost_currency > $2")
	err := p.db.Get(&user, query, currency, cost)
	if err != nil {
		return user, err
	}
	_, err = p.db.Exec("DELETE FROM change_currency WHERE currency = $1 AND cost_currency > $2 ", currency, cost)
	return user, err

}

func (p *ChangeCostPostgres) DeleteUserChatIdChangeCostDB(chatId int) error {
	_, err := p.db.Exec("DELETE FROM change_currency WHERE chat_id = $1", chatId)
	return err
}

func (p *ChangeCostPostgres) AllUserChatIdChangeCostDB() ([]int64, error) {
	var user []int64
	query := fmt.Sprint("SELECT DISTINCT chat_id FROM change_currency")
	err := p.db.Select(&user, query)
	return user, err
}
