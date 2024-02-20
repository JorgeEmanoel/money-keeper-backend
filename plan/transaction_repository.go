package plan

import (
	"database/sql"
	"fmt"

	"github.com/JorgeEmanoel/money-keeper-backend/model"
)

type TransactionRepository struct {
	Db *sql.DB
}

func MakeTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{
		Db: db,
	}
}

func (r *TransactionRepository) Store(name, description, direction, period, currency, status string, value float64, ownerId int) (int, error) {
	result, err := r.Db.Exec("INSERT INTO transactions (name, description, direction, period, currency, status, value, user_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", name, description, direction, period, currency, status, value, ownerId)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (r *TransactionRepository) Delete(id int) error {
	result, err := r.Db.Exec("DELETE FROM transactions WHERE id = ?", id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows < 1 {
		return fmt.Errorf("no transactions found with specified id %d", id)
	}

	return nil
}

func (r *TransactionRepository) GetByUserId(userId int) ([]model.Transaction, error) {
	result, err := r.Db.Query("SELECT id, name, description, direction, period, currency, status, value FROM transactions WHERE user_id = ?", userId)
	if err != nil {
		return []model.Transaction{}, err
	}

	var transactions []model.Transaction

	for result.Next() {
		var s model.Transaction
		result.Scan(&s.Id, &s.Name, &s.Description, &s.Direction, &s.Period, &s.Currency, &s.Status, &s.Value)

		transactions = append(transactions, s)
	}

	return transactions, nil
}

func (r *TransactionRepository) GetOutcomingByUserId(userId int, period string) ([]model.Transaction, error) {
	result, err := r.Db.Query("SELECT id, name, description, direction, period, currency, status, value FROM transactions WHERE user_id = ? AND direction = 'outcome' AND period = ?", userId, period)
	if err != nil {
		return []model.Transaction{}, err
	}

	var transactions []model.Transaction

	for result.Next() {
		var s model.Transaction
		result.Scan(&s.Id, &s.Name, &s.Description, &s.Direction, &s.Period, &s.Currency, &s.Status, &s.Value)

		transactions = append(transactions, s)
	}

	return transactions, nil
}

func (r *TransactionRepository) GetIncomingByUserId(userId int, period string) ([]model.Transaction, error) {
	result, err := r.Db.Query("SELECT id, name, description, direction, period, currency, status, value FROM transactions WHERE user_id = ? AND direction = 'income' AND period = ?", userId, period)
	if err != nil {
		return []model.Transaction{}, err
	}

	var transactions []model.Transaction

	for result.Next() {
		var s model.Transaction
		result.Scan(&s.Id, &s.Name, &s.Description, &s.Direction, &s.Period, &s.Currency, &s.Status, &s.Value)

		transactions = append(transactions, s)
	}

	return transactions, nil
}

func (r *TransactionRepository) ChangeStatus(id int, status string) error {
	_, err := r.Db.Exec("UPDATE transactions SET status = ? WHERE id = ?", status, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *TransactionRepository) GetByUserIdFromPeriod(userId int, period string) ([]model.Transaction, error) {
	stmt, _ := r.Db.Prepare("SELECT id, name, description, direction, period, currency, status, value FROM transactions WHERE user_id = ? AND period = ?")
	result, err := stmt.Query(userId, period)
	if err != nil {
		return []model.Transaction{}, err
	}

	var transactions []model.Transaction

	for result.Next() {
		var s model.Transaction
		result.Scan(&s.Id, &s.Name, &s.Description, &s.Direction, &s.Period, &s.Currency, &s.Status, &s.Value)

		transactions = append(transactions, s)
	}

	return transactions, nil
}

func (r *TransactionRepository) CountByUserIdFromPeriod(userId int, period string) (int, error) {
	var total int

	stmt, _ := r.Db.Prepare("SELECT COUNT(*) FROM transactions WHERE user_id = ? AND period = ?")
	err := stmt.QueryRow(userId, period).Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}
