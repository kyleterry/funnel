package services

import (
	"github.com/jinzhu/gorm"
	"github.com/kyleterry/funnel/data"
	"github.com/pkg/errors"
)

type ItemService struct {
	db *data.FunnelDB
}

func NewItemService(db *data.FunnelDB) *ItemService {
	return &ItemService{db: db}
}

func (s *ItemService) Exists(item *data.Item) (bool, error) {
	if err := s.db.Conn.Where("link = ? AND feed_id = ?", item.Link, item.FeedID).Find(&item).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *ItemService) Create(item *data.Item) (bool, error) {
	exists, err := s.Exists(item)
	if err != nil {
		return false, errors.Wrap(err, "failed to check if item exists")
	}
	if exists {
		return false, nil
	}
	if err := s.db.Conn.Create(item).Error; err != nil {
		return false, err
	}

	return true, nil
}

func (s *ItemService) Deliver(item *data.Item) error {
	var subs []*data.Subscription

	if err := s.db.Conn.Where("feed_id = ?", item.FeedID).Find(&subs).Error; err != nil {
		return err
	}

	for _, sub := range subs {
		itemMeta := data.ItemMeta{
			UserID: sub.UserID,
			ItemID: item.ID,
		}
		if err := s.db.Conn.Create(&itemMeta).Error; err != nil {
			return err
		}
	}

	return nil
}

func (s *ItemService) Delete(item *data.Item) error {
	return s.db.Conn.Delete(item).Error
}

func (s *ItemService) GetAllForUser(user *data.User) ([]data.Item, error) {
	var items []data.Item

	err := s.db.Conn.
		Joins("JOIN subscriptions ON subscriptions.user_id = ?", user.ID).
		Joins("JOIN feeds ON subscriptions.feed_id = feeds.id").
		Where("items.feed_id = feeds.id").Find(&items).Error

	if err != nil {
		return nil, err
	}

	return items, nil
}
