//go:generate go-bindata -o ./web/bindata.go -pkg web templates/... static/...

package main

import (
	"log"
	"net/http"

	"github.com/kyleterry/funnel/config"
	"github.com/kyleterry/funnel/data"
	"github.com/kyleterry/funnel/web"
	"github.com/kyleterry/funnel/worker"
)

func main() {
	c, err := config.New()

	if err != nil {
		log.Fatal(err)
	}

	db, err := data.New(c)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Bringing inline worker online")
	w := worker.New(c, db)
	go w.Start()

	funnel := web.New(c, db)

	log.Println("Starting web")
	log.Printf("Listening on http://%s\n", c.HTTPAddr)
	if err := http.ListenAndServe(c.HTTPAddr, funnel); err != nil {
		log.Fatal(err)
	}
}
