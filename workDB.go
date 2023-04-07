package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

// база данных для уведомлений ежечасно
func addNewUser(chatId int, db *sql.DB) error {
	rows, err := db.Query("select * from notif_currency where chatId = ?", chatId)
	if err != nil {
		return err
	}
	var us userСurrency
	for rows.Next() {
		err = rows.Scan(&us.UsersId, &us.ChatId, &us.Сurrency)
		if err != nil {
			return err
		}
	}
	//если пользователя нет - добавляем его в БД
	if us.UsersId == 0 && us.ChatId == 0 {
		_, err = db.Exec("insert into notif_currency chatId values ?", chatId)
		if err != nil {
			return err
		}
	}
	return nil
}
func changeInformation(chatId int, currency string, db *sql.DB) error {
	_, err := db.Exec("update notif_currency set Сurrency = ? where chatId = ?", currency, chatId)
	return err
}

// добавить обротку на отдельные валюты
func allChatIdCostDB(Currency string, db *sql.DB) ([]userСurrency, error) {
	us := make([]userСurrency, 0)
	rows, err := db.Query("select chatId , currency from notif_currency where currency = ?", Currency)
	if err != nil {
		return us, err
	}
	i := 0
	for rows.Next() {
		us = append(us, userСurrency{0, 0, ""})
		err = rows.Scan(&us[i].ChatId, &us[i].Сurrency)
		i++
		if err != nil {
			return us, err
		}
	}
	err = rows.Close()
	if err != nil {
		return us, err
	}

	return us, nil
}
func deleteUserChatIdCostDB(chatId int, db *sql.DB) error {
	_, err := db.Exec("delete from notif_currency where chatId = ?", chatId)
	return err
}
func allUserChatIdCostDB(db *sql.DB) ([]int64, error) {
	us := make([]int64, 0)
	rows, err := db.Query("select DISTINCT ChatId from notif_currency")
	if err != nil {
		return us, err
	}
	i := 0
	for rows.Next() {
		us = append(us, 0)
		err = rows.Scan(&us[i])
		i++
		if err != nil {
			return us, err
		}
	}
	err = rows.Close()
	return us, err
}

// база данных для уведомлении о дохожнении валюты до заданной цены
func addUserChangeCostDB(chatId int, Currency string, cost float64, db *sql.DB) error {
	_, err := db.Exec("insert into notif_change_cost (СhatId, Currency,Change_cost) values (?,?)", chatId, Currency, cost)
	return err
}

// добавить обротку на отдельные валюты ??? возможно пробрасывать валюту не нужно , просто сравнивать два название (но тогда два метода )
func allChatIdChangeCostDB(currency string, cost float64, db *sql.DB) ([]userCost, error) {
	us := make([]userCost, 0)
	rows, err := db.Query("select ChatId from notif_change_cost where (Currency = ? AND Change_cost > ?) ", currency, cost)
	if err != nil {
		return us, err
	}
	i := 0
	for rows.Next() {
		us = append(us, userCost{
			ChatId:   0,
			Сurrency: "",
			Cost:     0,
		})
		err = rows.Scan(&us[i].ChatId, &us[i].Сurrency, &us[i].Cost)
		i++
		if err != nil {
			return us, err
		}
	}
	//удаление всех кого касался запрос
	err = rows.Close()
	if err != nil {
		return us, err
	}
	return us, nil
}
func deleteUserChatIdChangeCostDB(chatId int, db *sql.DB) error {
	_, err := db.Exec("delete from notif_change_cost where chatId = ?", chatId)
	return err
}
func allUserChatIdChangeCostDB(db *sql.DB) ([]int64, error) {
	us := make([]int64, 0)
	rows, err := db.Query("select DISTINCT ChatId from notif_change_cost ")
	if err != nil {
		return us, err
	}
	i := 0
	for rows.Next() {
		us = append(us, 0)
		err = rows.Scan(&us[i])
		i++
		if err != nil {
			return us, err
		}
	}
	err = rows.Close()
	return us, err
}
