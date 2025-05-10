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
	ID       uint64 `json:"id"`
	Title    string `json:"title"`
	MaxScore int    `json:"max_score"`
}

type ResultCriteria struct {
	ResultID    uint64 `json:"result_id"`
	CriteriaID  uint64 `json:"criteria_id"`
	Score       int    `json:"score"`
	Explanation string `json:"explanation"`
}

type EssayCard struct {
	ID             uint64 `json:"id"`
	VariantID      uint64 `json:"variant_id"`
	VariantTitle   string `json:"variant_title"`
	AuthorNickname string `json:"author_nickname"`
	Likes          int    `json:"likes"`
	Score          int    `json:"score"`
	Status         string `json:"status"`
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
	K1_score        int     `json:"K1_score"`
	K2_score        int     `json:"K2_score"`
	K3_score        int     `json:"K3_score"`
	K4_score        int     `json:"K4_score"`
	K5_score        int     `json:"K5_score"`
	K6_score        int     `json:"K6_score"`
	K7_score        int     `json:"K7_score"`
	K8_score        int     `json:"K8_score"`
	K9_score        int     `json:"K9_score"`
	K10_score       int     `json:"K10_score"`
	K1_explanation  string  `json:"K1_explanation"`
	K2_explanation  string  `json:"K2_explanation"`
	K3_explanation  string  `json:"K3_explanation"`
	K4_explanation  string  `json:"K4_explanation"`
	K5_explanation  string  `json:"K5_explanation"`
	K6_explanation  string  `json:"K6_explanation"`
	K7_explanation  string  `json:"K7_explanation"`
	K8_explanation  string  `json:"K8_explanation"`
	K9_explanation  string  `json:"K9_explanation"`
	K10_explanation string  `json:"K10_explanation"`
	Score           *int    `json:"score,omitempty"`
	AppealText      *string `json:"appealText,omitempty"`
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

type EssayRequest struct {
	EssayID        uint64 `json:"essay_id"`
	EssayText      string `json:"essay_text"`
	VariantText    string `json:"variant_text"`
	AuthorPosition string `json:"author_position"`
}

type ResultDate struct {
	CompletedAt time.Time `json:"completed_at"`
	Score       int       `json:"score"`
}
