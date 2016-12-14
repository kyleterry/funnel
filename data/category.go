package data

import "github.com/jinzhu/gorm"

// Category represents a group of feeds
type Category struct {
	gorm.Model

	Name   string
	User   *User
	UserID uint64 `gorm:"index"`
	Feeds  []Feed `gorm:"many2many:feed_categories;"`
}
