package services

import (
	"database/sql"
	"errors"
	"essay/src/internal/database"
	"essay/src/internal/models"
	"time"
)

var ErrWrongID = errors.New("wrong id or user_id")

type EssayService struct {
	DB *sql.DB
}

func NewEssayService() *EssayService {
	return &EssayService{
		DB: database.GetPostgreSQLConnection(),
	}
}

// GetPublishedEssays retrieves all published essays.
func (s *EssayService) GetPublishedEssays() ([]models.Essay, error) {
	query := `SELECT id, essay_text, updated_at, status, is_published, user_id, variant_id FROM essay WHERE is_published = true`
	rows, err := s.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	essays := []models.Essay{}
	for rows.Next() {
		var essay models.Essay
		if err := rows.Scan(&essay.ID, &essay.EssayText, &essay.UpdatedAt, &essay.Status, &essay.IsPublished, &essay.UserID, &essay.VariantID); err != nil {
			return nil, err
		}
		essays = append(essays, essay)
	}

	return essays, nil
}

// GetEssayByID retrieves a published essay by its ID.
func (s *EssayService) GetEssayByID(id uint8) (*models.Essay, error) {
	query := `SELECT id, essay_text, updated_at, status, is_published, user_id, variant_id FROM essay WHERE id = $1`
	row := s.DB.QueryRow(query, id)

	var essay models.Essay
	if err := row.Scan(&essay.ID, &essay.EssayText, &essay.UpdatedAt, &essay.Status, &essay.IsPublished, &essay.UserID, &essay.VariantID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &essay, nil
}

// GetPublishedEssayByID retrieves a published essay by its ID.
func (s *EssayService) GetPublishedEssayByID(id uint8) (*models.Essay, error) {
	query := `SELECT id, essay_text, updated_at, status, is_published, user_id, variant_id FROM essay WHERE id = $1 AND is_published = true`
	row := s.DB.QueryRow(query, id)

	var essay models.Essay
	if err := row.Scan(&essay.ID, &essay.EssayText, &essay.UpdatedAt, &essay.Status, &essay.IsPublished, &essay.UserID, &essay.VariantID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &essay, nil
}

// GetUserEssays retrieves all essays for a specific user.
func (s *EssayService) GetUserEssays(userID uint8) ([]models.Essay, error) {
	query := `SELECT id, essay_text, updated_at, status, is_published, user_id, variant_id FROM essay WHERE user_id = $1`
	rows, err := s.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	essays := []models.Essay{}
	for rows.Next() {
		var essay models.Essay
		if err := rows.Scan(&essay.ID, &essay.EssayText, &essay.UpdatedAt, &essay.Status, &essay.IsPublished, &essay.UserID, &essay.VariantID); err != nil {
			return nil, err
		}
		essays = append(essays, essay)
	}

	return essays, nil
}

// CreateEssay creates a new essay in draft status.
func (s *EssayService) CreateEssay(essay *models.Essay) error {
	query := `INSERT INTO essay (essay_text, updated_at, status, is_published, user_id, variant_id) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := s.DB.Exec(query, essay.EssayText, time.Now(), "draft", false, essay.UserID, essay.VariantID)
	if err != nil {
		return err
	}
	return nil
}

// UpdateEssay updates an existing essay.
func (s *EssayService) UpdateEssay(essay *models.Essay) error {
	countQuery := `SELECT COUNT(*) FROM essay WHERE id = $1 AND user_id = $2`
	var count int

	err := s.DB.QueryRow(countQuery, essay.ID, essay.UserID).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		query := `UPDATE essay SET essay_text = $1 WHERE id = $2 AND user_id = $3`
		_, err = s.DB.Exec(query, essay.EssayText, essay.ID, essay.UserID)
		if err != nil {
			return err
		}
		return nil
	} else {
		return ErrWrongID
	}
}

// ChangeEssayStatus updates the status of an essay.
func (s *EssayService) ChangeEssayStatus(essayID uint8, userID uint8, status string) error {
	query := `UPDATE essay SET status = $1 WHERE id = $2 AND user_id = $3`
	_, err := s.DB.Exec(query, status, essayID, userID)
	if err != nil {
		return err
	}
	return nil
}

// PublishEssay marks an essay as published.
func (s *EssayService) PublishEssay(essayID uint8, userID uint8) error {
	query := `UPDATE essay SET is_published = true WHERE id = $1 AND user_id = $2`
	_, err := s.DB.Exec(query, essayID, userID)
	if err != nil {
		return err
	}
	return nil
}
