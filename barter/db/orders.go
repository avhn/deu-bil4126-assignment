package db

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
)

type Order struct {
	gorm.Model
	Email           string `gorm:"not null;index:orders_email_idx" json:"email"`
	Key             string `gorm:"not null;" json:"key"`
	GivenInventory  string `gorm:"not null" json:"given_inventory"`
	GivenItem       string `gorm:"not null;index:orders_given_item_idx" json:"given_item"`
	GivenAmount     int    `gorm:"not null" json:"given_amount"`
	WantedInventory string `gorm:"not null" json:"wanted_inventory"`
	WantedItem      string `gorm:"not null;index:orders_wanted_item_idx" json:"wanted_item"`
	WantedAmount    int    `gorm:"not null" json:"wanted_amount"`
}

// gorm tablename
func (r *Order) TableName() string {
	return "orders"
}

// implement Stringer interface
func (r *Order) String() string {
	return fmt.Sprintf("Order(%s[key:%s] wants %d %s from %s inventory for %d %s from %s inventory)",
		r.Email, r.Key,
		r.WantedAmount, r.WantedItem, r.WantedInventory,
		r.GivenAmount, r.GivenItem, r.GivenInventory)
}

// get a user with the args
func NewOrder(email string, key string, given_inventory string, given_item string, given_amount int,
	wanted_inventory string, wanted_item string, wanted_amount int) *Order {
	return &Order{
		Email:           email,
		Key:             key,
		GivenInventory:  given_inventory,
		GivenItem:       given_item,
		GivenAmount:     given_amount,
		WantedInventory: wanted_inventory,
		WantedItem:      wanted_item,
		WantedAmount:    wanted_amount,
	}
}

// create user
func (r *Order) Create() bool {
	if !db.NewRecord(*r) {
		log.Println(PrimaryKeyCollusionErr)
		return false
	}
	db.Create(r)
	if db.NewRecord(*r) {
		log.Println(UniqueCollusionErr)
		return false
	}
	return true

}

func GetOrders(u *User) []Order {
	var orders []Order
	db.Take(&orders, "key = ?", u.Key)
	if 0 < len(orders) {
		return orders
	}
	return nil
}

func FindOrders(wanted_inventory string, wanted_name string,
	given_inventory string, given_name string) []Order {
	var orders []Order
	db.Take(&orders, "given_inventory LIKE ? AND given_item LIKE ? AND wanted_inventory LIKE ? AND wanted_item LIKE ?",
		given_inventory, given_name, wanted_inventory, wanted_name)
	if 0 < len(orders) {
		return orders
	}
	return nil
}

func (r *Order) PermDel() {
	db.Unscoped().Delete(Order{}, "id = ?", r.ID)
}

func (r *Order) Update() {
	db.Save(r)
}

func (r *Order) PullUpdate() {
	updated := make([]Order, 1)
	db.Take(&updated, "id = ?", r.ID)
	if 0 < len(updated) {
		*r = updated[0]
	}
}
