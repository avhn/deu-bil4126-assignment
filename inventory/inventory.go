package inventory

import (
	"ebarter/inventory/db"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
)

type CalculateRequest struct {
	WantedItem   string  `json:"wanted_item"`
	WantedAmount int     `json:"wanted_amount"`
	Budget       float64 `json:"budget"`
}

// inventory methods and inventory reuseable interface
// add new item
func Add(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error while reading request body: %v", err)
		return
	}
	var i db.Item
	err = json.Unmarshal(body, &i)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("error while unmarshaling item: %v", err)
		return
	}
	check := db.GetItem(i.Name)
	if check != nil {
		w.WriteHeader(http.StatusConflict)
		log.Println("Item already exists.")
		return
	}
	// check conditions
	if i.PriceMax < i.PriceMin || i.PriceMin <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Didn't satisfy logical checks.")
		return
	}
	if !i.Create() {
		// not created
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Could not create.")
		return
	}
	// success
	w.WriteHeader(http.StatusCreated)
}

// update existing item
func Update(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error while reading request body: %v", err)
		return
	}
	var i db.Item
	err = json.Unmarshal(body, &i)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("error while unmarshaling item: %v", err)
		return
	}
	db_i := db.GetItem(i.Name)
	if db_i == nil {
		w.WriteHeader(http.StatusNotFound)
		log.Println("Item doesn't exists.")
		return
	}
	db_i.Name, db_i.PriceMin, db_i.PriceMax = i.Name, i.PriceMin, i.PriceMax
	db_i.Update()
	// check update
	i = *db_i
	check := db.GetItem(i.Name)
	if check == nil || check.PriceMax != i.PriceMax || check.PriceMin != check.PriceMin {
		// not updated
		w.WriteHeader(http.StatusExpectationFailed)
	}
	// passed
	w.WriteHeader(http.StatusOK)
}

// delete existing item
func Del(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error while reading request body: %v", err)
		return
	}
	var i db.Item
	err = json.Unmarshal(body, &i)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("error while unmarshaling item: %v", err)
		return
	}
	db_i := db.GetItem(i.Name)
	if db_i == nil {
		w.WriteHeader(http.StatusNotFound)
		log.Println("Item doesn't exists.")
		return
	}
	db_i.PermDel()
	// check
	if db.GetItem(db_i.Name) != nil {
		// not deleted
		w.WriteHeader(http.StatusExpectationFailed)
	}
	// passed
	w.WriteHeader(http.StatusOK)
}

// list all items in the inventory
func ListAll(w http.ResponseWriter, r *http.Request) {
	res := db.Find(&db.Item)
	fmt.Println(res)
}

// check inventory for item
func Check(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error while reading request body: %v", err)
		return
	}
	var i db.Item
	err = json.Unmarshal(body, &i)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("error while unmarshaling item: %v", err)
		return
	}
	db_i := db.GetItem(i.Name)
	if db_i == nil {
		w.WriteHeader(http.StatusNotFound)
		log.Println("Item doesn't exists.")
		return
	}
	// found
	w.WriteHeader(http.StatusConflict)
}

// find wanted item from db
// return wanted_item_min_price * wanted_amount
func Cost(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error while reading request body: %v", err)
		return
	}
	var cr CalculateRequest
	err = json.Unmarshal(body, &cr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("error while unmarshaling item: %v", err)
		return
	}
	db_i := db.GetItem(cr.WantedItem)
	if db_i == nil {
		w.WriteHeader(http.StatusNotFound)
		log.Println("Item doesn't exists.")
		return
	}
	if cr.WantedAmount <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Invalid amount argument, or doesn't exists.")
		return
	}
	// calculate cost
	res := db_i.PriceMin * float64(cr.WantedAmount)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"cost": %f,}`, res)))
}

// find wanted item from db
// use floor division while calculating
// return (budget // wanted_item_min_price)
func Calculate(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error while reading request body: %v", err)
		return
	}
	var cr CalculateRequest
	err = json.Unmarshal(body, &cr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("error while unmarshaling item: %v", err)
		return
	}
	db_i := db.GetItem(cr.WantedItem)
	if db_i == nil {
		w.WriteHeader(http.StatusNotFound)
		log.Println("Item doesn't exists.")
		return
	}
	if cr.Budget <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Invalid budget argument, or doesn't exists.")
		return
	}
	// calculate
	res := int(math.Floor(cr.Budget * float64(cr.WantedAmount)))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"result": %d,}`, res)))
}
