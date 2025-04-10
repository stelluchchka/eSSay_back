package services

import (
	"database/sql"
	"essay/src/internal/models"
	"log"
)

func (s *UserService) GetCounts() (int, int, int, error) {
	var variants_count int
	var essays_count int
	var users_count int

	variants_query := `SELECT COUNT(*) FROM variant WHERE is_public = TRUE`
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

func (s *UserService) CreateVariant(variant models.Variant) (int, error) {
	var insertedID int
	query := `
		INSERT INTO variant (variant_title, variant_text, author_position)  
		VALUES ($1, $2, $3)  
		RETURNING id;`
	err := s.DB.QueryRow(query, variant.VariantTitle, variant.VariantText, variant.AuthorPosition).Scan(&insertedID)

	if err != nil {
		return 0, err
	}

	return insertedID, nil
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

func (s *UserService) CreateResult(result *models.DetailedResult, essayID uint64) error {
	var resultID int
	err := s.DB.QueryRow(`
		INSERT INTO result (sum_score, essay_id) 
		VALUES ($1, $2) RETURNING id`,
		result.Score, essayID,
	).Scan(&resultID)
	if err != nil {
		return err
	}
	log.Println("resultID ", resultID)

	// Вставляем критерии оценки
	criteriaScores := []struct {
		Score       int
		Explanation string
		CriteriaID  int
	}{
		{result.K1_score, result.K1_explanation, 1},
		{result.K2_score, result.K2_explanation, 2},
		{result.K3_score, result.K3_explanation, 3},
		{result.K4_score, result.K4_explanation, 4},
		{result.K5_score, result.K5_explanation, 5},
		{result.K6_score, result.K6_explanation, 6},
		{result.K7_score, result.K7_explanation, 7},
		{result.K8_score, result.K8_explanation, 8},
		{result.K9_score, result.K9_explanation, 9},
		{result.K10_score, result.K10_explanation, 10},
	}

	for _, c := range criteriaScores {
		_, err = s.DB.Exec(`
			INSERT INTO result_criteria (result_id, criteria_id, score, explanation)
			VALUES ($1, $2, $3, $4)`,
			resultID, c.CriteriaID, c.Score, c.Explanation,
		)
		if err != nil {
			return err
		}
	}

	log.Println("Resilt saved with ID:", resultID)

	return nil
}
