package users

type User struct {
	UserId   int    `db:"user_id"`
	UserName string `db:"user_name"`
	LastName string `db:"last_name"`
	Email    string `db:"email"`
	Role     string `db:"role"`
}
