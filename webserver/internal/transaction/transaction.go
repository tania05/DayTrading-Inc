package transaction

import (
	"common/money"
	"common/logger"
	"common/context"
	_ "github.com/lib/pq"
	"fmt"
	"common/database"
)

type Holding interface {
	pay(database.Queryable, *context.Context) error
	receive(database.Queryable, *context.Context) error
}

type Transaction struct {
	Id         int
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

func (hold StockHolding) pay(target database.Queryable, ctx *context.Context) error {
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

func (hold StockHolding) receive(target database.Queryable, ctx *context.Context) error {
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

func (hold MoneyHolding) pay(target database.Queryable, ctx *context.Context) error {
	row := target.QueryRow(`
		UPDATE users
			SET money = money - $2
			WHERE Id = $1
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

func (hold MoneyHolding) receive(target database.Queryable, ctx *context.Context) error {
	row := target.QueryRow(`
		INSERT INTO users(Id,money)
			VALUES ($1, $2)
			ON CONFLICT(Id) DO UPDATE
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
	err := receivable.receive(database.GetDatabase(ctx.UserId), ctx)
	return err
}

//TODO
func attemptAllocate(ctx *context.Context, trans Transaction) (Transaction, error) {
	fmt.Println("Attempting allocation ", ctx, " trans ", trans)
	tx, err := database.GetDatabase(ctx.UserId).Begin()
	if err != nil {
		return Transaction{}, ctx.MakeError("Failed to initialize transaction context")
	}
	defer tx.Rollback()

	err = trans.payable.pay(tx, ctx)
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

	row := tx.QueryRow(`
		INSERT INTO transactions(user_id, money_amount, stock_sym, stock_amount, is_buy, created_at)
			VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP) RETURNING Id
	`, ctx.UserId, int(moneyHolding.Amount), stockHolding.StockSymbol, stockHolding.Amount, isBuy)

	var txId int
	err = row.Scan(&txId)
	if err != nil {
		return Transaction{}, ctx.MakeError(err.Error())
	}

	trans.Id = txId
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

	err := hold.pay(database.GetDatabase(ctx.UserId), ctx)
	if err != nil {
		return nil, err
	}

	return hold, nil
}

func HoldMoney(ctx *context.Context, amount money.Money) (Holding, error) {
	hold := MoneyHolding{UserId: ctx.UserId, Amount: amount}

	err := hold.pay(database.GetDatabase(ctx.UserId), ctx)
	if err != nil {
		return nil, err
	}

	return hold, nil
}

func CommitTransaction(ctx *context.Context, isBuy bool) error {
	return commitOrCancelTransaction(ctx, isBuy, true)
}

func CancelTransaction(ctx *context.Context, isBuy bool) error {
	return commitOrCancelTransaction(ctx, isBuy, false)
}

func CancelByTimeout(ctx *context.Context, txId int) error {
	tx, err := database.GetDatabase(ctx.UserId).Begin()
	if err != nil {
		return ctx.MakeError("Failed to initalize transaction")
	}
	defer tx.Rollback()

	row := tx.QueryRow(`
		DELETE FROM transactions
		  WHERE transactions.Id = $1
		RETURNING money_amount, stock_sym, stock_amount, is_buy
	`, txId)

	var moneyAmount int
	var stockSym string
	var stockAmount int
	var isBuy bool
	err = row.Scan(&moneyAmount, &stockSym, &stockAmount, &isBuy)
	if err != nil {
		return nil // just means this was finshed by normal means
	}

	moneyHolding := MoneyHolding{UserId: ctx.UserId, Amount: money.Money(moneyAmount)}
	stockHolding := StockHolding{UserId: ctx.UserId, StockSymbol: stockSym, Amount: stockAmount}

	if !isBuy { //cancelling a sell or commiting a buy
		err = stockHolding.receive(tx, ctx)
		if err != nil {
			return err
		}
	} else {
		err = moneyHolding.receive(tx, ctx)
		if err != nil {
			return err
		}
	}

	tx.Commit()
	return nil

}

func commitOrCancelTransaction(ctx *context.Context, isBuy bool, isCommit bool) error {
	tx, err := database.GetDatabase(ctx.UserId).Begin()
	if err != nil {
		return ctx.MakeError("Failed to initalize transaction")
	}
	defer tx.Rollback()

	row := tx.QueryRow(`
		WITH subquery AS (
			SELECT Id AS Id
			FROM transactions
			WHERE user_id = $1
				  AND is_buy = $2
			ORDER BY created_at DESC
			LIMIT 1
		)
		DELETE FROM transactions
		  WHERE transactions.Id = (SELECT Id from subquery)
		RETURNING money_amount, stock_sym, stock_amount
	`, ctx.UserId, isBuy)

	var moneyAmount int
	var stockSym string
	var stockAmount int
	err = row.Scan(&moneyAmount, &stockSym, &stockAmount)
	if err != nil {
		return ctx.MakeError("Failed to find recent transaction for user")
	}

	moneyHolding := MoneyHolding{UserId: ctx.UserId, Amount: money.Money(moneyAmount)}
	stockHolding := StockHolding{UserId: ctx.UserId, StockSymbol: stockSym, Amount: stockAmount}

	if isBuy == isCommit { //cancelling a sell or commiting a buy
		err = stockHolding.receive(tx, ctx)
		if err != nil {
			return err
		}
	} else {
		err = moneyHolding.receive(tx, ctx)
		if err != nil {
			return err
		}
	}

	tx.Commit()
	return nil
}

//TODO return and execute should be atomic
func Return(ctx *context.Context, holds ... Holding) error {

	tx, err := database.GetDatabase(ctx.UserId).Begin()
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
