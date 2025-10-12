package operetor

import (
	"database/sql"
	"userServer/internal/model/operator"

	_ "github.com/lib/pq"
)

type operetorRepository struct {
	db *sql.DB
}

func NewOperetorRepository(db *sql.DB) *operetorRepository {
	return &operetorRepository{db: db}
}

func (r *operetorRepository) Operator(id int64) (*operator.Model, error) {
	query := `SELECT name, is_active FROM operators WHERE id = $1`
	operator := &operator.Model{}

	err := r.db.QueryRow(query, id).Scan(
		&operator.Name,
		&operator.Is_active,
	)
	if err != nil {
		return nil, err
	}

	operator.ID = id
	return operator, nil
}

func (r *operetorRepository) Operators() ([]*operator.Model, error) {
	query := `SELECT id, name, is_active FROM operators`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, nil //TODO
	}
	defer rows.Close()

	var operators []*operator.Model

	for rows.Next() {
		operator := &operator.Model{}
		err := rows.Scan(&operator.ID, &operator.Name, &operator.Is_active)
		if err != nil {
			return nil, err //TODO
		}
		operators = append(operators, operator)
	}

	if err := rows.Err(); err != nil {
		return nil, err //TODO
	}
	return operators, nil
}
