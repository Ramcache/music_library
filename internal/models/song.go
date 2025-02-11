package models

import (
	"time"
)

type Song struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Group       string    `gorm:"size:255" json:"group"`
	Song        string    `gorm:"size:255" json:"song"`
	ReleaseDate string    `gorm:"size:50" json:"release_date"`
	Text        string    `gorm:"type:text" json:"text"`
	Link        string    `gorm:"type:text" json:"link"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
