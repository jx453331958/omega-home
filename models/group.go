package models

import "time"

type Group struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Icon      string    `json:"icon"`
	SortOrder int       `json:"sort_order"`
	Columns   int       `gorm:"default:3" json:"columns"`
	Services  []Service `json:"services,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
