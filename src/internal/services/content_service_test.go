package services

import (
	"database/sql"
	"essay/src/internal/models"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestContentService_GetLikesCount(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	service := NewContentService(db)
	essayID := uint64(1)
	expectedCount := 5

	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM "like" WHERE essay_id = \$1`).
		WithArgs(essayID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

	count, err := service.GetLikesCount(essayID)

	assert.NoError(t, err)
	assert.Equal(t, expectedCount, count)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestContentService_AddLike_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	service := NewContentService(db)
	userID := uint64(1)
	essayID := uint64(1)

	mock.ExpectQuery(`
		INSERT INTO "like" \(user_id, essay_id\)
		VALUES \(\$1, \$2\)
		ON CONFLICT \(user_id, essay_id\) DO NOTHING
		RETURNING true`).
		WithArgs(userID, essayID).
		WillReturnRows(sqlmock.NewRows([]string{"true"}).AddRow(true))

	err := service.AddLike(userID, essayID)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestContentService_AddLike_AlreadyExists(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	service := NewContentService(db)
	userID := uint64(1)
	essayID := uint64(1)

	mock.ExpectQuery(`
		INSERT INTO "like" \(user_id, essay_id\)
		VALUES \(\$1, \$2\)
		ON CONFLICT \(user_id, essay_id\) DO NOTHING
		RETURNING true`).
		WithArgs(userID, essayID).
		WillReturnError(sql.ErrNoRows)

	err := service.AddLike(userID, essayID)

	assert.Equal(t, ErrLikeAlreadyExists, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestContentService_GetComments(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	service := NewContentService(db)
	essayID := uint64(1)

	expectedComments := []models.Comment{
		{UserID: 1, EssayID: 1, CommentText: "First comment", CreatedAt: time.Now()},
		{UserID: 2, EssayID: 1, CommentText: "Second comment", CreatedAt: time.Now()},
	}

	rows := sqlmock.NewRows([]string{"user_id", "essay_id", "comment_text", "created_at"}).
		AddRow(expectedComments[0].UserID, expectedComments[0].EssayID, expectedComments[0].CommentText, expectedComments[0].CreatedAt).
		AddRow(expectedComments[1].UserID, expectedComments[1].EssayID, expectedComments[1].CommentText, expectedComments[1].CreatedAt)

	mock.ExpectQuery(`SELECT user_id, essay_id, comment_text, created_at FROM comment WHERE essay_id = \$1`).
		WithArgs(essayID).
		WillReturnRows(rows)

	comments, err := service.GetComments(essayID)

	assert.NoError(t, err)
	assert.Equal(t, expectedComments, comments)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestContentService_AddComment(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	service := NewContentService(db)
	userID := uint64(1)
	essayID := uint64(1)
	commentText := "New comment"

	mock.ExpectExec(`INSERT INTO comment \(user_id, essay_id, comment_text, created_at\) VALUES \(\$1, \$2, \$3, NOW\(\)\)`).
		WithArgs(userID, essayID, commentText).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := service.AddComment(userID, essayID, commentText)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
