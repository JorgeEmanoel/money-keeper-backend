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

func (r *TransactionRepository) Store(name, description, direction, reference, currency, status string, value, ownerId int) (int, error) {
	result, err := r.Db.Exec("INSERT INTO transactions (name, description, direction, reference, currency, status, value, user_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", name, description, direction, reference, currency, status, value, ownerId)
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
	result, err := r.Db.Query("SELECT id, name, description, direction, reference, currency, status, value FROM transactions WHERE user_id = ?", userId)
	if err != nil {
		return []model.Transaction{}, err
	}

	var transactions []model.Transaction

	for result.Next() {
		var s model.Transaction
		result.Scan(&s.Id, &s.Name, &s.Description, &s.Direction, &s.Reference, &s.Currency, &s.Status, &s.Value)

		transactions = append(transactions, s)
	}

	return transactions, nil
}

func (r *TransactionRepository) GetOutcomingByUserId(userId int) ([]model.Transaction, error) {
	result, err := r.Db.Query("SELECT id, name, description, direction, reference, currency, status, value FROM transactions WHERE user_id = ? AND direction = 'outcoming'", userId)
	if err != nil {
		return []model.Transaction{}, err
	}

	var transactions []model.Transaction

	for result.Next() {
		var s model.Transaction
		result.Scan(&s.Id, &s.Name, &s.Description, &s.Direction, &s.Reference, &s.Currency, &s.Status, &s.Value)

		transactions = append(transactions, s)
	}

	return transactions, nil
}

func (r *TransactionRepository) GetIncomingByUserId(userId int) ([]model.Transaction, error) {
	result, err := r.Db.Query("SELECT id, name, description, direction, reference, currency, status, value FROM transactions WHERE user_id = ? AND direction = 'incoming'", userId)
	if err != nil {
		return []model.Transaction{}, err
	}

	var transactions []model.Transaction

	for result.Next() {
		var s model.Transaction
		result.Scan(&s.Id, &s.Name, &s.Description, &s.Direction, &s.Reference, &s.Currency, &s.Status, &s.Value)

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

func (r *TransactionRepository) GetByUserIdFromReference(userId int, reference string) ([]model.Transaction, error) {
	result, err := r.Db.Query("SELECT id, name, description, direction, reference, currency, status, value FROM transactions WHERE user_id = ? AND reference = '?'", userId, reference)
	if err != nil {
		return []model.Transaction{}, err
	}

	var transactions []model.Transaction

	for result.Next() {
		var s model.Transaction
		result.Scan(&s.Id, &s.Name, &s.Description, &s.Direction, &s.Reference, &s.Currency, &s.Status, &s.Value)

		transactions = append(transactions, s)
	}

	return transactions, nil
}

func (r *TransactionRepository) CountByUserIdFromReference(userId int, reference string) (int, error) {
	result, err := r.Db.Query("SELECT count(*) as total FROM transactions WHERE user_id = ? AND reference = '?'", userId, reference)
	if err != nil {
		return 0, err
	}

	total := 0

	result.Scan(&total)

	return total, nil
}
