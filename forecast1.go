// <suppliers>
// 	<supplier name="GlassesRUs" age="4"/>
// 	<supplier name="Redstick" age="20"/>
// </suppliers>

package main

import (
	"database/sql"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
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

func processXML(filename string, db *sql.DB) string {
	// Open our xmlFile
	xmlFile, err := os.Open(filename)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}

	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(xmlFile)

	// we initialize our Supplier array
	var suppliers Suppliers
	// we unmarshal our byteArray which contains our
	// xmlFiles content into 'supplier' which we defined above
	xml.Unmarshal(byteValue, &suppliers)

	// output XML file content
	for i := 0; i < len(suppliers.Supplier); i++ {
		name := suppliers.Supplier[i].Name
		age, _ := strconv.Atoi(suppliers.Supplier[i].Age)
		row, err := db.Query(fmt.Sprintf("INSERT INTO suppliers (name, age) VALUES ('%s', %d)", name, age))
		fmt.Println(err)
		fmt.Println(row)
	}

	// defer the closing of our xmlFile so that we can parse it later on
	defer xmlFile.Close()
	return fmt.Sprintf("Successfully Opened %s", filename)
}

func main() {
	// create MySQL database.
	db, err := sql.Open("mysql", "root@/")
	db.Exec("CREATE DATABASE IF NOT EXISTS webshopdb")
	db.Exec("USE webshopdb")
	if err != nil {
		panic(err.Error())
	}

	db.Exec("CREATE TABLE IF NOT EXISTS suppliers (id INTEGER PRIMARY KEY AUTO INCREMENR, name TEXT, age INTEGER)")
	db.Exec("SELECT COUNT(*) FROM suppliers;")

	if len(os.Args) == 1 {
		fmt.Println("No argument to process")
	} else {
		fmt.Println(processXML(os.Args[1], db))
	}

	defer db.Close()
	os.Exit(42)
}
