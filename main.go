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
	c := config.New()

	db, err := data.New(c.DatabaseFile)
	if err != nil {
		log.Fatal(err)
	}

	w := worker.New(c, db)
	go w.Start()

	funnel := web.New(c, db)

	log.Printf("Listening on http://%s\n", c.HTTPAddr)
	if err := http.ListenAndServe(c.HTTPAddr, funnel); err != nil {
		log.Fatal(err)
	}
}
