package user

import (
	"database/sql"
	"userServer/internal/model/user"

	_ "github.com/lib/pq"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository() *userRepository {
	connStr := "host=localhost port=5432 user=postgres password=1234  dbname=postgres sslmode=disable" //TODO
	db, _ := sql.Open("postgres", connStr)

	return &userRepository{db: db}
}

func (r *userRepository) Create(user user.Model) error {
	query := `INSERT INTO users (id, username , operator_id) VALUES ($1, $2, $3) RETURNING id`
	return r.db.QueryRow(query, user.ID, user.Username, user.MobileOperatorID).Scan(&user.ID)
}

func (r *userRepository) User(id int64) (*user.FullModel, error) {
	query := `SELECT username, create_at, operator_id, total_sum, is_trial FROM users WHERE id = $1`
	user := &user.FullModel{}
	err := r.db.QueryRow(query, id).Scan(&user.Username, &user.CreateAt, &user.MobileOperator.ID, &user.TotalSum, &user.IsTrial)
	if err != nil { //TODO
		return nil, err //TODO
	}
	user.ID = id
	return user, nil
}

func (r *userRepository) Update(user user.Model) error {
	if user.Username != "" {
		query := `UPDATE users SET username = $1 WHERE id = $2`
		_, err := r.db.Exec(query, user.Username, user.ID)
		return err
	}
	if user.ID != 0 {
		query := `UPDATE users SET operator_id = $1 WHERE id = $2`
		_, err := r.db.Exec(query, user.MobileOperatorID, user.ID)
		return err
	}
	if !user.IsTrial {
		query := `UPDATE users SET is_trial = $1 WHERE id = $2`
		_, err := r.db.Exec(query, user.IsTrial, user.ID)
		return err
	}
	return nil // TODO
}
