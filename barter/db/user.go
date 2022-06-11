package db

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Email string `gorm:"unique;not null;index:users_email_idx" json:"email"`
	Key   string `gorm:"unique;not null;index:users_key_idx" json:"key"`
}

// gorm tablename
func (u *User) TableName() string {
	return "users"
}

// implement Stringer interface
func (u *User) String() string {
	return fmt.Sprintf("User(email: %s, key: %s)", u.Email, u.Key)
}

// get a user with the args
func NewUser(email string) *User {
	return &User{
		Email: email,
		Key:   uuid.New().String(),
	}
}

// create user
func (u *User) Create() bool {
	if !db.NewRecord(*u) {
		log.Println(PrimaryKeyCollusionErr)
		return false
	}
	db.Create(u)
	if db.NewRecord(*u) {
		log.Println(UniqueCollusionErr)
		return false
	}
	return true

}

// get existing user
// return nil if not exists
func GetUser(key string) *User {
	users := make([]User, 1)
	db.Take(&users, "key = ?", key)
	if 0 < len(users) {
		return &users[0]
	}
	return nil
}

// del user from database permanently
func (u *User) PermDel() {
	db.Unscoped().Delete(User{}, "email LIKE ?", u.Email)
}

// write user to db as is
func (u *User) Update() {
	db.Save(u)
}

// query same user and update by copying new object by pointer
func (u *User) PullUpdate() {
	updated := GetUser(u.Email)
	if updated == nil {
		log.Panicf("can't update %v, no record exists w/ email", u)
	}
	*u = *updated
}
