package repositoryOperetor

import (
	"database/sql"
	"errors"
	"fmt"
	"userServer/internal/models/modelOperator"

	_ "github.com/lib/pq"
)

type operetorRepository struct {
	db *sql.DB
}

func NewOperetorRepository() *operetorRepository {
	connStr := "host=localhost port=5432 user=postgres password=1234  dbname=postgres sslmode=disable" //TODO ПИЗДЕЦ
	db, _ := sql.Open("postgres", connStr)

	return &operetorRepository{db: db}
}

func (r *operetorRepository) Get(id int64) (*modelOperator.Operator, error) {
	query := `SELECT name, is_active FROM operators WHERE id = $1`
	operator := &modelOperator.Operator{}

	err := r.db.QueryRow(query, id).Scan(
		&operator.Name,
		&operator.Is_active,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("error scanning plan: %w", err)
	}

	operator.ID = id
	return operator, nil
}

func (r *operetorRepository) GetAll() ([]*modelOperator.Operator, error) {
	query := `SELECT id, name, is_active FROM operators`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, nil //TODO
	}
	defer rows.Close()

	var operators []*modelOperator.Operator

	for rows.Next() {
		operator := &modelOperator.Operator{}
		err := rows.Scan(&operator.ID, &operator.Name, &operator.Is_active)
		if err != nil {
			return nil, nil //TODO
		}
		operators = append(operators, operator)
	}

	if err := rows.Err(); err != nil {
		return nil, nil //TODO
	}
	return operators, nil
}
