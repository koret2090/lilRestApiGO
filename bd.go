package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

type Person struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Nickname string `json:"nickname"`
}

const (
	user     = "postgres"
	password = "1234"
	dbname   = "persons"
)

func main() {
	fmt.Println("Server is started")

	http.HandleFunc("/", GetAllPersons)
	http.HandleFunc("/insert", AddPerson)
	http.HandleFunc("/update", UpdatePerson)
	http.HandleFunc("/delete", DeletePerson)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func OpenConnection() *sql.DB {
	connStr := fmt.Sprintf("user=%s password=%s dbname =%s "+
		"sslmode=disable", user, password, dbname)

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}

func GetAllPersons(w http.ResponseWriter, r *http.Request) {
	db := OpenConnection()
	defer db.Close()

	rows, err := db.Query("select * from person")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var persons []Person

	for rows.Next() {
		var person Person
		rows.Scan(&person.Id, &person.Name, &person.Nickname)
		persons = append(persons, person)
	}
	personsBytes, err := json.MarshalIndent(persons, "", "\t")

	w.Header().Set("Content-Type", "application/json")
	w.Write(personsBytes)

}

func AddPerson(w http.ResponseWriter, r *http.Request) {
	db := OpenConnection()

	var p Person
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sqlStatement := "insert into person (id, name, nickname) values ($1, $2, $3)"
	_, err = db.Exec(sqlStatement, p.Id, p.Name, p.Nickname)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	defer db.Close()
}

func UpdatePerson(w http.ResponseWriter, r *http.Request) {
	db := OpenConnection()

	var p Person
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sqlStatement := "update person set name = $2, nickname = $3 where id = $1"
	_, err = db.Exec(sqlStatement, p.Id, p.Name, p.Nickname)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	defer db.Close()
}

func DeletePerson(w http.ResponseWriter, r *http.Request) {
	db := OpenConnection()

	var p Person
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sqlStatement := "delete from person where id = $1"
	_, err = db.Exec(sqlStatement, p.Id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	defer db.Close()
}
