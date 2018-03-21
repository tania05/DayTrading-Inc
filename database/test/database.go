package main

import (
	"fmt"
	"time"
	"database/sql"
	_ "github.com/lib/pq"
)
func waitForConnection() *sql.DB {
	dbinfo := "user=postgres password=mysecretpass dbname=quote host=localhost port=5432 sslmode=disable"

	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		panic(err)
	}

	for {
		fmt.Printf("Beginning test ping using configuration %s\n", dbinfo)
		pingErr := db.Ping()
		if pingErr != nil {
			fmt.Println(pingErr)
			fmt.Println("Retrying in 1s")
			time.Sleep(time.Second * 1)
			continue
		}

		break
	}

	fmt.Println("database ready... ping was successful")

	return db
}

type queryable interface {
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
}

func moneyReceive(db queryable, id string, money int) (int, error) {
	row := db.QueryRow(`
		INSERT INTO users(id,money)
			VALUES ($1, $2)
			ON CONFLICT(id) DO UPDATE
				SET money = users.money + $2
			RETURNING money;
	`, id, money)

	var newBalance int
	err := row.Scan(&newBalance)
	return newBalance, err
}

func stockReceive(db queryable, userId string, stockSym string, amount int) (int, error) {
	row := db.QueryRow(`
		INSERT INTO stocks(user_id, stock_sym, amount)
			VALUES ($1, $2, $3)
			ON CONFLICT(user_id, stock_sym) DO UPDATE
				SET amount = stocks.amount + $3
			RETURNING amount;
	`, userId, stockSym, amount)

	var newAmount int
	err := row.Scan(&newAmount)
	return newAmount, err
}

func moneyPay(db queryable, id string, money int) (int, error) {
	row := db.QueryRow(`
		UPDATE users
			SET money = money - $2
			WHERE id = $1
			RETURNING money;
	`, id, money)

	var newBalance int
	err := row.Scan(&newBalance)
	return newBalance, err

}

func stockPay(db queryable, userId string, stockSym string, amount int) (int, error) {
	row := db.QueryRow(`
		UPDATE stocks
			SET amount = amount - $3
			WHERE user_id = $1
				AND stock_sym = $2
			RETURNING amount;
	`, userId, stockSym, amount)

	var newAmount int
	err := row.Scan(&newAmount)
	return newAmount, err
}

func main() {
	db := waitForConnection()

	v, err := stockReceive(db, "user_b", "DEF", 200)
	if err != nil {
		panic(err)
	}

	fmt.Printf("v %d\n", v)
}
