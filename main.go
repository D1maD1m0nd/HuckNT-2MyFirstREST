package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// "encoding/json"
// "log"
// "math/rand"
// "net/http"
// "strconv"

// "github.com/gorilla/mux"

type Result struct {
	ID        string `json:"id"`
	Referee   string `json:"referee"`
	Score     int32  `json:"score"`
	Sportsman string `json:"sportsman"`
	TypeScore string `json:"typeScore"`
}

var results []Result

//API
func getResults(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func getResult(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range results {
		if item.Sportsman == params["sportsman"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Result{})
}
func createResult(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var result Result
	_ = json.NewDecoder(r.Body).Decode(&result)
	result.ID = strconv.Itoa(rand.Intn(1000000))
	results = append(results, result)
	json.NewEncoder(w).Encode(result)
}

// работа с БД
func pushDataDB() {
	connStr := "user=postgres password=mypass dbname=productdb sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	result, err := db.Exec("insert into Products (model, company, price) values ('iPhone X', $1, $2)",
		"Apple", 72000)
	if err != nil {
		panic(err)
	}
	fmt.Println(result.LastInsertId()) // не поддерживается
	fmt.Println(result.RowsAffected()) // количество добавленных строк
}

func getDataDB() {
	// connStr := "user=postgres password=mypass dbname=productdb sslmode=disable"
	// db, err := sql.Open("postgres", connStr)
	// if err != nil {
	// 	panic(err)
	// }
	// defer db.Close()

	// rows, err := db.Query("select * from Products")
	// if err != nil {
	// 	panic(err)
	// }
	// defer rows.Close()
	// products := []product{}

	// for rows.Next() {
	// 	p := product{}
	// 	err := rows.Scan(&p.id, &p.model, &p.company, &p.price)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		continue
	// 	}
	// 	products = append(products, p)
	// }
	// for _, p := range products {
	// 	fmt.Println(p.id, p.model, p.company, p.price)
	// }
}
func main() {
	r := mux.NewRouter()

	results = append(results, Result{ID: "2", Referee: "Анатолий Сумаилов Тигранович", Score: 9, Sportsman: "Дарья Кунцевич", TypeScore: "d1"})
	r.HandleFunc("/results", getResults).Methods("GET")
	r.HandleFunc("/results/{id}", getResult).Methods("GET")
	r.HandleFunc("/results", createResult).Methods("POST")
	log.Fatal(http.ListenAndServe(":8000", r))
}
