package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type userCost struct {
	UsersId  int    `db:"user_id"`
	ChatId   int64  `db:"chat_id"`
	Сurrency string `db:"currency"`
}
type userChangeCost struct {
	UsersId  int     `db:"user_id"`
	ChatId   int64   `db:"chat_id"`
	Сurrency string  `db:"currency"`
	Cost     float64 `db:"cost_currency"`
}

func AddNewUser(chatId int, currency string, db *sqlx.DB) error {
	_, err := db.Exec("INSERT INTO currency (chat_id,currency) VALUES ($1,$2)", chatId, currency)
	if err != nil {
		return err
	}
	return nil
}

func AllChatIdCostDB(currency string, db *sqlx.DB) ([]userCost, error) {
	var user []userCost
	query := fmt.Sprint("SELECT chat_id FROM currency WHERE currency = $1")
	err := db.Select(&user, query, currency)
	return user, err
}
func DeleteUserChatIdCostDB(chatId int, db *sqlx.DB) error {
	_, err := db.Exec("DELETE FROM currency WHERE chat_id = $1", chatId)
	return err
}
func AllUserChatIdCostDB(db *sqlx.DB) ([]int64, error) {
	var user []int64
	query := fmt.Sprint("SELECT DISTINCT chat_id FROM currency")
	err := db.Select(&user, query)
	return user, err

}

func AddUserChangeCostDB(chatId int, Currency string, cost float64, db *sqlx.DB) error {
	_, err := db.Exec("INSERT INTO change_currency (chat_id, currency, cost_currency) VALUES ($1,$2,$3)", chatId, Currency, cost)
	return err
}

func AllChatIdChangeCostDB(currency string, cost float64, db *sqlx.DB) ([]userChangeCost, error) {
	var user []userChangeCost
	query := fmt.Sprint("SELECT chat_id FROM change_currency WHERE currency = $1 AND cost_currency > $2")
	err := db.Get(&user, query, currency, cost)
	if err != nil {
		return user, err
	}
	_, err = db.Exec("DELETE FROM change_currency WHERE currency = $1 AND cost_currency > $2 ", currency, cost)
	return user, err

}
func DeleteUserChatIdChangeCostDB(chatId int, db *sqlx.DB) error {
	_, err := db.Exec("DELETE FROM change_currency WHERE chat_id = $1", chatId)
	return err
}
func AllUserChatIdChangeCostDB(db *sqlx.DB) ([]int64, error) {
	var user []int64
	query := fmt.Sprint("SELECT DISTINCT chat_id FROM change_currency")
	err := db.Select(&user, query)
	return user, err

}
