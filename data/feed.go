package data

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Feed struct {
	gorm.Model

	Title        string
	URL          string
	SiteURL      string
	Favicon      string
	Domain       string
	RequiresAuth bool
	Auth         string
	HotLinking   bool
	Unread       int
	LastFetch    time.Time
	Categories   []Category `gorm:"many2many:feed_categories;"`
}
