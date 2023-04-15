package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type CostPostgres struct {
	db *sqlx.DB
}

func NewCostPostgres(db *sqlx.DB) *CostPostgres {
	return &CostPostgres{db: db}
}

type userCost struct {
	UsersId  int    `db:"user_id"`
	ChatId   int64  `db:"chat_id"`
	Сurrency string `db:"currency"`
}

func (p *CostPostgres) AddNewUser(chatId int, currency string) error {
	_, err := p.db.Exec("INSERT INTO currency (chat_id,currency) VALUES ($1,$2)", chatId, currency)
	if err != nil {
		return err
	}
	return nil
}

func (p *CostPostgres) AllChatIdCostDB(currency string) ([]userCost, error) {
	var user []userCost
	query := fmt.Sprint("SELECT chat_id FROM currency WHERE currency = $1")
	err := p.db.Select(&user, query, currency)
	return user, err
}

func (p *CostPostgres) DeleteUserChatIdCostDB(chatId int) error {
	_, err := p.db.Exec("DELETE FROM currency WHERE chat_id = $1", chatId)
	return err
}

func (p *CostPostgres) AllUserChatIdCostDB() ([]int64, error) {
	var user []int64
	query := fmt.Sprint("SELECT DISTINCT chat_id FROM currency")
	err := p.db.Select(&user, query)
	return user, err
}
