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
			fmt.Println("Retrying in 5s")
			time.Sleep(time.Second * 1)
			continue
		}

		break
	}

	fmt.Println("database ready... ping was successful")

	return db
}

/*
	db.Query(`
		INSERT INTO users(id,money)
			VALUES ($1, $2)
			ON CONFLICT(id) DO UPDATE
				SET money = users.money + $2
			RETURNING id,money;
	`, "abcdef", 400)

	row, err := db.Query(`
		INSERT INTO stocks(user_id, stock_sym, amount)
			VALUES ($1, $2, $3)
			ON CONFLICT(user_id, stock_sym) DO UPDATE
				SET amount = stocks.amount + $3
			RETURNING user_id, stock_sym, amount;
	`, "abcdef", "ABC", 120)

	row, err := db.Query(`
		UPDATE users
			SET money = money - $2
			WHERE id = $1
			RETURNING id,money;
	`, "abcdef", 250)

	row, err := db.Query(`
		UPDATE stocks
			SET amount = amount - $3
			WHERE user_id = $1
				AND stock_sym = $2
			RETURNING user_id,stock_sym,amount;
	`, "abcdef", "ABC", 200)
*/

func main() {
	db := waitForConnection()


	if err != nil {
		panic(err)
	}

	fmt.Println(row)
}
