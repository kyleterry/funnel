package worker

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/SlyMarbo/rss"
	"github.com/kyleterry/funnel/config"
	"github.com/kyleterry/funnel/data"
)

type Worker struct {
	config *config.Config
	db     *data.FunnelDB
}

func New(c *config.Config, db *data.FunnelDB) *Worker {
	w := &Worker{c, db}

	return w
}

func (w Worker) Start() {
	for {
		feeds := []data.Feed{}
		w.db.Conn.Find(&feeds)

		wg := sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)
			go w.Work(&wg, &feed)
		}
		fmt.Println(feeds)

		wg.Wait()
		<-time.After(time.Second * 5)
	}
}

func (w Worker) Work(wg *sync.WaitGroup, feed *data.Feed) {
	rawFeed, err := rss.Fetch(feed.URL)
	if err != nil {
		log.Printf("There was an error fetching the feed %s, %s", feed.Title, err)
	}

	for _, rawItem := range rawFeed.Items {
		item := &data.Item{
			FeedID:      feed.ID,
			Title:       rawItem.Title,
			Description: rawItem.Summary,
			Link:        rawItem.Link,
		}

		err := w.db.Conn.Create(item).Error

		if err != nil {
			log.Printf("There was an error fetching the feed %s, %s", feed.Title, err)
		}

	}

	wg.Done()
}
