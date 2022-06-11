package db

import (
	"errors"
	"log"
	"os"
	"path"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var db *gorm.DB

var (
	// TODO: sqliteFilename reuseable
	// sqlite
	sqliteFilename = "inventory_name.db"
)

// errors
var (
	PrimaryKeyCollusionErr = errors.New("Primary key not empty")
	UniqueCollusionErr     = errors.New("Unique identifier collusion")
)

// initializes development database
// may use sqlite at production on early phase
func InitSQLiteDB() {
	// set dbpath to $HOME/$sqliteFilename
	home, err := os.UserHomeDir()
	if err != nil {
		log.Panicf("can't get home folder path: %v", err)
	}
	sqlitePath := path.Join(home, sqliteFilename)
	db, err = gorm.Open("sqlite3", sqlitePath)
	if err != nil {
		log.Panicf("can't open database(%v): %v", sqlitePath, err)
	}
}

// close all connections
func CloseDB() {
	db.Close()
}

// initialize driver
func init() {
	InitSQLiteDB()
	db.AutoMigrate(&Item{})
}
