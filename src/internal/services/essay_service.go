package services

import (
	"database/sql"
	"errors"
	"essay/src/internal/models"
	"fmt"
	"time"
)

var ErrWrongID = errors.New("wrong id")

type EssayService struct {
	DB *sql.DB
}

func NewEssayService(db *sql.DB) *EssayService {
	return &EssayService{
		DB: db,
	}
}

// GetPublishedEssaysCount retrieves all published essays.
func (s *EssayService) GetEssaysCount() (int, error) {
	var count int

	query := `SELECT COUNT(*) FROM essay WHERE is_published = true`
	err := s.DB.QueryRow(query).Scan(&count)

	if err != nil {
		return 0, err
	}

	return count, nil
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
		if err := rows.Scan(&essay.ID, &essay.EssayText, &essay.CompletedAt, &essay.Status, &essay.IsPublished, &essay.UserID, &essay.VariantID); err != nil {
			return nil, err
		}
		essays = append(essays, essay)
	}

	return essays, nil
}

// GetPublishedEssays retrieves all appeal essays.
func (s *EssayService) GetAppealEssays() ([]models.Essay, error) {
	query := `SELECT id, essay_text, updated_at, status, is_published, user_id, variant_id FROM essay WHERE status = 'appeal'`
	rows, err := s.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	essays := []models.Essay{}
	for rows.Next() {
		var essay models.Essay
		if err := rows.Scan(&essay.ID, &essay.EssayText, &essay.CompletedAt, &essay.Status, &essay.IsPublished, &essay.UserID, &essay.VariantID); err != nil {
			return nil, err
		}
		essays = append(essays, essay)
	}

	return essays, nil
}

// GetEssayByID retrieves an essay by its ID.
func (s *EssayService) GetEssayByID(id uint8) (*models.Essay, error) {
	query := `SELECT id, essay_text, updated_at, status, is_published, user_id, variant_id FROM essay WHERE id = $1`
	row := s.DB.QueryRow(query, id)

	var essay models.Essay
	if err := row.Scan(&essay.ID, &essay.EssayText, &essay.CompletedAt, &essay.Status, &essay.IsPublished, &essay.UserID, &essay.VariantID); err != nil {
		return nil, err
	}

	return &essay, nil
}

// GetDetailedEssayByID retrieves an essay by its ID.
func (s *EssayService) GetDetailedEssayByID(id uint8) (*models.DetailedEssay, error) {
	var essay models.DetailedEssay
	err := s.DB.QueryRow("SELECT e.id, variant_id, essay_text, updated_at, status, is_published, user_id, nickname FROM essay e JOIN \"user\" u ON e.user_id = u.id WHERE e.id = $1", id).Scan(
		&essay.ID,
		&essay.VariantID,
		&essay.EssayText,
		&essay.CompletedAt,
		&essay.Status,
		&essay.IsPublished,
		&essay.AuthorID,
		&essay.AuthorNickname,
	)
	if err != nil {
		return nil, err
	}

	var variant models.Variant
	err = s.DB.QueryRow("SELECT variant_title, variant_text FROM variant WHERE id = $1", essay.VariantID).Scan(
		&variant.VariantTitle,
		&variant.VariantText,
	)
	if err != nil {
		return nil, fmt.Errorf("variant fetching error: %w", err)
	}
	essay.VariantTitle = variant.VariantTitle
	essay.VariantText = variant.VariantText

	err = s.DB.QueryRow("SELECT COUNT(*) FROM \"like\" WHERE essay_id = $1", essay.ID).Scan(&essay.Likes)
	if err != nil {
		return nil, fmt.Errorf("like fetching error: %w", err)
	}

	rows, err := s.DB.Query("SELECT c.id, u.nickname, comment_text, created_at FROM comment c JOIN \"user\" u ON c.user_id = u.id WHERE essay_id = $1 ORDER BY created_at DESC", essay.ID)
	if err != nil {
		return nil, fmt.Errorf("comment fetching error: %w", err)
	}
	defer rows.Close()
	var comments []models.DetailedEssayComment
	for rows.Next() {
		var comment models.DetailedEssayComment
		err := rows.Scan(&comment.ID, &comment.AuthorNickname, &comment.CommentText, &comment.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("comment scanning error: %w", err)
		}
		comments = append(comments, comment)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("comment iterating error: %w", err)
	}
	essay.Comments = comments

	query := `
	SELECT 
		COALESCE(SUM(CASE WHEN rc.criteria_id = 1 THEN rc.score ELSE 0 END), 0) AS K1_score,
		COALESCE(SUM(CASE WHEN rc.criteria_id = 2 THEN rc.score ELSE 0 END), 0) AS K2_score,
		COALESCE(SUM(CASE WHEN rc.criteria_id = 3 THEN rc.score ELSE 0 END), 0) AS K3_score,
		COALESCE(SUM(CASE WHEN rc.criteria_id = 4 THEN rc.score ELSE 0 END), 0) AS K4_score,
		COALESCE(SUM(CASE WHEN rc.criteria_id = 5 THEN rc.score ELSE 0 END), 0) AS K5_score,
		COALESCE(SUM(CASE WHEN rc.criteria_id = 6 THEN rc.score ELSE 0 END), 0) AS K6_score,
		COALESCE(SUM(CASE WHEN rc.criteria_id = 7 THEN rc.score ELSE 0 END), 0) AS K7_score,
		COALESCE(SUM(CASE WHEN rc.criteria_id = 8 THEN rc.score ELSE 0 END), 0) AS K8_score,
		COALESCE(SUM(CASE WHEN rc.criteria_id = 9 THEN rc.score ELSE 0 END), 0) AS K9_score,
		COALESCE(SUM(CASE WHEN rc.criteria_id = 10 THEN rc.score ELSE 0 END), 0) AS K10_score,
		SUM(rc.score) AS Score
	FROM 
		result r
	LEFT JOIN 
		result_criteria rc ON r.id = rc.result_id
	WHERE 
		r.essay_id = $1
	GROUP BY 
		r.id
	`

	rows, err = s.DB.Query(query, essay.ID)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	var detailedResults []models.DetailedResult
	for rows.Next() {
		var result models.DetailedResult
		err := rows.Scan(
			&result.K1_score,
			&result.K2_score,
			&result.K3_score,
			&result.K4_score,
			&result.K5_score,
			&result.K6_score,
			&result.K7_score,
			&result.K8_score,
			&result.K9_score,
			&result.K10_score,
			&result.Score,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning results: %w", err)
		}
		detailedResults = append(detailedResults, result)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating results: %w", err)
	}

	essay.Results = detailedResults

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
		if err := rows.Scan(&essay.ID, &essay.EssayText, &essay.CompletedAt, &essay.Status, &essay.IsPublished, &essay.UserID, &essay.VariantID); err != nil {
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
