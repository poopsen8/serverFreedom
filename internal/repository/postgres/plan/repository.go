package plan

import (
	"database/sql"
	"userServer/internal/model/plan"

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

func (r *planRepository) Plan(id int64) (*plan.Model, error) {
	query := `SELECT name, duration, price, discount, is_private FROM plans WHERE id = $1`
	plan := &plan.Model{}

	err := r.db.QueryRow(query, id).Scan(
		&plan.Name,
		&plan.Duration,
		&plan.Price,
		&plan.Discount,
		&plan.IsPrivate,
	)
	if err != nil {
		return nil, err
	}

	plan.ID = id
	return plan, nil
}

func (r *planRepository) GetAll() ([]*plan.Model, error) {
	query := `SELECT id, name, duration, price, discount, is_private FROM plans`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err //TODO
	}
	defer rows.Close()

	var plans []*plan.Model

	for rows.Next() {
		plan := &plan.Model{}
		err := rows.Scan(&plan.ID, &plan.Name, &plan.Duration, &plan.Price, &plan.Discount, &plan.IsPrivate)
		if err != nil {
			return nil, nil //TODO
		}
		plans = append(plans, plan)
	}

	if err := rows.Err(); err != nil {
		return nil, err //TODO
	}
	return plans, nil
}
