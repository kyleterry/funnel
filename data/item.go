package data

import (
	"time"

	"github.com/jinzhu/gorm"
)

// Item is a single feed item. This is usually a post belonging to the feed
type Item struct {
	gorm.Model

	Feed        *Feed
	FeedID      uint `gorm:"index"`
	UID         string
	Title       string
	Author      string
	Description string
	Link        string
	PostedAt    time.Time
}
