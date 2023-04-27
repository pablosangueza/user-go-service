package users

import (
	"context"

	"github.com/jmoiron/sqlx"
)

const (
	addUserCmd = `INSERT INTO gouser ("user_name", "last_name", "email", "role")
		VALUES ($1, $2, $3, $4) RETURNING user_id;`

	updateUserCmd = `UPDATE gouser SET user_name=$1, last_name=$2, email=$3, role=$4  WHERE user_id=$5;`

	deleteUserCmd = `DELETE FROM gouser  WHERE user_id=$1;`
)

type SaveUserCommand interface {
	SaveUser(context.Context, User) (int, error)
	UpdateUser(context.Context, User) (int64, error)
	DeleteUser(context.Context, int) (int64, error)
}

type saveUserCommand struct {
	db *sqlx.DB
}

func NewUserCommand(db *sqlx.DB) SaveUserCommand {
	return &saveUserCommand{
		db: db,
	}
}

// SaveUser implements SaveUserCommand
func (s *saveUserCommand) SaveUser(ctx context.Context, user User) (int, error) {
	var user_id int
	err := s.db.QueryRowxContext(ctx, addUserCmd, user.UserName, user.LastName, user.Email, user.Role).Scan(&user_id)
	if err != nil {
		return user_id, err
	}

	return user_id, nil
}

// UpdateUser implements SaveUserCommand
func (s *saveUserCommand) UpdateUser(ctx context.Context, user User) (int64, error) {

	res, err := s.db.ExecContext(ctx, updateUserCmd, user.UserName, user.LastName, user.Email, user.Role, user.UserId)
	if err != nil {
		return 0, err
	}
	rows, err := res.RowsAffected()

	return rows, nil

}

// DeleteUser implements SaveUserCommand
func (s *saveUserCommand) DeleteUser(ctx context.Context, userId int) (int64, error) {

	res, err := s.db.ExecContext(ctx, deleteUserCmd, userId)
	if err != nil {
		return 0, err
	}
	rows, err := res.RowsAffected()
	return rows, nil

}
