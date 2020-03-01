// <suppliers>
// 	<supplier name="GlassesRUs" age="4"/>
// 	<supplier name="Redstick" age="20"/>
// </suppliers>

// NOTES
// assign -> =
// initialise a new variable and assign -> :=

package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

const (
	dbUser     = "root"
	dbPassword = ""
	dbName     = "webshopdb"
)

// Suppliers xml structure
type Suppliers struct {
	XMLName  xml.Name   `xml:"suppliers"`
	Supplier []Supplier `xml:"supplier"`
}

// Supplier xml structure
type Supplier struct {
	XMLName xml.Name `xml:"supplier"`
	Name    string   `xml:"name,attr"`
	Age     string   `xml:"age,attr"`
}

// ========= API BEGIN

// Stock structure with glasses and frames
type Stock struct {
	Glasses uint8 `json:"glasses"`
	Frames  uint8 `json:"frames"`
}

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

// func handleError(err error) {
// 	if err != nil {
// 		panic(err.Error())
// 	}
// }

func loadSuppliers(w http.ResponseWriter, r *http.Request) {
	var stock Stock
	w.Header().Set("Content-Type", "application/xml")
	// w.WriteHeader(http.StatusOK)
	reqBody, err := ioutil.ReadAll(r.Body)
	handleError(err)

	xml.Unmarshal(reqBody, &stock)
	w.WriteHeader(http.StatusResetContent)
	xml.NewEncoder(w).Encode(stock)

}

// HandleRequests : init http endpoints
func HandleRequests() {
	// fs := http.FileServer(http.Dir("assets/"))
	// http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	http.HandleFunc("/", HomePage)
	http.HandleFunc("/suppliers/load", loadSuppliers)
	http.HandleFunc("/forecast/stock", GetCurrentStock)
	log.Fatal(http.ListenAndServe(":8081", nil))
}

// ========= API END

func readXML(filename string) Suppliers {
	// Open our xmlFile
	xmlFile, err := os.Open(filename)
	// if we os.Open returns an error then handle it
	handleError(err)

	// read our opened xmlFile as a byte array.
	byteValue, err := ioutil.ReadAll(xmlFile)
	handleError(err)

	// we initialize our Supplier array
	var suppliers Suppliers
	// we unmarshal our byteArray which contains our
	// xmlFiles content into 'supplier' which we defined above
	xml.Unmarshal(byteValue, &suppliers)

	// output XML file content
	// for i := 0; i < len(suppliers.Supplier); i++ {
	// 	name := suppliers.Supplier[i].Name
	// 	age := suppliers.Supplier[i].Age
	// 	// row, err := db.Query(fmt.Sprintf("INSERT INTO suppliers (name, age) VALUES ('%s', %d)", name, age))
	// 	fmt.Println(name)
	// 	fmt.Println(age)
	// }

	// defer the closing of our xmlFile so that we can parse it later on
	defer xmlFile.Close()
	return suppliers
}

func getTotalSuppliers() int8 {
	var count int8
	totalSuppliers := db.QueryRow("SELECT COUNT(id) as count FROM suppliers;")
	totalSuppliers.Scan(&count)
	return count
}

func DBstats() {
	totalSuppliers := getTotalSuppliers()
	fmt.Printf("Total Suppliers: %d\n\n", totalSuppliers)
}

func insertIntoDB(suppliers Suppliers) {
	for i := 0; i < len(suppliers.Supplier); i++ {
		age, _ := strconv.Atoi(suppliers.Supplier[i].Age)
		name := suppliers.Supplier[i].Name

		_, err := db.Exec(fmt.Sprintf("INSERT INTO suppliers (name, age) VALUES ('%s', %d) ON DUPLICATE KEY UPDATE age = %d;", name, age, age))
		handleError(err)
	}

	fmt.Printf("Database successfully updated!\n\n")
}

func handleError(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func configDB() {
	var err error
	// create MySQL database.
	dbinfo := fmt.Sprintf("%s@/%s", dbUser, dbName)
	db, err = sql.Open("mysql", dbinfo)

	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dbName))
	handleError(err)

	// select created database
	_, err = db.Exec(fmt.Sprintf("USE %s", dbName))
	handleError(err)

	// create table Suppliers(id, name, age, created_at, updated_at)
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS suppliers (id INTEGER PRIMARY KEY AUTO_INCREMENT, name varchar(255) NOT NULL UNIQUE, age INTEGER NOT NULL, created_at datetime DEFAULT NULL DEFAULT NOW(), updated_at datetime DEFAULT NULL DEFAULT NOW() ON UPDATE NOW())")
	handleError(err)

	// just a sanity check
	_, err = db.Query("SELECT COUNT(*) FROM suppliers;")
	handleError(err)
}

// DBListSuppliers list all suppliers in database
func DBListSuppliers() {
	rows, err := db.Query("SELECT name, age FROM suppliers;")
	handleError(err)
	defer rows.Close()
	var (
		name string
		age  int8
	)
	fmt.Println()
	for rows.Next() {
		err := rows.Scan(&name, &age)
		handleError(err)
		fmt.Println(name, age)
	}
	fmt.Println("\n\n")
}

func showMenu() {
	fmt.Println("1- Load Suppliers")
	fmt.Println("2- List suppliers")
	fmt.Println("3- Age of supplier")
	fmt.Println("4- Exit")
	fmt.Print("Enter option: ")
}

func readUserOption() int {
	var input int
	fmt.Scanln(&input)
	return input
}

func clearScreen() {
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	c.Run()
}

func processSubMenuOption(option int) {
	// read from XML file
	if option == 1 {
		fmt.Println("Insert XML file name: ")
		reader := bufio.NewReader(os.Stdin)
		xmlFileName, _ := reader.ReadString('\n')
		xmlData := readXML(xmlFileName)
		insertIntoDB(xmlData)
	} else if option == 2 {
		// read from API
		HandleRequests()
	}
}

func processMenuOption(option int) {
	// clear screen
	clearScreen()
	switch option {
	case 1:
		subOption := 0
		for ok := true; ok; ok = (subOption != 3) {
			clearScreen()
			fmt.Println("1- From file")
			fmt.Println("2- From API")
			subOption = readUserOption()
			processSubMenuOption(subOption)
		}
	case 2:
		fmt.Println("List Suppliers")
		DBListSuppliers()
	case 3:
		fmt.Println("two")
	case 4:
		defer db.Close()
		os.Exit(42)
	}
}

func main() {

	// if len(os.Args) < 2 {
	// 	fmt.Println("Missing XML file")
	// } else {
	configDB()
	// xmlData := readXML(os.Args[1])

	// insertIntoDB(xmlData)
	DBstats()

	option := 0
	for ok := true; ok; ok = (option != 3) {
		showMenu()
		option = readUserOption()
		processMenuOption(option)
	}
	// }
	defer db.Close()
}
