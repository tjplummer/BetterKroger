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
	Name  string `json:"name"`
	Code  string `json:"code"`
	Price string `json:"price"`
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
		err = rows.Scan(&item.Name, &item.Code, &item.Price)

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

func GetItemByCode(code string) (Item, error) {
	statement, err := DB.Prepare("SELECT * FROM ITEMS WHERE CODE = ?")

	if err != nil {
		return Item{}, err
	}

	item := Item{}

	sqlErr := statement.QueryRow(code).Scan(&item.Name, &item.Code, &item.Price)

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

	statement, err := transaction.Prepare("INSERT INTO ITEMS (NAME, CODE, PRICE) VALUES (?, ?, ?))")

	if err != nil {
		return false, err
	}

	defer statement.Close()

	_, err = statement.Exec(item.Name, item.Code, item.Price)

	if err != nil {
		return false, err
	}

	transaction.Commit()

	return true, nil
}

func UpdatePrice(id int, amount int) (bool, error) {
	transaction, err := DB.Begin()

	if err != nil {
		return false, err
	}

	statement, err := transaction.Prepare("UPDATE ITEMS SET PRICE = ? WHERE CODE = ?")

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

	statement, err := transaction.Prepare("DELETE FROM ITEMS WHERE CODE = ?")

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
