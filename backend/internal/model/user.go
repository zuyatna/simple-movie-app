package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid; primaryKey; default:gen_random_uuid()"`
	Username  string    `json:"username" gorm:"type:varchar(20); unique; not null"`
	Password  string    `json:"password" gorm:"type:varchar(20); not null"`
	Role      string    `json:"role" gorm:"type:varchar(10); default:'user'"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
