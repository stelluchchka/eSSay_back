package services

import (
	"database/sql"
	"errors"
	"essay/src/internal/models"
)

var ErrLikeAlreadyExists = errors.New("like already exists")

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

func (s *ContentService) GetVariantByID(variantID uint8) (models.Variant, error) {
	var variant models.Variant

	query := `SELECT id, variant_title, variant_text FROM variant WHERE id = $1`
	err := s.DB.QueryRow(query, variantID).Scan(&variant.ID, &variant.VariantTitle, &variant.VariantText)

	if err != nil {
		return variant, err
	}

	return variant, nil
}

func (s *ContentService) GetLikesCount(essayID uint8) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM "like" WHERE essay_id = $1`
	err := s.DB.QueryRow(query, essayID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (s *ContentService) AddLike(userID uint8, essayID uint8) error {
	query := `
		INSERT INTO "like" (user_id, essay_id)
		VALUES ($1, $2)
		ON CONFLICT (user_id, essay_id) DO NOTHING
		RETURNING true`

	var isInserted bool
	err := s.DB.QueryRow(query, userID, essayID).Scan(&isInserted)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrLikeAlreadyExists
		}
		return err
	}
	return nil
}

func (s *ContentService) GetComments(essayID uint8) ([]models.Comment, error) {
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

func (s *ContentService) AddComment(userID uint8, essayID uint8, text string) error {
	query := `INSERT INTO comment (user_id, essay_id, comment_text, created_at) VALUES ($1, $2, $3, NOW())`
	_, err := s.DB.Exec(query, userID, essayID, text)
	return err
}
