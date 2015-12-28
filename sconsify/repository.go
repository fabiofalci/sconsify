package sconsify

import (
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/fabiofalci/sconsify/infrastructure"
)

func InitRepository() {
	db, _ := gorm.Open("sqlite3", infrastructure.GetDatabaseFileLocation())
	db.DB()
	db.DB().Ping()
	db.DB().SetMaxIdleConns(1)
	db.DB().SetMaxOpenConns(2)

	db.AutoMigrate(&Album{}, &Artist{}, &Playlist{}, &Track{})
}

