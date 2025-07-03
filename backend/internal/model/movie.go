package model

import (
	"time"

	"gorm.io/gorm"
)

type Movie struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Title       string         `json:"title" gorm:"not null"`
	Description string         `json:"description" gorm:"type:text"`
	Year        int            `json:"year" gorm:"type:integer"`
	Rating      float64        `json:"rating" gorm:"type:float"`
	Poster      []byte         `json:"-" gorm:"type:bytea"`
	PosterURL   string         `json:"poster_url" gorm:"-"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}
