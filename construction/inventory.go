package construction

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
	ItemName string  `json:"item"`
	Amount   int     `json:"amount"`
	Budget   float64 `json:"budget"`
}

type RawItem struct {
	Name     string  `json:"name"`
	PriceMin float64 `json:"price_min"`
	PriceMax float64 `json:"price_max"`
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
		return
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
		return
	}
	// passed
	w.WriteHeader(http.StatusOK)
}

// list all items in the inventory
func ListAll(w http.ResponseWriter, r *http.Request) {
	items := db.GetAllItems()
	res := make([]RawItem, len(items))
	for i := 0; i < len(items); i++ {
		res[i].Name, res[i].PriceMin, res[i].PriceMax = items[i].Name, items[i].PriceMin, items[i].PriceMax
	}
	resp, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Can't marshall response!")
		return
	}
	w.WriteHeader(http.StatusOK)
	// make resp = {"list": resp}
	w.Write(append(append([]byte(`{"list": `), resp...), []byte(`}`)...))
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
	ri := RawItem{db_i.Name, db_i.PriceMin, db_i.PriceMax}
	resp, err := json.Marshal(ri)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Can't marshall response!")
		return
	}
	w.Write(resp)
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
	db_i := db.GetItem(cr.ItemName)
	if db_i == nil {
		w.WriteHeader(http.StatusNotFound)
		log.Println("Item doesn't exists.")
		return
	}
	if cr.Amount <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Invalid amount argument, or doesn't exists.")
		return
	}
	// calculate cost
	res := db_i.PriceMin * float64(cr.Amount)
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
	db_i := db.GetItem(cr.ItemName)
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
	res := int(math.Floor(cr.Budget / db_i.PriceMin))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"amount": %d,}`, res)))
}
