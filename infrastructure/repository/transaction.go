package repository

import (
	"database/sql"
	"errors"

	"github.com/MrHenri/fullCycle-challenge/domain"
)

type TransactionRepositoryDb struct {
	db *sql.DB
}

func NewTransactionRepositoryDb(db *sql.DB) *TransactionRepositoryDb {
	return &TransactionRepositoryDb{db: db}
}

func (t *TransactionRepositoryDb) SaveTransaction(fromAccount domain.Account, toAccount domain.Account, transaction domain.Transaction) error {
	stmt, err := t.db.Prepare(`INSERT INTO transfers(id, from_account, to_account, amount) 
								VALUES ($1, $2, $3, $4)`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(
		transaction.ID,
		transaction.From,
		transaction.To,
		transaction.Amount,
	)

	if err != nil {
		return err
	}

	err = t.updateBalance(fromAccount, -transaction.Amount)
	if err != nil {
		return err
	}

	err = t.updateBalance(toAccount, transaction.Amount)
	if err != nil {
		return err
	}
	return nil
}

func (t *TransactionRepositoryDb) updateBalance(account domain.Account, amount float64) error {
	result := account.Balance + amount
	_, err := t.db.Exec(`UPDATE accounts SET balance = $1 WHERE id = $2`, result, account.ID)
	if err != nil {
		return err
	}

	return nil
}

func (t *TransactionRepositoryDb) GetAccount(accountNumber string) (domain.Account, error) {
	var a domain.Account
	stmt, err := t.db.Prepare("SELECT id, number, balance from accounts WHERE number = $1")

	if err != nil {
		return a, err
	}

	if err = stmt.QueryRow(accountNumber).Scan(&a.ID, &a.Number, &a.Balance); err != nil {
		return a, errors.New("account does not exist")
	}

	return a, nil
}

func (t *TransactionRepositoryDb) CreateAccount(account domain.Account) error {
	stmt, err := t.db.Prepare("INSERT INTO accounts (id, number, balance) VALUES($1, $2, $3)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(
		account.ID,
		account.Number,
		account.Balance,
	)

	if err != nil {
		return err
	}

	err = stmt.Close()

	if err != nil {
		return err
	}

	return nil
}
