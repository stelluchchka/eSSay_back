package services

import (
	"database/sql"
	"errors"
	"essay/src/internal/models"
	"time"
)

var ErrWrongID = errors.New("wrong id")
var ErrNoRows = errors.New("query doesn't return a row")

type EssayService struct {
	DB *sql.DB
}

func NewEssayService(db *sql.DB) *EssayService {
	return &EssayService{
		DB: db,
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
			return nil, ErrNoRows
		}
		return nil, err
	}

	return &essay, nil
}

// GetPublishedEssayByID retrieves a published essay by its ID.
func (s *EssayService) GetPublishedEssayByID(id uint8) (*models.PublishedEssay, error) {
	query := `
    SELECT 
        e.id, e.essay_text, 
        e.updated_at AT TIME ZONE 'UTC' AS updated_at,
        e.status, e.is_published, e.variant_id, 
        u.nickname,
        COUNT(l.user_id) AS likes_count,
        c.id AS comment_id,
        c.comment_text AS comment_text,
        c.created_at AS comment_created_at,
        cu.nickname AS comment_nickname
    FROM 
        essay e
    JOIN 
        "user" u ON e.user_id = u.id
    LEFT JOIN 
        "like" l ON l.essay_id = e.id
    LEFT JOIN 
        comment c ON c.essay_id = e.id
    LEFT JOIN 
        "user" cu ON cu.id = c.user_id  -- To get the nickname of the comment author
    WHERE 
        e.id = $1 AND e.is_published = true
    GROUP BY 
        e.id, c.id, u.nickname, cu.nickname
    ORDER BY 
        c.created_at
    `
	rows, err := s.DB.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var response models.PublishedEssay
	var comments []models.PublishedEssayComment
	for rows.Next() {
		var comment models.PublishedEssayComment
		if err := rows.Scan(
			&response.ID,
			&response.EssayText,
			&response.UpdatedAt,
			&response.Status,
			&response.IsPublished,
			&response.VariantID,
			&response.Nickname,
			&response.Likes,
			&comment.ID,
			&comment.CommentText,
			&comment.CreatedAt,
			&comment.Nickname,
		); err != nil {
			// TODO: нормально обработать случай если нет комментов к сочинению
			comments = []models.PublishedEssayComment{}
			break
		}
	}
	response.Comments = comments

	if response.ID == 0 {
		return nil, ErrNoRows
	}
	return &response, nil
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
