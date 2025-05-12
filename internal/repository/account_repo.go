package repository

import (
	"banking-api/internal/models"
	"database/sql"
	"errors"
)

type AccountRepository struct {
	DB *sql.DB
}

func NewAccountRepository(db *sql.DB) *AccountRepository {
	return &AccountRepository{DB: db}
}

func (r *AccountRepository) CreateAccount(userID int64) (*models.Account, error) {
	query := `INSERT INTO accounts (user_id, balance) VALUES ($1, 0) RETURNING id, created_at`
	var account models.Account
	err := r.DB.QueryRow(query, userID).Scan(&account.ID, &account.CreatedAt)
	if err != nil {
		return nil, err
	}
	account.UserID = userID
	account.Balance = 0
	return &account, nil
}

func (r *AccountRepository) TopUpAccount(accountID, userID int64, amount float64) error {
	query := `UPDATE accounts SET balance = balance + $1 WHERE id = $2 AND user_id = $3`
	res, err := r.DB.Exec(query, amount, accountID, userID)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *AccountRepository) TransferFunds(fromID, toID, userID int64, amount float64) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Проверка владения и баланса
	var balance float64
	err = tx.QueryRow(`SELECT balance FROM accounts WHERE id = $1 AND user_id = $2`, fromID, userID).Scan(&balance)
	if err != nil {
		return err
	}
	if balance < amount {
		return errors.New("недостаточно средств")
	}

	// Списание
	_, err = tx.Exec(`UPDATE accounts SET balance = balance - $1 WHERE id = $2`, amount, fromID)
	if err != nil {
		return err
	}

	// Зачисление
	_, err = tx.Exec(`UPDATE accounts SET balance = balance + $1 WHERE id = $2`, amount, toID)
	if err != nil {
		return err
	}

	// Запись в транзакции
	_, err = tx.Exec(`
		INSERT INTO transactions (from_account_id, to_account_id, amount) 
		VALUES ($1, $2, $3)`,
		fromID, toID, amount)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *AccountRepository) GetFirstAccountByUserID(userID int64) (int64, error) {
	var accountID int64
	err := r.DB.QueryRow(`SELECT id FROM accounts WHERE user_id = $1 ORDER BY id LIMIT 1`, userID).Scan(&accountID)
	return accountID, err
}

func (r *AccountRepository) GetUserIDByAccountID(accountID int64) (int64, error) {
	var userID int64
	err := r.DB.QueryRow(`SELECT user_id FROM accounts WHERE id = $1`, accountID).Scan(&userID)
	return userID, err
}
