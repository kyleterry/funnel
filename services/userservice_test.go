package services

import (
	"testing"

	"github.com/kyleterry/funnel/data"
	"github.com/kyleterry/funnel/testutils"
	"github.com/stretchr/testify/assert"
)

func TestUserExists(t *testing.T) {
	var cases = []struct {
		User   *data.User
		exists bool
	}{
		{testutils.Build("user").(*data.User), true},
		{&data.User{Email: "kyle@example.com"}, false},
	}

	userService := NewUserService(db)

	for _, c := range cases {
		res, err := userService.Exists(c.User)
		assert.NoError(t, err)
		assert.Equal(t, c.exists, res)
	}
}
