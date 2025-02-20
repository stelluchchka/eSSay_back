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

func (s *UserService) GetUserInfoByID(id uint64) (*models.UserInfo, error) {
	user := &models.UserInfo{}

	query := `SELECT 
	u.id, u.mail, u.nickname, u.is_moderator, u.count_checks,
	COUNT(e.id) AS count_essays, 
	COUNT(CASE WHEN e.is_published THEN 1 END) AS count_published_essays,
	COALESCE(AVG(r.sum_score), 0) AS average_result
	FROM "user" u
	LEFT JOIN essay e ON u.id = e.user_id
	LEFT JOIN result r ON e.id = r.essay_id
	WHERE u.id = $1
	GROUP BY u.id, u.mail, u.nickname, u.is_moderator, u.count_checks`

	err := s.DB.QueryRow(query, id).Scan(
		&user.ID, &user.Mail, &user.Nickname, &user.IsModerator,
		&user.CountChecks, &user.CountEssays, &user.CountPublishedEssays, &user.AverageResult)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

// func (s *UserService) GetUserByID(id uint64) (*models.User, error) {
// 	user := &models.User{}

// 	query := `SELECT
// 	u.id, u.mail, u.nickname, u.is_moderator, u.count_checks
// 	FROM "user" u
// 	WHERE u.id = $1`

// 	err := s.DB.QueryRow(query, id).Scan(
// 		&user.ID, &user.Mail, &user.Nickname, &user.IsModerator,
// 		&user.CountChecks)

// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return nil, nil
// 		}
// 		return nil, err
// 	}

// 	return user, nil
// }

func (s *UserService) GetUsersCount() (int, error) {
	var count int

	query := `SELECT COUNT(*) FROM "user"`
	err := s.DB.QueryRow(query).Scan(&count)

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (s *UserService) GetNickname(mail string) (string, error) {
	var nickname string

	query := `SELECT nickname FROM "user" WHERE mail = $1`
	err := s.DB.QueryRow(query, mail).Scan(&nickname)

	if err != nil {
		return "", err
	}

	return nickname, nil
}

func (s *UserService) CreateUser(user *models.User) error {
	hashedPassword := hashPassword(user.Password)

	query := `INSERT INTO "user" (mail, nickname, "password") VALUES ($1, $2, $3)`
	_, err := s.DB.Exec(query, user.Mail, user.Nickname, hashedPassword)
	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			return ErrDuplicateEmail
		}
		return err
	}

	return nil
}

func (s *UserService) UpdateUser(mail string, nickname string, id uint64) error {
	query := `UPDATE "user" SET mail = $1, nickname = $2 WHERE id = $3`
	_, err := s.DB.Exec(query, mail, nickname, id)
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
