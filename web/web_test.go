package web

import (
	"net/http/httptest"
	"os"
	"testing"

	"github.com/kyleterry/funnel/config"
	"github.com/kyleterry/funnel/data"
)

var server *httptest.Server

func setup() error {
	db, err := data.New(":memory:", true)
	if err != nil {
		return err
	}
	app := New(&config.Config{}, db)
	server = httptest.NewServer(app)

	return nil
}

func teardown() {
	if server != nil {
		server.Close()
		server = nil
	}
}

func run(m *testing.M) int {
	err := setup()
	if err != nil {
		panic(err)
	}
	defer teardown()

	return m.Run()
}

func TestMain(m *testing.M) {
	os.Exit(run(m))
}
