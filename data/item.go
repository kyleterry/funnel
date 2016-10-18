package data

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Item struct {
	gorm.Model

	FeedID      uint
	Feed        *Feed
	UID         string
	Title       string
	Author      string
	Description string
	Link        string
	Saved       bool
	ReadTime    time.Time
}
