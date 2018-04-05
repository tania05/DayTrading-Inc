package database

import (
	"common/config"
	"database/sql"
	_ "github.com/lib/pq"
	"fmt"
	"time"
	"os"
	"strconv"
	"common/hashing"
)


type Queryable interface {
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
}

var databases []*sql.DB
func waitForConnection() {
	databaseCount := os.Getenv("DATABASE_COUNT")
	count, err := strconv.Atoi(databaseCount)
	if err != nil {
		panic("Could not read environment variable")
	}

	databases = make([]*sql.DB, count)
	for i := 0; i < count; i++ {
		dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s host=%s%d port=%d sslmode=%s",
			config.GlobalConfig.Database.Username,
			config.GlobalConfig.Database.Password,
			config.GlobalConfig.Database.Database,
			config.GlobalConfig.Database.Domain, i,
			config.GlobalConfig.Database.Port,
			config.GlobalConfig.Database.SSLMode)

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

		fmt.Printf("database %d ready... ping was successful\n", i)

		databases[i] = db
	}
}

func GetDatabase(userId string) *sql.DB {
	index := hashing.ModuloHash(userId, len(databases))
	fmt.Printf("User has hased to database %d\n", index)
	return databases[index]
}
func init() {
	waitForConnection()
}
