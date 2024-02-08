package models

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func StartDB() error {
	db, err := sql.Open("sqlite3", "./inventory.db")
	if err != nil {
		return err
	}

	DB = db
	return nil
}

type Item struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	FV       string `json:"type"`
	Quantity int    `json:"quantity"`
}

func GetItems() ([]Item, error) {
	rows, err := DB.Query("SELECT * FROM ITEMS")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	items := make([]Item, 0)

	for rows.Next() {
		item := Item{}
		err = rows.Scan(&item.Id, &item.Name, &item.FV, &item.Quantity)

		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	err = rows.Err()

	if err != nil {
		return nil, err
	}

	return items, err
}

func GetItemById(id int) (Item, error) {
	statement, err := DB.Prepare("SELECT * FROM ITEMS WHERE ID = ?")

	if err != nil {
		return Item{}, err
	}

	item := Item{}

	sqlErr := statement.QueryRow(id).Scan(&item.Id, &item.Name, &item.FV, &item.Quantity)

	if sqlErr != nil {
		if sqlErr == sql.ErrNoRows {
			return Item{}, nil
		}
		return Item{}, sqlErr
	}
	return item, nil
}

func AddItem(item Item) (bool, error) {
	transaction, err := DB.Begin()

	if err != nil {
		return false, err
	}

	statement, err := transaction.Prepare("INSERT INTO ITEMS (NAME, FV, QUANTITY) VALUES (?, ?, ?))")

	if err != nil {
		return false, err
	}

	defer statement.Close()

	_, err = statement.Exec(item.Name, item.FV, item.Quantity)

	if err != nil {
		return false, err
	}

	transaction.Commit()

	return true, nil
}

func UpdateQuantity(id int, amount int) (bool, error) {
	transaction, err := DB.Begin()

	if err != nil {
		return false, err
	}

	statement, err := transaction.Prepare("UPDATE ITEMS SET QUANTITY = ? WHERE ID = ?")

	if err != nil {
		return false, err
	}

	defer statement.Close()

	_, err = statement.Exec(amount, id)

	if err != nil {
		return false, err
	}

	transaction.Commit()

	return true, nil
}

func RemoveItem(id int) (bool, error) {
	transaction, err := DB.Begin()

	if err != nil {
		return false, err
	}

	statement, err := transaction.Prepare("DELETE FROM ITEMS WHERE ID = ?")

	if err != nil {
		return false, err
	}

	defer statement.Close()

	_, err = statement.Exec(id)

	if err != nil {
		return false, err
	}

	transaction.Commit()

	return true, nil
}
