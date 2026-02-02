package models

type Bookmark struct {
	ID        uint   `gorm:"primarykey" json:"id"`
	Name      string `gorm:"not null" json:"name"`
	URL       string `gorm:"not null" json:"url"`
	Icon      string `json:"icon"`
	SortOrder int    `json:"sort_order"`
}
