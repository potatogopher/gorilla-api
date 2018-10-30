package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"

	// Empty import because driver
	_ "github.com/lib/pq"
)

const (
	port = ":8081"

	dbPort   = 5432
	host     = "localhost"
	user     = "nicholas.rucci"
	password = ""
	dbname   = "test"
)

// DB is the psql db instance
var DB *sqlx.DB

func init() {
	psqlInfo := fmt.Sprintf("host=%s port=%d  dbname=%s sslmode=disable", host, dbPort, dbname)

	db, err := sqlx.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Println("error connecting to db")
	}

	DB = db
}

// HealthzHandler ...
func HealthzHandler(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("ok"))
	if err != nil {
		log.Fatalln("An error occured")
	}
}

// Person is the person type
type Person struct {
	Name string
}

// GetPeopleHandler ...
func GetPeopleHandler(w http.ResponseWriter, r *http.Request) {
	var people []Person
	err := DB.QueryRowx("SELECT * FROM people").StructScan(&people)
	if err != nil {
		log.Fatalln("An error occured when querying for people: ", err)
	}

	b, err := json.Marshal(&people)
	if err != nil {
		log.Fatalln("An error occured when marshaling data")
	}

	_, err = w.Write(b)
	if err != nil {
		log.Fatalln("An error occured")
	}
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/healthz", HealthzHandler)
	router.HandleFunc("/people", GetPeopleHandler)
	err := http.ListenAndServe(port, router)
	if err != nil {
		log.Fatalln("An error occured when trying to start the server")
	}
}
