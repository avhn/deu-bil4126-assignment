package barter

import (
	"ebarter/barter/db"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

// matching order with inventories

type OrderResponse struct {
	acquired_wanted_item_amount int `json:"acquired_wanted_item_amount"`
	surplus_given_item_amount   int `json:"surplus_given_item_amount"`
	inorder_given_item_amount   int `json:"inorder_given_item_amount"`
}

var (
	// inventories defined by name
	inventories                = []string{"food", "electronics", "construction", "art"}
	invalidOrderFormatResponse = []byte(`{"message": "Invalid order."}`)
)

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
	// ask inventories for prices
	// match the order
	// return OrderResponse

}
