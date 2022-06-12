package barter

import (
	"ebarter/barter/db"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type Credentials struct {
	Email string `json:"email"`
	Key   string `json:"key"`
}

var (
	userAlreadyExistsResponse       = []byte(`{"message": "User already exists"}`)
	invalidCredentialFormatResponse = []byte(`{"message": "Invalid credential format"}`)
)

func Signup(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error while reading request body: %v", err)
		return
	}
	var c Credentials
	err = json.Unmarshal(body, &c)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(invalidCredentialFormatResponse)
		log.Printf("error while unmarshaling credentials: %v", err)
		return
	}
	u := db.NewUser(c.Email)
	if !u.Create() {
		w.WriteHeader(http.StatusConflict)
		w.Write(userAlreadyExistsResponse)
		log.Printf("User already exists.")
		return
	}
	c.Key = u.Key
	resp, err := json.Marshal(c)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Can't marshall response!")
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}
