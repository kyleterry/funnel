package worker

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/kyleterry/funnel/config"
	"github.com/kyleterry/funnel/data"
	"github.com/kyleterry/funnel/services"
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
	feedService := services.NewFeedService(w.db)

	for {
		log.Println("Fetching feeds")
		feeds, err := feedService.GetFeeds()
		if err != nil {
			log.Print(err)
		} else {
			wg := sync.WaitGroup{}
			for _, feed := range feeds {
				log.Printf("[Start] Fetching %s", feed.URL)
				wg.Add(1)
				go w.Work(&wg, feed)
			}
			wg.Wait()
		}

		<-time.After(time.Second * 10)
	}
}

func (w Worker) Work(wg *sync.WaitGroup, feed data.Feed) {
	log.Printf("[Work] Fetching %s", feed.URL)
	feedService := services.NewFeedService(w.db)
	rawFeed, err := feedService.Fetch(&feed)
	if err != nil {
		log.Printf("There was an error fetching the feed %s, %s", feed.Title, err)
		wg.Done()
		return
	}

	fmt.Println(rawFeed)

	itemService := services.NewItemService(w.db)
	for _, rawItem := range rawFeed.Items {
		item := &data.Item{
			FeedID:      feed.ID,
			Title:       rawItem.Title,
			Description: rawItem.Summary,
			Link:        rawItem.Link,
			PostedAt:    rawItem.Date,
		}

		created, err := itemService.Create(item)
		if err != nil {
			log.Printf("Error creating item: %s", err)
			wg.Done()
			return
		}

		if created {
			err = itemService.Deliver(item)
			if err != nil {
				log.Printf("Error delivering item to user: %s", err)
				wg.Done()
				return
			}
		}
	}

	wg.Done()
}
