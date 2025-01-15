package services

import (
	"database/sql"
	"essay/src/internal/database"
	"essay/src/internal/models"
	"log"
)

type UserService struct {
	DB *sql.DB
}

func NewUserService() *UserService {
	return &UserService{
		DB: database.GetPostgreSQLConnection(),
	}
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
	query := `INSERT INTO "user" (mail, nickname, password) VALUES ($1, $2, $3)`
	_, err := s.DB.Exec(query, user.Mail, user.Nickname, user.Password)
	if err != nil {
		log.Println("Error creating user:", err)
		return err
	}
	return nil
}

func (s *UserService) Close() {
	if err := s.DB.Close(); err != nil {
		log.Println("Error closing database connection:", err)
	}
}
