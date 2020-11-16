package controllers

import (
	"database/sql"
	userRepository "jwt-course/repository/user"
	"log"
	"os"
)

type Controller struct{
	db *sql.DB
	logger *log.Logger
	userRepositoryObj userRepository.UserRepository
}

func NewController(db *sql.DB) Controller {
	return Controller{
		db: db,
		logger: log.New(os.Stdout, "controller: ", log.LstdFlags),
		userRepositoryObj: userRepository.NewUserRepository(db),
	}
}
