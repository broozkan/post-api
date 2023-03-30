package models

type Post struct {
	Title     string `json:"title"`
	Author    string `json:"author"`
	Link      string `json:"link"`
	Subreddit string `json:"subreddit"`
	Content   string `json:"content"`
	Score     int    `json:"score"`
	Promoted  bool   `json:"promoted"`
	NSFW      bool   `json:"nsfw"`
}
