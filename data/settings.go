package data

import "github.com/jinzhu/gorm"

// Settings represents the global configuration
type Settings struct {
	gorm.Model

	RefreshRate int
}
