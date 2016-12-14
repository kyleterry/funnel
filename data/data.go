package data

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite" // Needed for the sqlite driver
)

// FunnelDB wraps a gorm.DB connection and provides a useful interface for
// interacting with the database.
type FunnelDB struct {
	Conn *gorm.DB
}

// New returns a new instance of FunnelDB or an error
func New(path string, debug bool) (*FunnelDB, error) {
	conn, err := gorm.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	conn.LogMode(debug)

	fdb := &FunnelDB{
		Conn: conn,
	}

	// Temporary migrations
	fdb.Conn.AutoMigrate(
		&Feed{},
		&Item{},
		&User{},
		&Subscription{},
		&SubscriptionCategory{},
		&ItemMeta{},
		&Settings{})

	return fdb, nil
}

func (f *FunnelDB) Close() error {
	return f.Conn.Close()
}
