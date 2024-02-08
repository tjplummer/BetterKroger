package models

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func StartDB() error {
	db, err := sql.Open("sqlite3", "./bk.db")
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
	rows, err := DB.Query("SELECT * FROM INVENTORY")

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
	statement, err := DB.Prepare("SELECT * FROM INVENTORY WHERE CODE = ?")

	if err != nil {
		return Item{}, err
	}

	code, err = EnsureCode(code)

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

	statement, err := transaction.Prepare("INSERT INTO INVENTORY (CODE, NAME, PRICE) VALUES (?, ?, ?))")

	if err != nil {
		return false, err
	}

	defer statement.Close()

	item.Code, err = EnsureCode(item.Code)

	if err != nil {
		return false, err
	}

	convert := ConvertToInt(item.Price)

	_, err = statement.Exec(item.Code, item.Name, convert)

	if err != nil {
		return false, err
	}

	transaction.Commit()

	return true, nil
}

func UpdatePrice(code string, amount string) (bool, error) {
	transaction, err := DB.Begin()

	if err != nil {
		return false, err
	}

	statement, err := transaction.Prepare("UPDATE INVENTORY SET PRICE = ? WHERE CODE = ?")

	if err != nil {
		return false, err
	}

	defer statement.Close()

	code, err = EnsureCode(code)

	if err != nil {
		return false, err
	}

	convert := ConvertToInt(amount)

	_, err = statement.Exec(convert, code)

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

	statement, err := transaction.Prepare("DELETE FROM INVENTORY WHERE CODE = ?")

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

func ConvertToUSD(amount int) string {
	return fmt.Sprintf("$%v", amount)
}

func ConvertToInt(amount string) int {
	i, err := strconv.Atoi(amount[1:])

	// This is NOT good and should not stay
	if err != nil {
		panic(err)
	}

	return i
}

func EnsureCode(code string) (string, error) {
	// This just removes them if they are there. We will add them regardless
	s := strings.Trim(code, "-")
	i := 4

	if len(s) != 16 {
		return s, errors.New("Code is not the required 16 characters in length!")
	}

	// There is definitely more elegant ways to do this
	return (s[:i] + "-" + s[:i] + "-" + s[:i] + "-" + s[:i]), nil
}
