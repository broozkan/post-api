package models

import "time"

type (
	PostRow struct {
		Post Post `json:"postData"`
	}

	CountResult struct {
		Count int `json:"count"`
	}

	Post struct {
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

	AdPlacement struct {
		PostIndex   int
		Position    string
		Ad          *Ad
		PreviousNSF bool
	}

	Ad struct {
		ID           string    `json:"id"`
		Post         *Post     `json:"post"`
		Title        string    `json:"title"`
		Link         string    `json:"link"`
		ImageURL     string    `json:"image_url"`
		TargetGeo    string    `json:"target_geo"`
		TargetAgeMin int       `json:"target_age_min"`
		TargetAgeMax int       `json:"target_age_max"`
		CreatedAt    time.Time `json:"created_at"`
		ExpiresAt    time.Time `json:"expires_at"`
	}

	ListPostsResponse struct {
		Posts      []*Post `json:"posts"`
		Page       int     `json:"page"`
		TotalPages int     `json:"totalPages"`
	}
)
