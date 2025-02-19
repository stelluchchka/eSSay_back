package services

import (
	"database/sql"
	"errors"
	"essay/src/internal/models"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Не удалось создать мок базы данных: %v", err)
	}
	defer db.Close()

	userService := NewUserService(db)

	// Тест успешного создания пользователя
	user := &models.User{
		Mail:     "test@example.com",
		Nickname: "testuser",
		Password: "password123",
	}

	// Мокируем успешный запрос на создание пользователя
	mock.ExpectExec(`INSERT INTO "user" \(mail, nickname, password\)`).
		WithArgs(user.Mail, user.Nickname, hashPassword(user.Password)).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = userService.CreateUser(user)
	assert.NoError(t, err)

	// Проверка на дублирование email
	mock.ExpectExec(`INSERT INTO "user" \(mail, nickname, password\)`).
		WithArgs(user.Mail, user.Nickname, hashPassword(user.Password)).
		WillReturnError(errors.New("duplicate key value violates unique constraint"))

	err = userService.CreateUser(user)
	assert.Equal(t, ErrDuplicateEmail, err)

	// Проверка на ошибку при выполнении запроса
	mock.ExpectExec(`INSERT INTO "user" \(mail, nickname, password\)`).
		WithArgs(user.Mail, user.Nickname, hashPassword(user.Password)).
		WillReturnError(errors.New("some error"))

	err = userService.CreateUser(user)
	assert.Error(t, err)
}

func TestAuthenticate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Не удалось создать мок базы данных: %v", err)
	}
	defer db.Close()

	userService := NewUserService(db)

	// Тест успешной аутентификации
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, password, is_moderator FROM "user" WHERE mail = $1`)).
		WithArgs("test@example.com").
		WillReturnRows(sqlmock.NewRows([]string{"id", "password", "is_moderator"}).
			AddRow(1, hashPassword("password123"), false))

	user, err := userService.Authenticate("test@example.com", "password123")
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, uint64(1), user.ID)

	// Тест неудачной аутентификации из-за неправильного пароля
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, password, is_moderator FROM "user" WHERE mail = $1`)).
		WithArgs("test@example.com").
		WillReturnRows(sqlmock.NewRows([]string{"id", "password", "is_moderator"}).
			AddRow(1, hashPassword("password123"), false))

	user, err = userService.Authenticate("test@example.com", "wrongpassword")
	assert.Equal(t, ErrInvalidCredentials, err)
	assert.Nil(t, user)

	// Тест неудачной аутентификации из-за отсутствия пользователя
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, password, is_moderator FROM "user" WHERE mail = $1`)).
		WithArgs("nonexistent@example.com").
		WillReturnError(sql.ErrNoRows)

	user, err = userService.Authenticate("nonexistent@example.com", "password123")
	assert.Equal(t, ErrInvalidCredentials, err)
	assert.Nil(t, user)
}

func TestGetUserByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Не удалось создать мок базы данных: %v", err)
	}
	defer db.Close()

	userService := NewUserService(db)

	// Тест успешного получения пользователя по ID
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT "id", "mail", "nickname", "password", "is_moderator", "count_checks" FROM "user" WHERE "id" = $1`)).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "mail", "nickname", "password", "is_moderator", "count_checks"}).
			AddRow(1, "test@example.com", "testuser", hashPassword("password123"), false, 10))

	user, err := userService.GetUserByID(1)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, uint64(1), user.ID)

	// Тест случая, когда пользователь не найден
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT "id", "mail", "nickname", "password", "is_moderator", "count_checks" FROM "user" WHERE "id" = $1`)).
		WithArgs(999).
		WillReturnError(sql.ErrNoRows)

	user, err = userService.GetUserByID(999)
	assert.NoError(t, err)
	assert.Nil(t, user)

	// Тест ошибки при выполнении запроса
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT "id", "mail", "nickname", "password", "is_moderator", "count_checks" FROM "user" WHERE "id" = $1`)).
		WithArgs(1).
		WillReturnError(errors.New("some error"))

	user, err = userService.GetUserByID(1)
	assert.Error(t, err)
	assert.Nil(t, user)
}
