package testutils

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/kyleterry/funnel/data"
	"github.com/modocache/gory"
)

var db *data.FunnelDB

func Init(fdb *data.FunnelDB) {
	db = fdb

	gory.Define("feed", data.Feed{}, func(factory gory.Factory) {
		factory["Title"] = gory.Sequence(func(n int) interface{} {
			return fmt.Sprintf("Example Title %d", n)
		})
		factory["URL"] = gory.Sequence(func(n int) interface{} {
			return fmt.Sprintf("http://example.com/%d/feed.rss", n)
		})
		factory["Description"] = "Very tasty description of a feed"
		factory["SiteURL"] = "http://example.com"
		factory["Favicon"] = ""
		factory["Domain"] = "example.com"
		factory["RequiresAuth"] = false
		factory["Auth"] = ""
		factory["HotLinking"] = true
		factory["Unread"] = 0
		factory["LastFetch"] = time.Now()
	})

	gory.Define("item", data.Item{}, func(factory gory.Factory) {
		factory["Title"] = gory.Sequence(func(n int) interface{} {
			return fmt.Sprintf("Example Item Title %d", n)
		})
		factory["Link"] = gory.Sequence(func(n int) interface{} {
			return fmt.Sprintf("http://example.com/%d/view", n)
		})
		factory["Description"] = "Very tasty description of an item"
		factory["Author"] = "Kyle Terry"
		factory["Feed"] = gory.Build("feed").(*data.Feed)
	})

	gory.Define("user", data.User{}, func(factory gory.Factory) {
		factory["Email"] = gory.Sequence(func(n int) interface{} {
			return fmt.Sprintf("user%d@example.com", n)
		})
		factory["Password"] = ""
		factory["Admin"] = false
	})
}

func Build(name string) interface{} {
	fixture := gory.Build(name)
	if fixture == nil {
		panic("Unknown fixture: " + name)
	}

	db.Conn.Debug().Create(fixture)
	err := db.Conn.Error
	if err != nil {
		panic(err)
	}

	return fixture
}

func BuildWithParams(name string, p map[string]interface{}) interface{} {
	fixture := gory.BuildWithParams(name, gory.Factory(p))
	if fixture == nil {
		panic("Unknown fixture: " + name)
	}

	db.Conn.Debug().Create(fixture)
	err := db.Conn.Error
	if err != nil {
		panic(err)
	}

	return fixture
}

func StartTestServer() *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `
<?xml version="1.0" encoding="UTF-8" ?>
<rss version="2.0">
    <channel>
        <title>Example RSS</title>
        <link>http://example.com/feed.rss</link>
        <description>Such a smooooooth description</description>
        <item>
            <title>Hi</title>
            <link>http://www.w3schools.com/xml/xml_rss.asp</link>
            <description>New RSS tutorial on W3Schools</description>
        </item>
        <item>
            <title>Hey</title>
            <link>http://www.w3schools.com/xml</link>
            <description>New XML tutorial on W3Schools</description>
        </item>
    </channel>
</rss>
	`)
	}))

	return ts
}

func StopTestServer(ts *httptest.Server) {
	ts.Close()
}
