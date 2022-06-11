package notify

// simulate email and notification
import (
	"ebarter/barter/db"
	"fmt"
)

// email and notify
func Notify(u *db.User) {
	fmt.Fprintf("%s", u.Email)
}
