package data

import "time"

type SubscriptionCategory struct {
	FeedID     uint
	CategoryID uint
}

type Subscription struct {
	FeedID     uint
	UserID     uint
	Categories []Category `gorm:"many2many:subscription_categories"`
}

type ItemMeta struct {
	ItemID   uint
	UserID   uint
	User     *User
	Saved    bool
	ReadTime time.Time `gorm:"null"`
}
