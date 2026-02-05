package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"time"
)

type CalculateFeeRequest struct {
	CartValue        int    `json:"cart_value"`        // in cents
	DeliveryDistance int    `json:"delivery_distance"` // in meters
	NumberOfItems    int    `json:"number_of_items"`
	Time             string `json:"time"` // ISO 8601 format
}

type CalculateFeeResponse struct {
	DeliveryFee int `json:"delivery_fee"` // in cents
}

func main() {
	http.HandleFunc("POST /calculate-delivery-fee", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	var order CalculateFeeRequest

	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println(order)
	fee := CalculateDeliveryFee(&order)
	resp := CalculateFeeResponse{DeliveryFee: fee}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)

}

func CalculateDeliveryFee(order *CalculateFeeRequest) int {

	// The delivery is free (0€) when the cart value is equal or more than 200€
	if order.CartValue >= 20000 {
		return 0
	}

	fee := 0

	// If the cart value is less than 10€, a small order surcharge is added to the delivery price
	if order.CartValue < 1000 {
		fee += 1000 - order.CartValue
	}

	// The delivery fee is 2€ for the first 1000 meters. For every additional 500 meters, 1€ is added to the delivery fee
	if order.DeliveryDistance <= 500 {
		fee += 100
	} else if order.DeliveryDistance <= 1000 {
		fee += 200
	} else {
		fee += 200 + int(math.Ceil(float64(order.DeliveryDistance-1000)/500))*100
	}

	// If the number of items is 5 or more, an additional surcharge of 0.50€ per item is added for each item above 4
	if order.NumberOfItems >= 5 {
		fee += (order.NumberOfItems - 4) * 50
	}

	// If the number of items is 12 or more, an additional bulk fee of 1.20€ is added to the delivery fee
	if order.NumberOfItems >= 12 {
		fee += 120
	}

	// During the Friday rush (between 3 PM and 7 PM on Fridays), the delivery fee (including possible surcharges) will be multiplied by 1.2x
	t, err := time.Parse(time.RFC3339, order.Time)
	if err != nil {
		log.Printf("Error parsing time: %v", err)
		return 0
	}

	if t.Weekday() == time.Friday && t.Hour() >= 15 && t.Hour() < 19 {
		fee = int(float64(fee) * 1.2)
	}

	// The delivery fee can never be more than 15€, including possible surcharges.
	if fee > 1500 {
		fee = 1500
	}

	return fee
}
