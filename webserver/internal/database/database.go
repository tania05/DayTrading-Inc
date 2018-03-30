package database

import (
	"webserver/internal/money"
	"webserver/internal/logger"
	"webserver/internal/context"
	"common/config"
	"sync"

	"database/sql"
	_ "github.com/lib/pq"
	"fmt"
	"time"
	"os"
	"strconv"
	"hash/fnv"
)

type Holding interface {
	pay(queryable, *context.Context) error
	receive(queryable, *context.Context) error
}

type Transaction struct {
	payable    Holding
	receivable Holding
}

type StockHolding struct {
	UserId      string
	StockSymbol string
	Amount      int
}

type MoneyHolding struct {
	UserId string
	Amount money.Money
}

type queryable interface {
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
}

var mutex sync.Mutex
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

func getDatabase(userId string) *sql.DB {
	digest := fnv.New32a()
	digest.Write([]byte(userId))

	hash := digest.Sum32()
	hashInt := int(hash)
	if hashInt < 0 {
		hashInt = hashInt * -1
	}
	index := hashInt % len(databases)
	fmt.Printf("User has hased to database %d\n", index)
	return databases[index]
}

func init() {
	waitForConnection()
}

func (hold StockHolding) pay(target queryable, ctx *context.Context) error {
	row := target.QueryRow(`
		UPDATE stocks
			SET amount = amount - $3
			WHERE user_id = $1
				AND stock_sym = $2
			RETURNING amount;
	`, hold.UserId, hold.StockSymbol, hold.Amount)

	var newAmount int
	err := row.Scan(&newAmount)

	if err != nil {
		return ctx.MakeError(err.Error())
	}
	return nil
}

func (hold StockHolding) receive(target queryable, ctx *context.Context) error {
	row := target.QueryRow(`
		INSERT INTO stocks(user_id, stock_sym, amount)
			VALUES ($1, $2, $3)
			ON CONFLICT(user_id, stock_sym) DO UPDATE
				SET amount = stocks.amount + $3
			RETURNING amount;
	`, hold.UserId, hold.StockSymbol, hold.Amount)

	var newAmount int
	err := row.Scan(&newAmount)
	if err != nil {
		return ctx.MakeError(err.Error())
	}
	return nil
}

func (hold MoneyHolding) pay(target queryable, ctx *context.Context) error {
	row := target.QueryRow(`
		UPDATE users
			SET money = money - $2
			WHERE id = $1
			RETURNING money;
	`, hold.UserId, int(hold.Amount))

	var newBalance int
	err := row.Scan(&newBalance)

	if err != nil {
		return ctx.MakeError(err.Error())
	}
	ctx.Funds = money.Money(newBalance)
	ctx.MakeAccountTransactionLog(logger.RemoveAction)
	return nil

}

func (hold MoneyHolding) receive(target queryable, ctx *context.Context) error {
	row := target.QueryRow(`
		INSERT INTO users(id,money)
			VALUES ($1, $2)
			ON CONFLICT(id) DO UPDATE
				SET money = users.money + $2
			RETURNING money;
	`, hold.UserId, int(hold.Amount))

	var newBalance int
	err := row.Scan(&newBalance)
	if err != nil {
		return ctx.MakeError(err.Error())
	}
	ctx.Funds = money.Money(newBalance)
	ctx.MakeAccountTransactionLog(logger.AddAction)
	return nil
}

func AddFunds(ctx *context.Context, amount money.Money) error {
	receivable := MoneyHolding{UserId: ctx.UserId, Amount: amount}
	err := receivable.receive(getDatabase(ctx.UserId), ctx)
	return err
}

//TODO
func attemptAllocate(ctx *context.Context, trans Transaction) (Transaction, error) {

	tx, err := getDatabase(ctx.UserId).Begin()
	if err != nil {
		return Transaction{}, ctx.MakeError("Failed to initialize transaction context")
	}
	defer tx.Rollback()

	err = trans.payable.pay(getDatabase(ctx.UserId), ctx)
	if err != nil {
		return Transaction{}, ctx.MakeError(err.Error())
	}

	var moneyHolding MoneyHolding
	var stockHolding StockHolding
	var isBuy bool
	var ok bool

	//ehh this is bad
	//and we assume that both payable and recieve are the same person
	//and that one holding of each type
	//TODO
	switch p := trans.payable.(type) {
	case MoneyHolding:
		moneyHolding = p
		isBuy = true
		stockHolding, ok = trans.receivable.(StockHolding)
		if !ok {
			panic("Bad assumption")
		}
		break
	case StockHolding:
		stockHolding = p
		isBuy = false
		moneyHolding, ok = trans.receivable.(MoneyHolding)
		if !ok {
			panic("Bad assumption")
		}
	default:
		panic("Unknown holding type")
	}

	_, err = getDatabase(ctx.UserId).Exec(`
		INSERT INTO transactions(user_id, money_amount, stock_sym, stock_amount, is_buy, created_at)
			VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP)
	`, ctx.UserId, int(moneyHolding.Amount), stockHolding.StockSymbol, stockHolding.Amount, isBuy)

	if err != nil {
		return Transaction{}, ctx.MakeError(err.Error())
	}

	tx.Commit()

	return trans, nil
}

func AllocateFunds(ctx *context.Context, amount money.Money, stockAmount int) (Transaction, error) {
	if amount == 0 || stockAmount == 0 {
		return Transaction{}, ctx.MakeError("Both amount of funds and stock amount must be non-zero")
	}
	payable := MoneyHolding{UserId: ctx.UserId, Amount: amount}
	receivable := StockHolding{UserId: ctx.UserId, StockSymbol: ctx.StockSymbol, Amount: stockAmount}
	trans := Transaction{payable: payable, receivable: receivable}

	return attemptAllocate(ctx, trans)
}

func AllocateStocks(ctx *context.Context, stockAmount int, amount money.Money) (Transaction, error) {
	if amount == 0 || stockAmount == 0 {
		return Transaction{}, ctx.MakeError("Both amount of funds and stock amount must be non-zero")
	}
	payable := StockHolding{UserId: ctx.UserId, StockSymbol: ctx.StockSymbol, Amount: stockAmount}
	receivable := MoneyHolding{UserId: ctx.UserId, Amount: amount}
	trans := Transaction{payable: payable, receivable: receivable}

	return attemptAllocate(ctx, trans)
}

func HoldStocks(ctx *context.Context, amount int) (Holding, error) {
	hold := StockHolding{UserId: ctx.UserId, StockSymbol: ctx.StockSymbol, Amount: amount}

	err := hold.pay(getDatabase(ctx.UserId), ctx)
	if err != nil {
		return nil, err
	}

	return hold, nil
}

func HoldMoney(ctx *context.Context, amount money.Money) (Holding, error) {
	hold := MoneyHolding{UserId: ctx.UserId, Amount: amount}

	err := hold.pay(getDatabase(ctx.UserId), ctx)
	if err != nil {
		return nil, err
	}

	return hold, nil
}

//TODO return and execute should be atomic
func Return(ctx *context.Context, holds ... Holding) error {

	tx, err := getDatabase(ctx.UserId).Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, hold := range holds {
		h := Holding(hold)
		err := h.receive(tx, ctx)
		if err != nil {
			return err
		}
	}

	tx.Commit()
	return nil
}

func Commit(ctx *context.Context, trans Transaction) error {
	return trans.receivable.receive(getDatabase(ctx.UserId), ctx)
}

func Cancel(ctx *context.Context, trans Transaction) error {
	return trans.payable.receive(getDatabase(ctx.UserId), ctx)
}
