package plan

import (
	"database/sql"
	"fmt"

	"github.com/JorgeEmanoel/money-keeper-backend/model"
)

type SkeletonRepository struct {
	Db *sql.DB
}

func MakeSkeletonRepository(db *sql.DB) *SkeletonRepository {
	return &SkeletonRepository{
		Db: db,
	}
}

func (r *SkeletonRepository) Store(name, description, behaviour, frequency, value, currency, planId, ownerId int) (int, error) {
	result, err := r.Db.Exec("INSERT INTO skeletons (name, description, behaviour, frequency, value, currency, plan_id, user_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", name, description, behaviour, frequency, value, currency, planId, ownerId)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (r *SkeletonRepository) Delete(id int) error {
	result, err := r.Db.Exec("DELETE FROM skeletons WHERE id = ?", id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows < 1 {
		return fmt.Errorf("no skeletons found with specified id %d", id)
	}

	return nil
}

func (r *SkeletonRepository) GetByUserId(userId int) ([]model.Skeleton, error) {
	result, err := r.Db.Query("SELECT id, name, description, behaviour, frequency, value, currency FROM skeletons WHERE user_id = ?", userId)
	if err != nil {
		return []model.Skeleton{}, err
	}

	var skeletons []model.Skeleton

	for result.Next() {
		var s model.Skeleton
		result.Scan(&s.Id, &s.Name, &s.Description, &s.Behaviour, &s.Frequency, &s.Value, &s.Currency)

		skeletons = append(skeletons, s)
	}

	return skeletons, nil
}

func (r *SkeletonRepository) ChangeStatus(id int, status string) error {
	_, err := r.Db.Exec("UPDATE skeletons SET status = ? WHERE id = ?", status, id)
	if err != nil {
		return err
	}

	return nil
}