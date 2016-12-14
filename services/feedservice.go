package services

import (
	"errors"
	"net/url"

	"github.com/SlyMarbo/rss"
	"github.com/jinzhu/gorm"
	"github.com/kyleterry/funnel/data"
	"github.com/kyleterry/funnel/helpers"
)

type FeedService struct {
	db *data.FunnelDB
}

func NewFeedService(db *data.FunnelDB) *FeedService {
	return &FeedService{db: db}
}

func (s *FeedService) Exists(feed *data.Feed) (bool, error) {
	if err := s.db.Conn.Where("url = ?", feed.URL).Find(&feed).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *FeedService) SubscribeUser(user *data.User, feed *data.Feed) (*data.Subscription, error) {
	subscription := data.Subscription{}
	for _, err := range s.db.Conn.Where("user_id = ? and feed_id = ?", user.ID, feed.ID).First(&subscription).GetErrors() {
		if err == gorm.ErrRecordNotFound {
			subscription.UserID = user.ID
			subscription.FeedID = feed.ID
			if err := s.db.Conn.Create(&subscription).Error; err != nil {
				return nil, err
			}
			return &subscription, nil
		}
		return nil, err
	}
	return &subscription, nil
}

func (s *FeedService) Create(rawurl string) (*data.Feed, error) {
	purl, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}

	if purl.Scheme != "http" && purl.Scheme != "https" {
		return nil, errors.New("Cannot add a non-http style feed url")
	}

	rawFeed, err := rss.FetchByFunc(helpers.Fetch, purl.String())
	if err != nil {
		return nil, err
	}

	feed := &data.Feed{}
	feed.Title = rawFeed.Title
	feed.URL = purl.String()
	feed.Domain = purl.Host
	feed.Description = rawFeed.Description

	if err := s.db.Conn.Create(feed).Error; err != nil {
		return nil, err
	}

	return feed, nil
}

func (s *FeedService) GetFeeds() ([]data.Feed, error) {
	var feeds []data.Feed
	if err := s.db.Conn.Find(&feeds).Error; err != nil {
		return feeds, err
	}
	return feeds, nil
}

func (s *FeedService) Fetch(feed *data.Feed) (*rss.Feed, error) {
	rawFeed, err := rss.FetchByFunc(helpers.Fetch, feed.URL)
	if err != nil {
		return nil, err
	}
	return rawFeed, nil
}
