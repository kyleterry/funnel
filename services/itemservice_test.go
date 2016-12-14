package services

import (
	"fmt"
	"testing"

	"github.com/Xe/uuid"
	"github.com/kyleterry/funnel/data"
	"github.com/kyleterry/funnel/testutils"
	"github.com/stretchr/testify/assert"
)

func TestItemExists(t *testing.T) {
	var cases = []struct {
		item   *data.Item
		exists bool
	}{
		{testutils.Build("item").(*data.Item), true},
		{&data.Item{Link: "https://example.com/does-not-exist-yet"}, false},
	}

	itemService := NewItemService(db)

	for _, c := range cases {
		res, err := itemService.Exists(c.item)
		assert.NoError(t, err)
		assert.Equal(t, c.exists, res)
	}
}

func TestCreateItem(t *testing.T) {
	feed := testutils.Build("feed").(*data.Feed)

	itemService := NewItemService(db)

	url := fmt.Sprintf("http://example.com/feed/%s.html", uuid.New())

	item := &data.Item{
		Feed:        feed,
		Title:       "A test item 4000",
		Author:      "Kyle Terry",
		Description: "A test item description",
		Link:        url,
	}

	created, err := itemService.Create(item)
	assert.True(t, created)
	assert.NoError(t, err)

	exists, err := itemService.Exists(item)
	assert.NoError(t, err)
	assert.True(t, exists)

	// NO DUPES!
	item.ID = 0

	created, err = itemService.Create(item)
	assert.False(t, created)
}

func TestItemDelivery(t *testing.T) {
	item := testutils.Build("item").(*data.Item)
	user := testutils.Build("user").(*data.User)

	feedService := NewFeedService(db)
	_, err := feedService.SubscribeUser(user, item.Feed)
	assert.NoError(t, err)

	itemService := NewItemService(db)
	err = itemService.Deliver(item)
	assert.NoError(t, err)

	itemMeta := data.ItemMeta{}
	err = db.Conn.Where("user_id = ? AND item_id = ?", user.ID, item.ID).Find(&itemMeta).Error
	assert.NoError(t, err)

	assert.Equal(t, user.ID, itemMeta.UserID)
	assert.Equal(t, item.ID, itemMeta.ItemID)
	assert.Equal(t, false, itemMeta.Saved)
}

func TestItemDeletion(t *testing.T) {
	item := testutils.Build("item").(*data.Item)

	itemService := NewItemService(db)

	err := itemService.Delete(item)
	assert.NoError(t, err)

	assert.True(t, db.Conn.First(&item, item.ID).RecordNotFound())
}

func TestGetItemsForUser(t *testing.T) {
	user := testutils.Build("user").(*data.User)
	feed1 := testutils.Build("feed").(*data.Feed)
	// feed2 := testutils.Build("feed").(*data.Feed)

	feedService := NewFeedService(db)
	itemService := NewItemService(db)

	_, err := feedService.SubscribeUser(user, feed1)
	assert.NoError(t, err)

	for i := 0; i < 10; i++ {
		item := testutils.BuildWithParams("item", map[string]interface{}{
			"FeedID": feed1.ID,
		}).(*data.Item)
		err := itemService.Deliver(item)
		assert.NoError(t, err)
	}

	// check to see if the method can only fetch items for feed 1
	items, err := itemService.GetAllForUser(user)
	assert.NoError(t, err)
	assert.Len(t, items, 10)
}
