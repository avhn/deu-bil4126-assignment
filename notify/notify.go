package notify

// simulate email and notification
import (
	barterdb "ebarter/barter/db"
	"log"
)

type Notification struct {
	U1            *barterdb.User
	U2            *barterdb.User
	GotInventory  string
	GotItem       string
	GotAmount     int
	GaveInventory string
	GaveItem      string
	GaveAmount    int
}

// email and notify
func Notify(n Notification) {
	// mock notification
	log.Printf("%v gave %d amount %s from %s to %v and got %d amount %s from %s in return.",
		n.U1, n.GotAmount, n.GotItem, n.GotInventory, n.U2, n.GaveAmount, n.GaveItem, n.GaveInventory)
}
