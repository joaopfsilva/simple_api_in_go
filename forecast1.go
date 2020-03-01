// <suppliers>
// 	<supplier name="GlassesRUs" age="4"/>
// 	<supplier name="Redstick" age="20"/>
// </suppliers>

// NOTES
// assign -> =
// initialise a new variable and assign -> :=

package main

import (
	"database/sql"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"

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
	Age     int8     `xml:"age,attr"`
}

func readXML(filename string, db *sql.DB) Suppliers {
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

func checkCount(rows *sql.Rows) (count int) {
	rows.Scan(&count)
	return count
}

func DBstats() {
	res, _ := db.Query("SELECT COUNT(*) FROM suppliers;")
	totalSuppliers := checkCount(res)
	fmt.Printf("Total Suppliers: %d\n\n", totalSuppliers)
}

func insertIntoDB(suppliers Suppliers) {
	var errors int = 0
	for i := 0; i < len(suppliers.Supplier); i++ {
		_, err := db.Exec("INSERT INTO suppliers (name, age) VALUES (?, ?)", suppliers.Supplier[i].Name, suppliers.Supplier[i].Age)
		if err != nil {
			errors++
		}
	}
	fmt.Printf("Successfully inserted values %d of %d\n", len(suppliers.Supplier)-errors, len(suppliers.Supplier))
}

func handleError(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func configDB() {
	var err error
	// create MySQL database.
	dbinfo := fmt.Sprintf("%s@/", dbUser)
	db, err = sql.Open("mysql", dbinfo)

	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dbName))
	handleError(err)

	_, err = db.Exec(fmt.Sprintf("USE %s", dbName))
	handleError(err)

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS suppliers (id INTEGER PRIMARY KEY AUTO_INCREMENT, name varchar(255) NOT NULL UNIQUE, age INTEGER NOT NULL, created_at datetime DEFAULT NULL DEFAULT NOW(), updated_at datetime DEFAULT NULL DEFAULT NOW() ON UPDATE NOW())")
	handleError(err)

	// just a check
	_, err = db.Exec("SELECT COUNT(*) FROM suppliers;")
	handleError(err)
}

func showMenu() {
	var input string

	fmt.Println("1- List suppliers")
	fmt.Println("2- Age of supplier")
	fmt.Print("Enter option: ")

	fmt.Scanln(&input)
	fmt.Print(input)
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Missing XML file")
	} else {
		configDB()
		xmlData := readXML(os.Args[1], db)
		insertIntoDB(xmlData)
		DBstats()
		showMenu()
	}

	defer db.Close()
	os.Exit(42)
}
