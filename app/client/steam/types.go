package steam

import "time"

type Comment struct {
	Author       string    `json:"author"`
	AuthorURL    string    `json:"author_url"`
	AvatarURL    string    `json:"avatar_url"`
	Content      string    `json:"content"`
	Timestamp    time.Time `json:"timestamp"`
	TimestampRaw string    `json:"timestamp_raw"`
	CommentID    string    `json:"comment_id"`
}

type ParseResult struct {
	Comments []Comment `json:"comments"`
}
