package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Stock structure with glasses and frames
type Stock struct {
	Glasses uint8 `json:"glasses"`
	Frames  uint8 `json:"frames"`
}

// Supplier structure to define supplier
// type Supplier struct {
// 	Name string `json:"name"`
// 	Age  uint8  `json:"age"`
// }

// type Suppliers []Supplier

// GetCurrentStock /forecast/stock/T
func GetCurrentStock(w http.ResponseWriter, r *http.Request) {
	stock := Stock{Glasses: 10, Frames: 99}
	prettyJSON, err := json.MarshalIndent(stock, "", "    ")
	// json.NewEncoder(w).Encode(stock
	if err != nil {
		log.Fatal("Failed to generate json", err)
	}

	fmt.Fprintf(w, string(prettyJSON))
}

// HomePage default page
func HomePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to GO")
}

func handleRequests() {
	http.HandleFunc("/", HomePage)
	http.HandleFunc("/forecast/stock", GetCurrentStock)
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func main() {
	handleRequests()
}
