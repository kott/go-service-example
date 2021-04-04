package articles

import "time"

// Article is the nominal object used for interacting with articles.
// This represents what is stored in the database.
type Article struct {
	ID string `json:"id"`

	Title string `json:"title"`
	Body  string `json:"body"`

	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
	DisabledAt *time.Time `json:"disabledAt,omitempty"`
}

// Articles is used to present a list of articles in a JSON response.
type Articles struct {
	Articles []Article `json:"articles"`
}

// ArticleCreateUpdate is the request body that is
// accepted for create and updates to articles.
type ArticleCreateUpdate struct {
	Title string `json:"title" binding:"required"`
	Body  string `json:"body" binding:"required"`
}
