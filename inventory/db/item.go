package db

import (
	"fmt"
	"log"
	"strings"

	"github.com/jinzhu/gorm"
)

type Item struct {
	gorm.Model
	Name     string  `gorm:"unique;not null;index:item_name_idx`
	PriceMin float64 `gorm:"not null"`
	PriceMax float64 `gorm:"not null"`
}

// gorm tablename
func (i *Item) TableName() string {
	return "items"
}

// implement Stringer interface
func (i *Item) String() string {
	return fmt.Sprintf("Item(%s [%f, %f])", i.Name, i.PriceMin, i.PriceMax)
}

// get an item with the args
func NewItem(name string, min float64, max float64) *Item {
	return &Item{
		Name:     name,
		PriceMin: min,
		PriceMax: max,
	}
}

// create item
func (i *Item) Create() bool {
	if !db.NewRecord(*i) {
		log.Println(PrimaryKeyCollusionErr)
		return false
	}
	db.Create(i)
	if db.NewRecord(*i) {
		log.Println(UniqueCollusionErr)
		return false
	}
	return true

}

// get existing item record
// return nil if not exists
func GetItem(name string) *Item {
	items := make([]Item, 1)
	db.Take(&items, "name LIKE ?", strings.ToLower(name))
	if 0 < len(items) {
		return &items[0]
	}
	return nil
}

// del item from database permanently
func (i *Item) PermDel() {
	db.Unscoped().Delete(Item{}, "name LIKE ?", i.Name)
}

// write item to db as is
func (i *Item) Update() {
	db.Save(i)
}

// query same item and update by copying new object by pointer
func (i *Item) PullUpdate() {
	updated := GetItem(i.Name)
	if updated == nil {
		log.Panicf("can't update %v, no record exists w/ email", i)
	}
	*i = *updated
}
