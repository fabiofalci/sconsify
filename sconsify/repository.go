package sconsify

import (
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/fabiofalci/sconsify/infrastructure"
)

var DB gorm.DB

var Create chan interface{}

func InitRepository() {
	DB, _ = gorm.Open("sqlite3", infrastructure.GetDatabaseFileLocation())
	DB.DB().SetMaxIdleConns(1)
	DB.DB().SetMaxOpenConns(2)

	DB.DropTable(&Album{}, &Artist{}, &Playlist{}, &Track{})
	DB.CreateTable(&Album{}, &Artist{}, &Playlist{}, &Track{})

	Create = make(chan interface{})
	go func() {
		for {
			select {
			case toCreate := <-Create:
				DB.Create(toCreate)
			}
		}
	}()
}

