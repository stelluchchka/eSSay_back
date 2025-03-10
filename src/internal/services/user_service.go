package services

import (
	"database/sql"
	"errors"
)

var ErrDuplicateEmail = errors.New("email already in use")
var ErrInvalidCredentials = errors.New("invalid email or password")
var ErrWrongID = errors.New("wrong id")
var ErrLikeAlreadyExists = errors.New("like already exists")
var ErrLikeNotFound = errors.New("like doesn't exists")

type UserService struct {
	DB *sql.DB
}

func NewUserService(db *sql.DB) *UserService {
	return &UserService{
		DB: db,
	}
}
