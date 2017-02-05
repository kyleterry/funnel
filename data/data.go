package data

import (
	"github.com/jinzhu/gorm"
	"github.com/kyleterry/funnel/config"
	_ "github.com/lib/pq"           // Needed for the postgres driver
	_ "github.com/mattn/go-sqlite3" // Needed for the sqlite driver
)

// FunnelDB wraps a gorm.DB connection and provides a useful interface for
// interacting with the database.
type FunnelDB struct {
	Conn *gorm.DB
}

// New returns a new instance of FunnelDB or an error
func New(conf *config.Config) (*FunnelDB, error) {
	conn, err := gorm.Open(conf.DatabaseEngine, conf.ConnectionString)
	if err != nil {
		return nil, err
	}

	conn.LogMode(conf.Debug)

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
