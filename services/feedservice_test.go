package services

import (
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/kyleterry/funnel/data"
	"github.com/kyleterry/funnel/testutils"
	"github.com/stretchr/testify/assert"
)

func TestFeedExists(t *testing.T) {
	var cases = []struct {
		feed   *data.Feed
		exists bool
	}{
		{testutils.Build("feed").(*data.Feed), true},
		{&data.Feed{URL: "http://example.com/hi/feed.rss"}, false},
	}

	feedService := NewFeedService(db)

	for _, c := range cases {
		res, err := feedService.Exists(c.feed)
		assert.NoError(t, err)
		assert.Equal(t, c.exists, res)
	}
}

func TestSubscribe(t *testing.T) {
	feed := testutils.Build("feed").(*data.Feed)
	user := testutils.Build("user").(*data.User)

	feedService := NewFeedService(db)
	_, err := feedService.SubscribeUser(user, feed)
	assert.NoError(t, err)

	// refetch the user
	db.Conn.First(&user, &user.ID)
	db.Conn.Model(&user).Association("Feeds").Find(&user.Feeds)

	assert.Len(t, user.Feeds, 1)
	assert.Equal(t, feed.Title, user.Feeds[0].Title)
}

func TestCreateFeed(t *testing.T) {
	ts := testutils.StartTestServer()
	defer testutils.StopTestServer(ts)

	var cases = []struct {
		url        string
		shouldFail bool
		errString  string
	}{
		{ts.URL, false, ""},
		{"https://www.reddit.com/.rss", false, ""},
		{"https://www.reddit.com/r/golang.rss", false, ""},
		{"https://google.com", true, ""},
		{"cats://bacon", true, "Cannot add a non-http style feed url"},
		{"ftp://1998.com", true, "Cannot add a non-http style feed url"},
		{"http:/www.reddit.com/.rss", true, ""},
	}

	feedService := NewFeedService(db)

	for _, c := range cases {
		feed, err := feedService.Create(c.url)
		if c.shouldFail {
			assert.Error(t, err)
			if c.errString != "" {
				assert.Equal(t, c.errString, err.Error())
			}
		} else {
			if assert.NoError(t, err, c.url) {
				assert.Equal(t, c.url, feed.URL)
			}
		}
	}

	// Now lets simulate a DB failure with an empty DB and no tables
	hold := db
	defer func() {
		db = hold
	}()

	conn, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		assert.NoError(t, err)
	}

	db = &data.FunnelDB{Conn: conn}
	_, err = feedService.Create(ts.URL)
	assert.Error(t, err)
}

func TestGetFeeds(t *testing.T) {
	testutils.Build("feed")
	testutils.Build("feed")

	feedService := NewFeedService(db)
	feeds, err := feedService.GetFeeds()
	assert.NoError(t, err)
	assert.True(t, len(feeds) > 0)
}
