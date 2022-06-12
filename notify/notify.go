package notify

// simulate email and notification
import (
	barterdb "ebarter/barter/db"
	"log"
)

type Notification struct {
	u1             *barterdb.User
	u2             *barterdb.User
	got_inventory  string
	got_item       string
	got_amount     int
	gave_inventory string
	gave_item      string
	gave_amount    int
}

// email and notify
func Notify(n Notification) {
	// mock notification
	log.Printf("%v gave %d amount %s from %s to %v and got %d amount %s from %s in return.",
		n.u1, n.got_amount, n.got_item, n.got_inventory, n.u2, n.gave_amount, n.gave_item, n.gave_inventory)
}
