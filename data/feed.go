package data

import (
	"time"

	"github.com/jinzhu/gorm"
)

// Feed is the a url pointing to an RSS or Atom feed
type Feed struct {
	gorm.Model

	Title        string
	URL          string `gorm:"not null:unique"`
	Description  string
	SiteURL      string
	Favicon      string
	Domain       string
	RequiresAuth bool
	Auth         string
	HotLinking   bool
	Unread       int
	LastFetch    time.Time
	AddedByID    uint
	AddedBy      *User
}
