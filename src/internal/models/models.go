package models

import "time"

type User struct {
	ID          uint8  `json:"id"`
	Mail        string `json:"mail"`
	Nickname    string `json:"nickname"`
	Password    string `json:"password"`
	IsModerator bool   `json:"is_moderator"`
	CountChecks int    `json:"count_checks"`
}

type Essay struct {
	ID          uint8     `json:"id"`
	EssayText   string    `json:"essay_text"`
	CompletedAt time.Time `json:"completed_at"`
	Status      string    `json:"status"`
	IsPublished bool      `json:"is_published"`
	UserID      uint8     `json:"user_id"`
	VariantID   uint8     `json:"variant_id"`
}

type Variant struct {
	ID             uint8  `json:"id"`
	VariantTitle   string `json:"variant_title"`
	VariantText    string `json:"variant_text"`
	AuthorPosition string `json:"author_position"`
}

type Comment struct {
	ID          uint8     `json:"id"`
	UserID      uint8     `json:"user_id"`
	EssayID     uint8     `json:"essay_id"`
	CommentText string    `json:"comment_text"`
	CreatedAt   time.Time `json:"created_at"`
}

type Like struct {
	UserID  uint8 `json:"user_id"`
	EssayID uint8 `json:"essay_id"`
}

type Result struct {
	ID         uint8  `json:"id"`
	SumScore   int    `json:"sum_score"`
	AppealText string `json:"appeal_text"`
	EssayID    uint8  `json:"essay_id"`
}

type Criteria struct {
	ID    uint8  `json:"id"`
	Title string `json:"title"`
}

type ResultCriteria struct {
	ResultID   uint8 `json:"result_id"`
	CriteriaID uint8 `json:"criteria_id"`
	Score      int   `json:"score"`
}

type DetailedEssay struct {
	ID             uint8                  `json:"id"`
	VariantID      uint8                  `json:"variant_id"`
	VariantTitle   string                 `json:"variant_title"`
	VariantText    string                 `json:"variant_text"`
	EssayText      string                 `json:"essay_text"`
	CompletedAt    time.Time              `json:"completed_at"`
	Status         string                 `json:"status"`
	IsPublished    bool                   `json:"is_published"`
	AuthorID       uint8                  `json:"author_id"`
	AuthorNickname string                 `json:"author_nickname"`
	Likes          int                    `json:"likes"`
	Comments       []DetailedEssayComment `json:"comments"`
	Results        []DetailedResult       `json:"results"`
}

type AppealEssay struct {
	ID           uint8            `json:"id"`
	VariantID    uint8            `json:"variant_id"`
	VariantTitle string           `json:"variant_title"`
	EssayText    string           `json:"essay_text"`
	CompletedAt  time.Time        `json:"completed_at"`
	Status       string           `json:"status"`
	Results      []DetailedResult `json:"results"`
}

type DetailedEssayComment struct {
	ID             uint8     `json:"id"`
	AuthorNickname string    `json:"author_nickname"`
	CommentText    string    `json:"comment_text"`
	CreatedAt      time.Time `json:"created_at"`
}

type DetailedResult struct {
	K1_score   int
	K2_score   int
	K3_score   int
	K4_score   int
	K5_score   int
	K6_score   int
	K7_score   int
	K8_score   int
	K9_score   int
	K10_score  int
	Score      int
	AppealText string
}
