package services

import (
	"os"
	"testing"

	"github.com/kyleterry/funnel/data"
	"github.com/kyleterry/funnel/testutils"
)

var db *data.FunnelDB

func setup() {
	var err error

	db, err = data.New(":memory:", true)
	if err != nil {
		panic(err)
	}

	testutils.Init(db)
}

func teardown() {
	err := db.Close()
	db = nil
	if err != nil {
		panic(err)
	}
}

func run(m *testing.M) int {
	setup()
	defer teardown()

	return m.Run()
}

func TestMain(m *testing.M) {
	os.Exit(run(m))
}
