package helpers

import (
	"net/http"
	"testing"

	"github.com/kyleterry/funnel/testutils"
	"github.com/stretchr/testify/assert"
)

func TestFetch(t *testing.T) {
	ts := testutils.StartTestServer()
	defer testutils.StopTestServer(ts)

	response, err := Fetch(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, response.StatusCode, http.StatusOK)
}
