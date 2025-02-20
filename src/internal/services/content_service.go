package services

import (
	"database/sql"
	"errors"
	"essay/src/internal/models"
)

var ErrLikeAlreadyExists = errors.New("like already exists")
var ErrLikeNotFound = errors.New("like doesn't exists")

type ContentService struct {
	DB *sql.DB
}

func NewContentService(db *sql.DB) *ContentService {
	return &ContentService{
		DB: db,
	}
}

func (s *ContentService) GetVariantsCount() (int, error) {
	var count int

	query := `SELECT COUNT(*) FROM variant`
	err := s.DB.QueryRow(query).Scan(&count)

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (s *ContentService) GetVariantByID(variantID uint64) (models.Variant, error) {
	var variant models.Variant

	query := `SELECT id, variant_title, variant_text FROM variant WHERE id = $1`
	err := s.DB.QueryRow(query, variantID).Scan(&variant.ID, &variant.VariantTitle, &variant.VariantText)

	if err != nil {
		return variant, err
	}

	return variant, nil
}

func (s *ContentService) GetLikesCount(essayID uint64) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM "like" WHERE essay_id = $1`
	err := s.DB.QueryRow(query, essayID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (s *ContentService) IsLiked(userID uint64, essayID uint64) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM "like" WHERE user_id = $1 AND essay_id = $2`
	err := s.DB.QueryRow(query, userID, essayID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (s *ContentService) AddLike(userID uint64, essayID uint64) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
        INSERT INTO "like" (user_id, essay_id)
        VALUES ($1, $2)
        ON CONFLICT (user_id, essay_id) DO NOTHING
        RETURNING true`

	var isInserted bool
	err = tx.QueryRow(query, userID, essayID).Scan(&isInserted)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrLikeAlreadyExists
		}
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *ContentService) DeleteLike(userID uint64, essayID uint64) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `DELETE FROM "like" WHERE user_id = $1 AND essay_id = $2`
	result, err := tx.Exec(query, userID, essayID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrLikeNotFound
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *ContentService) GetComments(essayID uint64) ([]models.Comment, error) {
	query := `SELECT user_id, essay_id, comment_text, created_at FROM comment WHERE essay_id = $1`
	rows, err := s.DB.Query(query, essayID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var comment models.Comment
		err := rows.Scan(&comment.UserID, &comment.EssayID, &comment.CommentText, &comment.CreatedAt)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

func (s *ContentService) AddComment(userID uint64, essayID uint64, text string) error {
	query := `INSERT INTO comment (user_id, essay_id, comment_text, created_at) VALUES ($1, $2, $3, NOW())`
	_, err := s.DB.Exec(query, userID, essayID, text)
	return err
}
