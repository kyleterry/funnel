package web

import (
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/elazarl/go-bindata-assetfs"
	"github.com/gorilla/schema"
	"github.com/kyleterry/funnel/data"
	"github.com/kyleterry/funnel/services"
)

var staticHandler = http.FileServer(
	&assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, AssetInfo: AssetInfo},
)

func (f Funnel) IndexHandler(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (f Funnel) NewFeedHandler(w http.ResponseWriter, r *http.Request) error {
	return renderTemplate(w, "feed-new", nil)
}

func (f Funnel) ListFeedHandler(w http.ResponseWriter, r *http.Request) error {
	feedService := services.NewFeedService(f.db)
	feeds, err := feedService.GetFeeds()
	if err != nil {
		return err
	}

	var context = struct {
		Feeds []data.Feed
	}{
		Feeds: feeds,
	}

	return renderTemplate(w, "feed-index", context)
}

func (f Funnel) CreateFeedHandler(w http.ResponseWriter, r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	var form struct {
		URL string `schema:"url"`
	}

	decoder := schema.NewDecoder()

	if err := decoder.Decode(&form, r.PostForm); err != nil {
		return err
	}

	if !govalidator.IsURL(form.URL) {
		panic("I hate myself")
	}

	feedService := services.NewFeedService(f.db)

	feed, err := feedService.Create(form.URL)
	if err != nil {
		return err
	}

	http.Redirect(w, r, f.reverse("feeds-view", "id", feed.ID), http.StatusSeeOther)
	return nil
}

func (f Funnel) ViewFeedHandler(w http.ResponseWriter, r *http.Request) error {
	return renderTemplate(w, "feed-view", nil)
}

func (f Funnel) TimelineHandler(w http.ResponseWriter, r *http.Request) error {
	itemService := NewItemService(f.db)
	items, err := itemService.GetAllForUser()
}
