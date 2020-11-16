package userRepository

import (
	"database/sql"
	"jwt-course/models"
	"log"
	"os"
)

type UserRepository struct {
	db *sql.DB
	logger *log.Logger
}

func NewUserRepository(db *sql.DB) UserRepository {
	return UserRepository{
		db:     db,
		logger: log.New(os.Stdout, "user-repository: ", log.LstdFlags),
	}
}

func (u UserRepository) SignUp(user models.User) (models.User, error) {
	query := "insert into users (email, password) values ($1, $2) RETURNING id;"
	err := u.db.QueryRow(query, user.Email, user.Password).Scan(&user.ID)
	if err != nil {
		u.logger.Println(err)
		return user, err
	}

	return user, nil
}

func (u UserRepository) Login(user models.User) (models.User, error) {
	row := u.db.QueryRow("select * from users where email=$1", user.Email)
	err := row.Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		u.logger.Println(err)
		return user, err
	}

	return user, nil
}
