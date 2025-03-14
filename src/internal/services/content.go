package services

import (
	"database/sql"
	"essay/src/internal/models"
)

func (s *UserService) GetCounts() (int, int, int, error) {
	var variants_count int
	var essays_count int
	var users_count int

	variants_query := `SELECT COUNT(*) FROM variant`
	err := s.DB.QueryRow(variants_query).Scan(&variants_count)
	if err != nil {
		return 0, 0, 0, err
	}

	essays_query := `SELECT COUNT(*) FROM essay`
	err = s.DB.QueryRow(essays_query).Scan(&essays_count)
	if err != nil {
		return 0, 0, 0, err
	}

	users_query := `SELECT COUNT(*) FROM "user"`
	err = s.DB.QueryRow(users_query).Scan(&users_count)
	if err != nil {
		return 0, 0, 0, err
	}

	return variants_count, essays_count, users_count, nil
}

func (s *UserService) GetVariantByID(variantID uint64) (models.Variant, error) {
	var variant models.Variant

	query := `SELECT id, variant_title, variant_text FROM variant WHERE id = $1`
	err := s.DB.QueryRow(query, variantID).Scan(&variant.ID, &variant.VariantTitle, &variant.VariantText)

	if err != nil {
		return variant, err
	}

	return variant, nil
}

func (s *UserService) GetLikesCount(essayID uint64) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM "like" WHERE essay_id = $1`
	err := s.DB.QueryRow(query, essayID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (s *UserService) IsLiked(userID uint64, essayID uint64) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM "like" WHERE user_id = $1 AND essay_id = $2`
	err := s.DB.QueryRow(query, userID, essayID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (s *UserService) AddLike(userID uint64, essayID uint64) error {
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

func (s *UserService) DeleteLike(userID uint64, essayID uint64) error {
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

func (s *UserService) GetComments(essayID uint64) ([]models.Comment, error) {
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

func (s *UserService) AddComment(userID uint64, essayID uint64, text string) (models.DetailedEssayComment, error) {
	query := `INSERT INTO comment (user_id, essay_id, comment_text, created_at) 
              VALUES ($1, $2, $3, NOW()) 
              RETURNING id, user_id, comment_text, created_at`

	var comment models.DetailedEssayComment
	err := s.DB.QueryRow(query, userID, essayID, text).Scan(&comment.ID, &comment.AuthorNickname, &comment.CommentText, &comment.CreatedAt)

	if err != nil {
		return models.DetailedEssayComment{}, err
	}

	var authorNickname string
	err = s.DB.QueryRow(`SELECT nickname FROM "user" WHERE id = $1`, userID).Scan(&authorNickname)
	if err != nil {
		return models.DetailedEssayComment{}, err
	}

	comment.AuthorNickname = authorNickname

	return comment, nil
}
