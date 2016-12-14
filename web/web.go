package web

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/kyleterry/funnel/config"
	"github.com/kyleterry/funnel/data"
)

type Funnel struct {
	router *mux.Router
	config *config.Config
	db     *data.FunnelDB
}

func New(c *config.Config, db *data.FunnelDB) *Funnel {
	f := &Funnel{
		config: c,
		router: mux.NewRouter(),
		db:     db,
	}

	f.init()

	return f
}

func (f *Funnel) init() {
	all := alice.New(LoggingHandler)

	initTemplates(f)

	f.router.Handle("/", all.Then(errorHandler(f.IndexHandler))).Methods("GET").Name("index")

	feedsrouter := f.router.PathPrefix("/feeds").Subrouter()
	feedsrouter.Handle("/", all.Then(errorHandler(f.ListFeedHandler))).
		Methods("GET").
		Name("feed-index")
	feedsrouter.Handle("/new", all.Then(errorHandler(f.NewFeedHandler))).
		Methods("GET").
		Name("feeds-new")
	feedsrouter.Handle("/create", all.Then(errorHandler(f.CreateFeedHandler))).
		Methods("POST").
		Name("feeds-create")
	feedsrouter.Handle("/{id:[0-9]+}", all.Then(errorHandler(f.ViewFeedHandler))).
		Methods("GET").
		Name("feeds-view")

	f.router.PathPrefix("/static").Handler(staticHandler)

	f.router.NotFoundHandler = errorHandler(func(w http.ResponseWriter, r *http.Request) error {
		w.WriteHeader(http.StatusNotFound)
		return renderTemplate(w, "404", nil)
	})
}

func (f *Funnel) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f.router.ServeHTTP(w, r)
}

type errorHandler func(http.ResponseWriter, *http.Request) error

func (fn errorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := fn(w, r)
	if err != nil {
		log.Printf("Got error while processing the request: %s\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
