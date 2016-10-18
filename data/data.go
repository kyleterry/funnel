package data

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type FunnelDB struct {
	Conn *gorm.DB
}

func New(path string) (*FunnelDB, error) {
	conn, err := gorm.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	fdb := &FunnelDB{
		Conn: conn,
	}

	// Temporary migrations
	fdb.Conn.AutoMigrate(
		&Feed{},
		&Item{},
		&User{},
		&FeedCategory{})

	return fdb, nil
}
