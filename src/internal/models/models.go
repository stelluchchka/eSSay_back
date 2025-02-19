package models

import "time"

type User struct {
	ID          uint64 `json:"id"`
	Mail        string `json:"mail"`
	Nickname    string `json:"nickname"`
	Password    string `json:"password"`
	IsModerator bool   `json:"is_moderator"`
	CountChecks int    `json:"count_checks"`
}

type Essay struct {
	ID          uint64    `json:"id"`
	EssayText   string    `json:"essay_text"`
	CompletedAt time.Time `json:"completed_at"`
	Status      string    `json:"status"`
	IsPublished bool      `json:"is_published"`
	UserID      uint64    `json:"user_id"`
	VariantID   uint64    `json:"variant_id"`
}

type Variant struct {
	ID             uint64 `json:"id"`
	VariantTitle   string `json:"variant_title"`
	VariantText    string `json:"variant_text"`
	AuthorPosition string `json:"author_position"`
}

type Comment struct {
	ID          uint64    `json:"id"`
	UserID      uint64    `json:"user_id"`
	EssayID     uint64    `json:"essay_id"`
	CommentText string    `json:"comment_text"`
	CreatedAt   time.Time `json:"created_at"`
}

type Like struct {
	UserID  uint64 `json:"user_id"`
	EssayID uint64 `json:"essay_id"`
}

type Result struct {
	ID         uint64 `json:"id"`
	SumScore   int    `json:"sum_score"`
	AppealText string `json:"appeal_text"`
	EssayID    uint64 `json:"essay_id"`
}

type Criteria struct {
	ID    uint64 `json:"id"`
	Title string `json:"title"`
}

type ResultCriteria struct {
	ResultID   uint64 `json:"result_id"`
	CriteriaID uint64 `json:"criteria_id"`
	Score      int    `json:"score"`
}

type DetailedEssay struct {
	ID             uint64                 `json:"id"`
	VariantID      uint64                 `json:"variant_id"`
	VariantTitle   string                 `json:"variant_title"`
	VariantText    string                 `json:"variant_text"`
	EssayText      string                 `json:"essay_text"`
	CompletedAt    time.Time              `json:"completed_at"`
	Status         string                 `json:"status"`
	IsPublished    bool                   `json:"is_published"`
	AuthorID       uint64                 `json:"author_id"`
	AuthorNickname string                 `json:"author_nickname"`
	Likes          int                    `json:"likes"`
	Comments       []DetailedEssayComment `json:"comments"`
	Results        []DetailedResult       `json:"results"`
}

type AppealEssay struct {
	ID           uint64           `json:"id"`
	VariantID    uint64           `json:"variant_id"`
	VariantTitle string           `json:"variant_title"`
	EssayText    string           `json:"essay_text"`
	CompletedAt  time.Time        `json:"completed_at"`
	Status       string           `json:"status"`
	Results      []DetailedResult `json:"results"`
}

type DetailedEssayComment struct {
	ID             uint64    `json:"id"`
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

type UserInfo struct {
	ID                   uint64  `json:"id"`
	Mail                 string  `json:"mail"`
	Nickname             string  `json:"nickname"`
	IsModerator          bool    `json:"is_moderator"`
	CountChecks          int     `json:"count_checks"`
	CountEssays          int     `json:"count_essays"`
	CountPublishedEssays int     `json:"count_published_essays"`
	AverageResult        float64 `json:"average_result"`
}
