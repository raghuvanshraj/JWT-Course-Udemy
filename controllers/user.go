package controllers

import (
	"database/sql"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"jwt-course/models"
	"jwt-course/utils"
	"net/http"
)

func (c Controller) SignUp() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		var user models.User
		var responseError models.Error
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			c.logger.Fatalln(err)
		}

		if user.Email == "" {
			responseError.Message = "email is missing"
			utils.RespondWithError(rw, http.StatusBadRequest, responseError)
			return
		}
		if user.Password == "" {
			responseError.Message = "password is missing"
			utils.RespondWithError(rw, http.StatusBadRequest, responseError)
			return
		}

		passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
		if err != nil {
			c.logger.Fatalln(err)
		}

		c.logger.Println("password hash generated")

		user.Password = string(passwordHash)
		user, err = c.userRepositoryObj.SignUp(user)
		if err != nil {
			responseError.Message = "server error"
			utils.RespondWithError(rw, http.StatusInternalServerError, responseError)
			return
		}

		user.Password = ""
		err = json.NewEncoder(rw).Encode(user)
		if err != nil {
			c.logger.Fatalln(err)
		}
	}
}

func (c Controller) Login() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		var (
			user          models.User
			responseJWT   models.JWT
			responseError models.Error
		)

		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			c.logger.Fatalln(err)
		}

		if user.Email == "" {
			responseError.Message = "email is missing"
			utils.RespondWithError(rw, http.StatusBadRequest, responseError)
			return
		}
		if user.Password == "" {
			responseError.Message = "password is missing"
			utils.RespondWithError(rw, http.StatusBadRequest, responseError)
			return
		}

		password := user.Password

		user, err = c.userRepositoryObj.Login(user)
		if err != nil {
			if err == sql.ErrNoRows {
				responseError.Message = "specified user does not exist"
				utils.RespondWithError(rw, http.StatusBadRequest, responseError)
			} else {
				c.logger.Fatalln(err)
			}
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
			responseError.Message = "password is invalid"
			utils.RespondWithError(rw, http.StatusUnauthorized, responseError)
			return
		}

		tokenString, err := utils.GenerateToken(user)
		if err != nil {
			c.logger.Fatalln(err)
		}

		responseJWT.Token = tokenString
		err = json.NewEncoder(rw).Encode(responseJWT)
		if err != nil {
			c.logger.Fatalln(err)
		}
	}
}
