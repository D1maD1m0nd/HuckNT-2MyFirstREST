package main

import (
	"encoding/json"
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

type SumResult struct {
	Sportsman  string `json:"sportsman"`
	FinalScore int32  `json:"finalScore"`
}

type Sportsman struct {
	Name string `json:"name`
}

var sumResult SumResult
var sportsman Sportsman
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
func getCurrentSportsman(w http.ResponseWriter, r *http.Request) {

}
func createResult(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var result Result
	_ = json.NewDecoder(r.Body).Decode(&result)
	result.ID = strconv.Itoa(rand.Intn(1000000))
	results = append(results, result)
	json.NewEncoder(w).Encode(result)
}
func sumarryResult(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var sumDFirstWave int32
	var sumDSecondWave int32

	var sumEFirstWave int32
	var sumESecondWave int32

	var finalResult int32
	//params := mux.Vars(r)
	for _, item := range results {
		if item.TypeScore == "d1" || item.TypeScore == "d2" {
			sumDFirstWave += item.Score

		}
		if item.TypeScore == "d3" || item.TypeScore == "d4" {
			sumDSecondWave += item.Score
		}
		if item.TypeScore == "e1" || item.TypeScore == "e2" {
			sumEFirstWave += item.Score
		} else {
			sumESecondWave += item.Score
		}

		finalResult = (sumDFirstWave/2 + sumDSecondWave/2) + (sumEFirstWave/2 + sumESecondWave/2)
	}
	sumResult.FinalScore = finalResult
	sumResult.Sportsman = results[0].Sportsman
	json.NewEncoder(w).Encode(&SumResult{})
}

func main() {
	r := mux.NewRouter()

	results = append(results, Result{ID: "2", Referee: "Анатолий Сумаилов Тигранович", Score: 9, Sportsman: "Дарья Кунцевич", TypeScore: "d1"})
	r.HandleFunc("/results", getResults).Methods("GET")
	r.HandleFunc("/results/sportsman", getCurrentSportsman).Methods("GET")
	r.HandleFunc("/results/{sportsman}", getResult).Methods("GET")
	if len(results) <= 10 {
		r.HandleFunc("/results", createResult).Methods("POST")
	}
	if len(results) == 10 {
		r.HandleFunc("/results/finalResult", sumarryResult).Methods("GET")

	}

	log.Fatal(http.ListenAndServe(":8000", r))
}
