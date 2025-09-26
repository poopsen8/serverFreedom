package repositoryPlan

import (
	"database/sql"
	"errors"
	"fmt"
	"userServer/internal/models/modelPlan"

	_ "github.com/lib/pq"
)

type planRepository struct {
	db *sql.DB
}

func NewPlanRepository() *planRepository {
	connStr := "host=localhost port=5432 user=postgres password=1234  dbname=postgres sslmode=disable" //TODO ПИЗДЕЦ
	db, _ := sql.Open("postgres", connStr)

	return &planRepository{db: db}
}

func (r *planRepository) Get(id int64) (*modelPlan.Plan, error) {
	query := `SELECT name, duration, price, discount, is_private FROM plans WHERE id = $1`
	plan := &modelPlan.Plan{}

	err := r.db.QueryRow(query, id).Scan(
		&plan.Name,
		&plan.Duration,
		&plan.Price,
		&plan.Discount,
		&plan.IsPrivate,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("error scanning plan: %w", err)
	}

	plan.ID = id
	return plan, nil
}

func (r *planRepository) GetAll() ([]*modelPlan.Plan, error) {
	query := `SELECT id, name, duration, price, discount, is_private FROM plans`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, nil //TODO
	}
	defer rows.Close()

	var plans []*modelPlan.Plan

	for rows.Next() {
		plan := &modelPlan.Plan{}
		err := rows.Scan(&plan.ID, &plan.Name, &plan.Duration, &plan.Price, &plan.Discount, &plan.IsPrivate)
		if err != nil {
			return nil, nil //TODO
		}
		plans = append(plans, plan)
	}

	if err := rows.Err(); err != nil {
		return nil, nil //TODO
	}
	return plans, nil
}
