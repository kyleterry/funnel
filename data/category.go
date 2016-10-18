package data

import "github.com/jinzhu/gorm"

type Category struct {
	gorm.Model

	Name  string
	Feeds []Feed `gorm:"many2many:feed_categories;"`
}
