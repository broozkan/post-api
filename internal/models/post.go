package models

type Post struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Author    string `json:"author"`
	Link      string `json:"link"`
	Subreddit string `json:"subreddit"`
	Content   string `json:"content"`
	Score     int    `json:"score"`
	Promoted  bool   `json:"promoted"`
	NSFW      bool   `json:"nsfw"`
}
