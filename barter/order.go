package barter

import (
	"bytes"
	"ebarter/barter/db"
	"ebarter/notify"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type OrderResponse struct {
	AcquiredWantedItemAmount int `json:"acquired_wanted_item_amount"`
	SurplusGivenItemAmount   int `json:"surplus_given_item_amount"`
	InorderGivenItemAmount   int `json:"inorder_given_item_amount"`
}

type ReceivedItem struct {
	Name     string  `json:"name"`
	PriceMin float64 `json:"price_min"`
	PriceMax float64 `json:"price_max"`
}

type CalculationResponse struct {
	ItemName string  `json:"item_name"`
	Amount   int     `json:"amount"`
	Cost     float64 `json:"cost"`
}

var (
	// inventories defined by name
	inventories                = []string{"food", "electronics", "construction", "art"}
	invalidOrderFormatResponse = []byte(`{"message": "Invalid order."}`)
	inventoryServicesPrefix    = "http://localhost:8080"
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
	if err != nil {
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
	// check amounts
	if o.GivenAmount < 1 || o.WantedAmount < 1 {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Amounts can't be lower than 1.")
		return
	}
	// ask inventories
	client := &http.Client{}
	jsonBodyWanted := []byte(`{"name": "` + o.WantedItem + `"}`)
	jsonBodyGiven := []byte(`{"name": "` + o.GivenItem + `"}`)
	urlGiven := strings.Join([]string{inventoryServicesPrefix, o.GivenInventory, "check"}, "/")
	urlWanted := strings.Join([]string{inventoryServicesPrefix, o.WantedInventory, "check"}, "/")
	reqGiven, err := http.NewRequest("GET", urlGiven, bytes.NewBuffer(jsonBodyGiven))
	reqWanted, err2 := http.NewRequest("GET", urlWanted, bytes.NewBuffer(jsonBodyWanted))
	if err != nil || err2 != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error while creating requests: %v | %v", err, err2)
	}
	respGiven, err := client.Do(reqGiven)
	respWanted, err2 := client.Do(reqWanted)
	if err != nil || err2 != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error while requesting from inventory: %v | %v", err, err2)
	}
	bodyGiven, err := ioutil.ReadAll(respGiven.Body)
	bodyWanted, err2 := ioutil.ReadAll(respWanted.Body)
	if err != nil || err2 != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error while reading response body: %v | %v", err, err2)
	}
	respGiven.Body.Close()
	respWanted.Body.Close()
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
	// returned Response
	var Response OrderResponse = OrderResponse{
		0,
		0,
		o.GivenAmount, // inorder_item_amount
	}
	// create but don't write newOrder
	newOrder := db.NewOrder(user.Email, user.Key,
		o.GivenInventory, givenItem.Name, o.GivenAmount,
		o.WantedInventory, wantedItem.Name, o.WantedAmount)
	// match order
	orders := db.FindOrders(o.GivenInventory, givenItem.Name,
		o.WantedInventory, wantedItem.Name)
	if orders != nil {
		// match orders
		for _, offer := range orders {
			interval_g1, interval_g2 := givenItem.PriceMin*float64(newOrder.GivenAmount), givenItem.PriceMax*float64(newOrder.GivenAmount)
			interval_w1, interval_w2 := wantedItem.PriceMin*float64(offer.GivenAmount), wantedItem.PriceMax*float64(offer.GivenAmount)
			// check barter
			if (interval_w1 <= interval_g1 && interval_g1 <= interval_w2) ||
				(interval_g1 <= interval_w1 && interval_w1 <= interval_g2) {
				// barter
				Response.AcquiredWantedItemAmount = offer.GivenAmount
				Response.InorderGivenItemAmount = 0
				Response.SurplusGivenItemAmount = Max(0, o.GivenAmount-offer.WantedAmount)
				u2 := db.GetUser(offer.Key)
				log.Println(offer)
				if u2 == nil {
					w.WriteHeader(http.StatusInternalServerError)
					log.Println("Can't find user of the order.")
					return
				}
				var n = notify.Notification{
					U1:            user,
					U2:            u2,
					GotInventory:  newOrder.WantedInventory,
					GotItem:       newOrder.WantedItem,
					GotAmount:     offer.GivenAmount,
					GaveInventory: newOrder.GivenInventory,
					GaveItem:      newOrder.GivenItem,
					GaveAmount:    newOrder.GivenAmount - Response.SurplusGivenItemAmount,
				}
				// notify and del record at the background
				go notify.Notify(n)
				go offer.PermDel()
				break
			}
		}
	}
	if Response.AcquiredWantedItemAmount == 0 {
		// barter didn't happen
		if !newOrder.Create() {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("Can't create order.")
			return
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
