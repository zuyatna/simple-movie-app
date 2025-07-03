package model

type Movie struct {
	ID          int     `json:"id" gorm:"primaryKey"`
	Title       string  `json:"title" gorm:"not null"`
	Description string  `json:"description"`
	Year        int     `json:"year"`
	Rating      float64 `json:"rating"`
	Poster      []byte  `json:"-" gorm:"type:bytea"`
	PosterURL   string  `json:"poster_url" gorm:"-"`
	CreatedAt   string  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   string  `json:"updated_at" gorm:"autoUpdateTime"`
}
