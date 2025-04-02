package services

import (
	"database/sql"
	"errors"
	"essay/src/internal/models"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestUserService_GetPublishedEssays(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	service := NewUserService(db)

	expectedEssays := []models.Essay{
		{ID: 1, EssayText: "Essay 1", CompletedAt: time.Now(), Status: "published", IsPublished: true, UserID: 1, VariantID: 1},
		{ID: 2, EssayText: "Essay 2", CompletedAt: time.Now(), Status: "published", IsPublished: true, UserID: 2, VariantID: 2},
	}

	rows := sqlmock.NewRows([]string{"id", "essay_text", "updated_at", "status", "is_published", "user_id", "variant_id"}).
		AddRow(expectedEssays[0].ID, expectedEssays[0].EssayText, expectedEssays[0].CompletedAt, expectedEssays[0].Status, expectedEssays[0].IsPublished, expectedEssays[0].UserID, expectedEssays[0].VariantID).
		AddRow(expectedEssays[1].ID, expectedEssays[1].EssayText, expectedEssays[1].CompletedAt, expectedEssays[1].Status, expectedEssays[1].IsPublished, expectedEssays[1].UserID, expectedEssays[1].VariantID)

	mock.ExpectQuery(`SELECT id, essay_text, updated_at, status, is_published, user_id, variant_id FROM essay WHERE is_published = true`).
		WillReturnRows(rows)

	essays, err := service.GetPublishedEssays()

	assert.NoError(t, err)
	assert.Equal(t, expectedEssays, essays)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserService_GetPublishedEssays_NoRows(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	service := NewUserService(db)

	mock.ExpectQuery(`SELECT id, essay_text, updated_at, status, is_published, user_id, variant_id FROM essay WHERE is_published = true`).
		WillReturnRows(sqlmock.NewRows(nil))

	essays, err := service.GetPublishedEssays()

	assert.NoError(t, err)
	assert.Empty(t, essays)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserService_GetAppealEssays(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	service := NewUserService(db)

	expectedEssays := []models.Essay{
		{ID: 1, EssayText: "Appeal Essay 1", CompletedAt: time.Now(), Status: "appeal", IsPublished: false, UserID: 1, VariantID: 1},
	}

	rows := sqlmock.NewRows([]string{"id", "essay_text", "updated_at", "status", "is_published", "user_id", "variant_id"}).
		AddRow(expectedEssays[0].ID, expectedEssays[0].EssayText, expectedEssays[0].CompletedAt, expectedEssays[0].Status, expectedEssays[0].IsPublished, expectedEssays[0].UserID, expectedEssays[0].VariantID)

	mock.ExpectQuery(`SELECT id, essay_text, updated_at, status, is_published, user_id, variant_id FROM essay WHERE status = 'appeal'`).
		WillReturnRows(rows)

	essays, err := service.GetAppealEssays()

	assert.NoError(t, err)
	assert.Equal(t, expectedEssays, essays)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserService_GetUserEssays(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	service := NewUserService(db)
	userID := uint64(1)

	expectedEssays := []models.Essay{
		{ID: 1, EssayText: "User Essay 1", CompletedAt: time.Now(), Status: "draft", IsPublished: false, UserID: userID, VariantID: 1},
		{ID: 2, EssayText: "User Essay 2", CompletedAt: time.Now(), Status: "published", IsPublished: true, UserID: userID, VariantID: 2},
	}

	rows := sqlmock.NewRows([]string{"id", "essay_text", "updated_at", "status", "is_published", "user_id", "variant_id"}).
		AddRow(expectedEssays[0].ID, expectedEssays[0].EssayText, expectedEssays[0].CompletedAt, expectedEssays[0].Status, expectedEssays[0].IsPublished, expectedEssays[0].UserID, expectedEssays[0].VariantID).
		AddRow(expectedEssays[1].ID, expectedEssays[1].EssayText, expectedEssays[1].CompletedAt, expectedEssays[1].Status, expectedEssays[1].IsPublished, expectedEssays[1].UserID, expectedEssays[1].VariantID)

	mock.ExpectQuery(`SELECT id, essay_text, updated_at, status, is_published, user_id, variant_id FROM essay WHERE user_id = \$1`).
		WithArgs(userID).
		WillReturnRows(rows)

	essays, err := service.GetUserEssays(userID)

	assert.NoError(t, err)
	assert.Equal(t, expectedEssays, essays)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserService_GetEssayByID_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	service := NewUserService(db)
	id := uint64(1)

	expectedEssay := models.Essay{
		ID:          id,
		EssayText:   "Sample Essay",
		CompletedAt: time.Now(),
		Status:      "draft",
		IsPublished: false,
		UserID:      1,
		VariantID:   1,
	}

	mock.ExpectQuery(`SELECT id, essay_text, updated_at, status, is_published, user_id, variant_id FROM essay WHERE id = \$1`).
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{"id", "essay_text", "updated_at", "status", "is_published", "user_id", "variant_id"}).
			AddRow(expectedEssay.ID, expectedEssay.EssayText, expectedEssay.CompletedAt, expectedEssay.Status, expectedEssay.IsPublished, expectedEssay.UserID, expectedEssay.VariantID))

	essay, err := service.GetEssayByID(id)

	assert.NoError(t, err)
	assert.Equal(t, &expectedEssay, essay)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserService_GetEssayByID_NoRows(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	service := NewUserService(db)
	id := uint64(1)

	mock.ExpectQuery(`SELECT id, essay_text, updated_at, status, is_published, user_id, variant_id FROM essay WHERE id = \$1`).
		WithArgs(id).
		WillReturnError(sql.ErrNoRows)

	essay, err := service.GetEssayByID(id)

	assert.Error(t, err)
	assert.Equal(t, sql.ErrNoRows, err)
	assert.Nil(t, essay)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// func TestUserService_GetDetailedEssayByID_NoRows(t *testing.T) {
// 	db, mock, _ := sqlmock.New()
// 	defer db.Close()

// 	service := NewUserService(db)
// 	essayID := uint64(1)

// 	mock.ExpectQuery(`SELECT * FROM essay e JOIN "user" u ON e.user_id = u.id WHERE e.id = \$1`).
// 		WithArgs(essayID).
// 		WillReturnError(sql.ErrNoRows)

// 	essay, err := service.GetDetailedEssayByID(essayID)

// 	assert.Nil(t, essay)
// 	assert.Equal(t, sql.ErrNoRows, err)
// 	assert.NoError(t, mock.ExpectationsWereMet())
// }

// func TestGetDetailedEssayByID(t *testing.T) {
// 	db, mock, _ := sqlmock.New()
// 	defer db.Close()

// 	service := NewUserService(db)
// 	ID := uint64(1)

// 	expectedComments := []models.DetailedEssayComment{
// 		{
// 			ID:             ID,
// 			AuthorNickname: "Commentator 1",
// 			CommentText:    "Great essay!",
// 			CreatedAt:      time.Now(),
// 		},
// 	}

// 	expectedResults := []models.DetailedResult{
// 		{
// 			K1_score:  0,
// 			K2_score:  0,
// 			K3_score:  0,
// 			K4_score:  0,
// 			K5_score:  0,
// 			K6_score:  0,
// 			K7_score:  0,
// 			K8_score:  0,
// 			K9_score:  0,
// 			K10_score: 0,
// 			Score:     0,
// 		},
// 	}

// 	expectedEssay := models.DetailedEssay{
// 		ID:             ID,
// 		VariantID:      1,
// 		VariantTitle:   "variant_title",
// 		VariantText:    "variant_text",
// 		EssayText:      "Sample Essay",
// 		CompletedAt:      time.Now(),
// 		Status:         "draft",
// 		IsPublished:    true,
// 		AuthorID:       ID,
// 		AuthorNickname: "John Doe",
// 		Likes:          5,
// 		Comments:       expectedComments,
// 		Results:        expectedResults,
// 	}

// 	rows := sqlmock.NewRows([]string{"id", "variant_id", "essay_text", "updated_at", "status", "is_published", "user_id", "nickname"}).
// 		AddRow(ID, 1, "Sample Essay", time.Now(), "draft", true, 101, "John Doe")

// 	mock.ExpectQuery(`SELECT e.id, variant_id, essay_text, updated_at, status, is_published, user_id, nickname FROM essay e JOIN \"user\" u ON e.user_id = u.id WHERE e.id = $1`).
// 		WithArgs(ID).
// 		WillReturnRows(rows)

// 	variantRows := sqlmock.NewRows([]string{"variant_title", "variant_text"}).AddRow("variant_title", "variant_text")
// 	mock.ExpectQuery(`SELECT variant_title, variant_text FROM variant WHERE id = $1`).
// 		WithArgs(ID).
// 		WillReturnRows(variantRows)

// 	likeCount := 5
// 	mock.ExpectQuery(`SELECT COUNT(*) FROM \"like\" WHERE essay_id = $1`).
// 		WithArgs(ID).
// 		WillReturnRows(sqlmock.NewRows([]string{"COUNT"}).AddRow(likeCount))

// 	commentRows := sqlmock.NewRows([]string{"user_id", "nickname", "comment_text", "created_at"}).
// 		AddRow(1, "Commentator 1", "Great essay!", time.Now().Add(-24*time.Hour))
// 	mock.ExpectQuery(`SELECT c.id, u.nickname, comment_text, created_at FROM comment c JOIN \"user\" u ON c.user_id = u.id WHERE essay_id = $1 ORDER BY created_at DESC`).
// 		WithArgs(ID).
// 		WillReturnRows(commentRows)

// 	resultRows := sqlmock.NewRows([]string{"K1_score", "K2_score", "K3_score", "K4_score", "K5_score", "K6_score", "K7_score", "K8_score", "K9_score", "K10_score", "Score"}).
// 		AddRow(0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0)
// 	mock.ExpectQuery(`SELECT
// 		COALESCE(SUM(CASE WHEN rc.criteria_id = 1 THEN rc.score ELSE 0 END), 0) AS K1_score,
// 		COALESCE(SUM(CASE WHEN rc.criteria_id = 2 THEN rc.score ELSE 0 END), 0) AS K2_score,
// 		COALESCE(SUM(CASE WHEN rc.criteria_id = 3 THEN rc.score ELSE 0 END), 0) AS K3_score,
// 		COALESCE(SUM(CASE WHEN rc.criteria_id = 4 THEN rc.score ELSE 0 END), 0) AS K4_score,
// 		COALESCE(SUM(CASE WHEN rc.criteria_id = 5 THEN rc.score ELSE 0 END), 0) AS K5_score,
// 		COALESCE(SUM(CASE WHEN rc.criteria_id = 6 THEN rc.score ELSE 0 END), 0) AS K6_score,
// 		COALESCE(SUM(CASE WHEN rc.criteria_id = 7 THEN rc.score ELSE 0 END), 0) AS K7_score,
// 		COALESCE(SUM(CASE WHEN rc.criteria_id = 8 THEN rc.score ELSE 0 END), 0) AS K8_score,
// 		COALESCE(SUM(CASE WHEN rc.criteria_id = 9 THEN rc.score ELSE 0 END), 0) AS K9_score,
// 		COALESCE(SUM(CASE WHEN rc.criteria_id = 10 THEN rc.score ELSE 0 END), 0) AS K10_score,
// 		SUM(rc.score) AS Score
// 	FROM
// 		result r
// 	LEFT JOIN
// 		result_criteria rc ON r.id = rc.result_id
// 	WHERE
// 		r.essay_id = $1
// 	GROUP BY
// 		r.id`).
// 		WithArgs(ID).
// 		WillReturnRows(resultRows)

// 	essay, err := service.GetDetailedEssayByID(ID)

// 	assert.Error(t, err)
// 	assert.Nil(t, essay)
// 	assert.Equal(t, &expectedEssay, essay)
// 	assert.NoError(t, mock.ExpectationsWereMet())
// }

func TestUserService_CreateEssay(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	service := NewUserService(db)
	newEssay := models.Essay{
		EssayText:   "New Essay",
		UserID:      1,
		VariantID:   1,
		IsPublished: false,
	}

	mock.ExpectQuery(`INSERT INTO essay \(essay_text, updated_at, status, is_published, user_id, variant_id\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6\) RETURNING id`).
		WithArgs(newEssay.EssayText, sqlmock.AnyArg(), "draft", false, newEssay.UserID, newEssay.VariantID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	id, err := service.CreateEssay(&newEssay)

	assert.NoError(t, err)
	assert.Equal(t, 1, id)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserService_ChangeEssayStatus(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	service := NewUserService(db)
	essayID := uint64(1)
	status := "published"

	mock.ExpectExec(`UPDATE essay SET status = \$1 WHERE id = \$2 AND user_id = \$3`).
		WithArgs(status, essayID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := service.ChangeEssayStatus(essayID, status)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserService_ChangeEssayStatus_Error(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	service := NewUserService(db)
	essayID := uint64(1)
	userID := uint64(1)
	status := "published"

	mock.ExpectExec(`UPDATE essay SET status = \$1 WHERE id = \$2 AND user_id = \$3`).
		WithArgs(status, essayID, userID).
		WillReturnError(errors.New("update failed"))

	err := service.ChangeEssayStatus(essayID, status)

	assert.Error(t, err)
	assert.EqualError(t, err, "update failed")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserService_UpdateEssay_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	service := NewUserService(db)
	updatedEssay := models.Essay{
		ID:        1,
		EssayText: "Updated Text",
		UserID:    1,
	}

	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM essay WHERE id = \$1 AND user_id = \$2`).
		WithArgs(updatedEssay.ID, updatedEssay.UserID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	mock.ExpectExec(`UPDATE essay SET essay_text = \$1 WHERE id = \$2 AND user_id = \$3`).
		WithArgs(updatedEssay.EssayText, updatedEssay.ID, updatedEssay.UserID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := service.UpdateEssay(&updatedEssay)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserService_UpdateEssay_WrongID(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	service := NewUserService(db)
	updatedEssay := models.Essay{
		ID:        1,
		EssayText: "Updated Text",
		UserID:    1,
	}

	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM essay WHERE id = \$1 AND user_id = \$2`).
		WithArgs(updatedEssay.ID, updatedEssay.UserID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	err := service.UpdateEssay(&updatedEssay)

	assert.Error(t, err)
	assert.Equal(t, ErrWrongID, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
func TestUserService_PublishEssay(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	service := NewUserService(db)
	essayID := uint64(1)
	userID := uint64(1)

	mock.ExpectExec(`UPDATE essay SET is_published = true WHERE id = \$1 AND user_id = \$2`).
		WithArgs(essayID, userID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := service.PublishEssay(essayID, userID)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserService_PublishEssay_Error(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	service := NewUserService(db)
	essayID := uint64(1)
	userID := uint64(1)

	mock.ExpectExec(`UPDATE essay SET is_published = true WHERE id = \$1 AND user_id = \$2`).
		WithArgs(essayID, userID).
		WillReturnError(errors.New("update failed"))

	err := service.PublishEssay(essayID, userID)

	assert.Error(t, err)
	assert.EqualError(t, err, "update failed")
	assert.NoError(t, mock.ExpectationsWereMet())
}
