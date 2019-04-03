package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var db *sql.DB

var server = "lugotcash.database.windows.net"
var port = 1433
var user = "XXXXXX"
var password = "XXXXXX"
var database = "lusCars"

type Aprticle struct {
	Title   string "json:'title'"
	Desc    string "json:'desc'"
	Content string "json:'content'"
}

type Car struct {
	CarNo   string "json:'carNo'"
	Year    int    "json:'year'"
	Make    string "json:'make'"
	Model   string "json:'model'"
	Milage  int    "json:'milage'"
	Owners  int    "json:'owners'"
	Acd_rpt int    "json:'accidentsReported'"
	Price   string "json:'price'"
}

//var db *sql.DB

//type Articles []Article

func allArticle(w http.ResponseWriter, r *http.Request) {
	articles := Aprticle{"Test title", "Test discription", "Hello world"}
	json.NewEncoder(w).Encode(articles)
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Homepage Endpoint cliped")
}
func getCars(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(getCarsA(db, params["make"], params["model"], params["year"]))
}

func postCar(w http.ResponseWriter, r *http.Request) {
	defer db.Close()
	//POST(db, params["make"], params["model"], params["year"])
}

func getCarsA(db *sql.DB, make string, model string, year string) []Car {

	var Query string

	//defer db.Close()

	if make == "all" {
		if year == "all" {
			//get all cars
			Query = "SELECT * FROM cars"
		} else {
			//get all with specific year
			Query = "SELECT * FROM cars where years = " + year
		}
	} else {
		if model == "all" {
			if year == "all" {
				//get all model for make
				Query = "SELECT * FROM cars where make = '" + make + "'"
			} else {
				//get all models with specific year
				Query = "SELECT * FROM cars where make = '" + make + "' AND years = " + year
			}
		} else {
			if year == "all" {
				//get specific model with all years
				Query = "SELECT * FROM cars where make = '" + make + "' AND model = '" + model + "'"
			} else {
				//get specific car
				Query = "SELECT * FROM cars where make = '" + make + "' AND model = '" + model + "' AND years = '" + year + "'"
			}
		}
	}
	results, err2 := db.Query(Query)
	cars := []Car{}
	for results.Next() {
		var car Car
		err2 = results.Scan(&car.CarNo, &car.Year, &car.Make, &car.Model, &car.Milage, &car.Owners, &car.Acd_rpt, &car.Price)
		if err2 != nil {
			panic(err2.Error())
		}
		cars = append(cars, car)
	}
	return cars
}

func handleRequest() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/cars/{make}/{model}/{year}", getCars).Methods("GET")
	//myRouter.HandleFunc("/cars/{make}/{model}/{year}", getCars).Methods("POST")
	fmt.Println("Running api.....")
	log.Fatal(http.ListenAndServe(":8081", myRouter))

}

func post(db *sql.DB, c Car) {
	iQuery := fmt.Sprintf("%s%s%s%d%s%d%s%d%s%d%s", "INSERT INTO cars(carNo, years, make, model, milage, owners, `accident report`, price)"+
		"VALUES (", c.CarNo, ",'", c.Year, "','"+c.Make+"','"+c.Model+"',", c.Milage, ",", c.Owners, ",", c.Acd_rpt, ",'"+c.Price+"')")
	insert, err2 := db.Query(iQuery)
	if err2 != nil {
		panic(err2.Error())
	}
	defer insert.Close()
}
func connect() {
	// Build connection string
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		server, user, password, port, database)

	var err error

	// Create connection pool
	db, err = sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatal("Error creating connection pool: ", err.Error())
	}
	ctx := context.Background()
	err = db.PingContext(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("Connected!\n")

}

func main() {

	connect()
	defer db.Close()
	handleRequest()
}
