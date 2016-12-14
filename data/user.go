package data

import "github.com/jinzhu/gorm"

// User represents the loginable users
type User struct {
	gorm.Model

	Email    string `gorm:"type:varchar(100);unique"`
	Password string
	Admin    bool
	Feeds    []Feed `gorm:"many2many:subscriptions;"`
}
