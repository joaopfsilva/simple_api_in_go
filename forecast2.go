// package main

// import (
// 	"encoding/json"
// 	"encoding/xml"
// 	"fmt"
// 	"io/ioutil"
// 	"log"
// 	"net/http"
// )

// // Stock structure with glasses and frames
// type Stock struct {
// 	Glasses uint8 `json:"glasses"`
// 	Frames  uint8 `json:"frames"`
// }

// // GetCurrentStock /forecast/stock/T
// func GetCurrentStock(w http.ResponseWriter, r *http.Request) {
// 	stock := Stock{Glasses: 10, Frames: 99}
// 	prettyJSON, err := json.MarshalIndent(stock, "", "    ")
// 	// json.NewEncoder(w).Encode(stock
// 	if err != nil {
// 		log.Fatal("Failed to generate json", err)
// 	}

// 	fmt.Fprintf(w, string(prettyJSON))
// }

// // HomePage default page
// func HomePage(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, "Welcome to GO")
// }

// func loadSuppliers(w http.ResponseWriter, r *http.Request) {
// 	var stock Stock
// 	w.Header().Set("Content-Type", "application/xml")
// 	// w.WriteHeader(http.StatusOK)
// 	reqBody, err := ioutil.ReadAll(r.Body)
// 	handleError(err)

// 	xml.Unmarshal(reqBody, &stock)
// 	w.WriteHeader(http.StatusResetContent)
// 	xml.NewEncoder(w).Encode(stock)

// }

// func handleRequests() {

// 	http.HandleFunc("/", HomePage)
// 	http.HandleFunc("/suppliers/load", loadSuppliers)
// 	http.HandleFunc("/forecast/stock", GetCurrentStock)
// 	log.Fatal(http.ListenAndServe(":8081", nil))
// }

// func main() {
// 	handleRequests()
// }
