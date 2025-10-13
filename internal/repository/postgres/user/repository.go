package user

import (
	"database/sql"
	"userServer/internal/model/user"

	_ "github.com/lib/pq"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *userRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user user.Model) error {
	query := `INSERT INTO users (id, username, operator_id) VALUES ($1, $2, $3) RETURNING id`
	return r.db.QueryRow(query, user.ID, user.Username, 0).Scan(&user.ID)
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

func (r *userRepository) Users() ([]*user.Model, error) {
	query := `SELECT id, username, create_at, operator_id, total_sum, is_trial FROM users`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*user.Model
	for rows.Next() {
		var user user.Model
		var operatorID sql.NullInt64

		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.CreateAt,
			&operatorID,
			&user.TotalSum,
			&user.IsTrial,
		)
		if err != nil {
			return nil, err
		}

		if operatorID.Valid {
			user.MobileOperatorID = operatorID.Int64
		} else {
			user.MobileOperatorID = 0
		}

		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) Update(user user.FullModel) error {
	query := `
        UPDATE users 
        SET username = $1, 
            create_at = $2, 
            mobile_operator_id = $3, 
            total_sum = $4, 
            is_trial = $5 
        WHERE id = $6`

	_, err := r.db.Exec(query,
		user.Username,
		user.CreateAt,
		user.MobileOperator.ID,
		user.TotalSum,
		user.IsTrial,
		user.ID,
	)

	return err
}
