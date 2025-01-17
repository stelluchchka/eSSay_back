package services

import (
	"crypto/sha256"
	"database/sql"
	"errors"
	"essay/src/internal/database"
	"essay/src/internal/models"
	"fmt"
	"log"
	"strings"
)

var ErrDuplicateEmail = errors.New("email already in use")

type UserService struct {
	DB *sql.DB
}

func NewUserService() *UserService {
	return &UserService{
		DB: database.GetPostgreSQLConnection(),
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
		log.Println("Error fetching user:", err)
		return nil, err
	}

	return user, nil
}

func (s *UserService) CreateUser(user *models.User) error {
	log.Printf("Attempting to create user: %s\n", user.Nickname)

	hashedPassword := hashPassword(user.Password)

	query := `INSERT INTO "user" (mail, nickname, password) VALUES ($1, $2, $3)`
	_, err := s.DB.Exec(query, user.Mail, user.Nickname, hashedPassword)
	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			log.Printf("Duplicate email detected for user %s: %v\n", user.Nickname, err)
			return ErrDuplicateEmail
		}
		log.Printf("Error creating user %s: %v\n", user.Nickname, err)
		return err
	}

	log.Printf("User %s successfully created\n", user.Nickname)
	return nil
}
