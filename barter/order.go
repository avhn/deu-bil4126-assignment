package barter

import (
	"ebarter/barter/db"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type OrderResponse struct {
	acquired_wanted_item_amount int `json:"acquired_wanted_item_amount"`
	surplus_given_item_amount   int `json:"surplus_given_item_amount"`
	inorder_given_item_amount   int `json:"inorder_given_item_amount"`
}

type ReceivedItem struct {
	Name     string  `json:"name"`
	PriceMin float64 `json:"price_min"`
	PriceMax float64 `json:"price_max"`
}

var (
	// inventories defined by name
	inventories                = []string{"food", "electronics", "construction", "art"}
	invalidOrderFormatResponse = []byte(`{"message": "Invalid order."}`)
	inventoryServicesPrefix    = "localhost:8080"
)

func inventoryExist(inventory_name string) bool {
	for _, s := range inventories {
		if s == inventory_name {
			return true
		}
	}
	return false
}

func Max(x int, y int) int {
	if x < y {
		return y
	}
	return x
}

// for a valid order:
// + amounts should to be > 0
// + inventory and item names should be true
// non-existing item or inventory orders are invalid.
func Order(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error while reading request body: %v", err)
		return
	}
	var o db.Order
	err = json.Unmarshal(body, &o)
	if err != nil || o.GivenAmount < 1 || o.WantedAmount < 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(invalidCredentialFormatResponse)
		log.Printf("error while unmarshaling: %v", err)
		return
	}
	// check user
	user := db.GetUser(o.Key)
	if user == nil || user.Email != o.Email {
		// invalid credentials
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Invalid credentials.")
		return
	}
	// normalize input strings
	o.GivenInventory = strings.ToLower(o.GivenInventory)
	o.GivenItem = strings.ToLower(o.GivenItem)
	o.WantedInventory = strings.ToLower(o.WantedInventory)
	o.WantedItem = strings.ToLower(o.WantedItem)
	// check inventory names
	if !inventoryExist(o.GivenInventory) || !inventoryExist(o.WantedInventory) {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Inventory doesn't exist.")
		return
	}
	// ask inventories
	respGiven, err := http.Get(strings.Join([]string{inventoryServicesPrefix, o.GivenInventory, "check", "?name=" + o.GivenItem}, "/"))
	respWanted, err2 := http.Get(strings.Join([]string{inventoryServicesPrefix, o.WantedInventory, "check", "?name=" + o.WantedItem}, "/"))
	if err != nil || err2 != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error while requesting from inventory: %v | %v", err, err2)
	}
	bodyGiven, err := ioutil.ReadAll(respGiven.Body)
	bodyWanted, err2 := ioutil.ReadAll(respWanted.Body)
	if err != nil || err2 != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error while reading respons body: %v | %v", err, err2)
	}
	var givenItem ReceivedItem
	var wantedItem ReceivedItem
	err = json.Unmarshal(bodyGiven, &givenItem)
	err2 = json.Unmarshal(bodyWanted, &wantedItem)
	if err != nil || err2 != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("error while unmarshaling item from resp: %v | %v", err, err2)
		return
	}
	// if maximum cost of given items < minimum cost of wanted items
	if givenItem.PriceMax*float64(o.GivenAmount) < wantedItem.PriceMin*float64(o.WantedAmount) {
		// not acceptable return the offer
		w.WriteHeader(http.StatusNotAcceptable)
		log.Println("Order not acceptable.")
		return
	}
	// create but don't write newOrder
	newOrder := db.NewOrder(user.Email, user.Key,
		o.GivenInventory, givenItem.Name, o.GivenAmount,
		o.WantedInventory, wantedItem.Name, o.WantedAmount)
	// match order
	//var Response OrderResponse
	//Response.inorder_given_item_amount = o.GivenAmount
	orders := db.FindOrders(o.GivenInventory, givenItem.Name,
		o.WantedInventory, wantedItem.Name)
	if orders == nil { // counter offer doesn't exist
		if !newOrder.Create() {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("Can't create order")
			return
		}
	} else { // match orders
		for _, offer := range orders {
			budget = float(givenItem.PriceMax) * newOrder.GivenAmount
			got

		}

	}
	// return OrderResponse
	resp, err := json.Marshal(Response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Can't marshall response!")
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)

}
