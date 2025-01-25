package services

import (
	"crypto/sha256"
	"database/sql"
	"errors"
	"essay/src/internal/models"
	"fmt"
	"strings"
)

var ErrDuplicateEmail = errors.New("email already in use")
var ErrInvalidCredentials = errors.New("invalid email or password")

type UserService struct {
	DB *sql.DB
}

func NewUserService(db *sql.DB) *UserService {
	return &UserService{
		DB: db,
	}
}

func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return fmt.Sprintf("%x", hash[:])
}

func (s *UserService) GetUserByID(id int) (*models.User, error) {
	user := &models.User{}

	query := `SELECT "id", "mail", "nickname", "password", "is_moderator", "count_checks" FROM "user" WHERE "id" = $1`
	err := s.DB.QueryRow(query, id).Scan(&user.ID, &user.Mail, &user.Nickname, &user.Password, &user.IsModerator, &user.CountChecks)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

func (s *UserService) CreateUser(user *models.User) error {
	hashedPassword := hashPassword(user.Password)

	query := `INSERT INTO "user" (mail, nickname, password) VALUES ($1, $2, $3)`
	_, err := s.DB.Exec(query, user.Mail, user.Nickname, hashedPassword)
	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			return ErrDuplicateEmail
		}
		return err
	}

	return nil
}

func (s *UserService) Authenticate(mail, password string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, password, is_moderator FROM "user" WHERE mail = $1`
	err := s.DB.QueryRow(query, mail).Scan(&user.ID, &user.Password, &user.IsModerator)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if hashPassword(password) != user.Password {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}
