package plan

import (
	"database/sql"
	"fmt"

	"github.com/JorgeEmanoel/money-keeper-backend/model"
)

type PlanRepository struct {
	Db *sql.DB
}

func MakePlanRepository(db *sql.DB) *PlanRepository {
	return &PlanRepository{
		Db: db,
	}
}

func (r *PlanRepository) Store(name, description, status string, ownerId int) (int, error) {
	result, err := r.Db.Exec("INSERT INTO plans (name, description, status, user_id) VALUES (?, ?, ?, ?)", name, description, status, ownerId)
	if err != nil {
		return 0, err
	}

	id, _ := result.LastInsertId()

	return int(id), nil
}

func (r *PlanRepository) Delete(id int) error {
	result, err := r.Db.Exec("DELETE FROM plans WHERE id = ?", id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows < 1 {
		return fmt.Errorf("no plans found with specified id %d", id)
	}

	return nil
}

func (r *PlanRepository) GetByUserId(userId int) ([]model.Plan, error) {
	result, err := r.Db.Query("SELECT id, name, description, status FROM plans WHERE user_id = ?", userId)
	if err != nil {
		return []model.Plan{}, err
	}

	var plans []model.Plan

	for result.Next() {
		var plan model.Plan
		result.Scan(&plan.Id, &plan.Name, &plan.Description, &plan.Status)

		plans = append(plans, plan)
	}

	return plans, nil
}

func (r *PlanRepository) FirstByUserId(userId int) (model.Plan, error) {
	result, err := r.Db.Query("SELECT id, name, description, status FROM plans WHERE user_id = ? LIMIT 1", userId)
	if err != nil {
		return model.Plan{}, err
	}

	var plan model.Plan
	result.Scan(&plan.Id, &plan.Name, &plan.Description, &plan.Status)

	return plan, nil
}
