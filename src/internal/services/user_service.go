package services

import (
	"database/sql"
	"errors"
)

var (
	ErrDuplicateEmail     = errors.New("email already in use")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrWrongID            = errors.New("wrong id")
	ErrLikeAlreadyExists  = errors.New("like already exists")
	ErrLikeNotFound       = errors.New("like doesn't exists")
	ErrNoChecksLeft       = errors.New("no checks left")
)

type UserService struct {
	DB *sql.DB
}

func NewUserService(db *sql.DB) *UserService {
	return &UserService{
		DB: db,
	}
}
