package models

import "time"

type Service struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	Name        string    `gorm:"not null" json:"name"`
	URL         string    `gorm:"not null" json:"url"`
	Icon        string    `json:"icon"`
	Description string    `json:"description"`
	GroupID     uint      `json:"group_id"`
	SortOrder   int       `json:"sort_order"`
	Target      string    `gorm:"default:_blank" json:"target"`
	StatusCheck bool      `gorm:"default:true" json:"status_check"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
