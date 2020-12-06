package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/rs/cors"
	"gopkg.in/mgo.v2"

	"gopkg.in/mgo.v2/bson"
)

type Result struct {
	ID        string `json:"id"`
	Referee   string `json:"referee"`
	Score     int32  `json:"score"`
	Sportsman string `json:"sportsman"`
	TypeScore string `json:"typeScore"`
}

type SumResult struct {
	BirthYear       int32   `json:"birthYear"`
	NameSportsmans  string  `json:"name"`
	CurrentCategory string  `json:"grade"`
	DesiredCategory string  `json:"gradeToClaim"`
	Coaches         string  `json:"coaches"`
	FirstType       float32 `json"firstType"`
	SecondType      float32 `json"secondType"`
	ThirdTType      float32 `json"thirdTType"`
	FourthType      float32 `json"fourthType"`
	SumOfPoint      float32 `json"sumOfPoint"`
}

type finalResultsSportsmans struct {
	Id              bson.ObjectId `bson:"_id"`
	BirthYear       int32         `bson:"birthYear"`
	NameSportsmans  string        `bson:"name"`
	CurrentCategory string        `bson:"grade"`
	DesiredCategory string        `bson:"gradeToClaim"`
	Coaches         string        `bson:"coaches"`
	FirstType       float32       `bson"firstType"`
	SecondType      float32       `bson"secondType"`
	ThirdTType      float32       `bson"thirdTType"`
	FourthType      float32       `bson"fourthType"`
	SumOfPoint      float32       `bson"fourthType"`
}

type UserDB struct {
	Id       bson.ObjectId `bson:"_id"`
	Login    string        `bson:"login"`
	Password string        `bson:"password"`
	Owner    string        `bson:"owner"`
}

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Owner    string `json:"owner"`
}

type Sportsman struct {
	Id              bson.ObjectId `json:"_id"`
	BirthYear       int32         `json:"birthYear"`
	NameSportsmans  string        `json:"name"`
	CurrentCategory string        `json:"grade"`
	DesiredCategory string        `json:"gradeToClaim"`
	Coaches         string        `json:"coaches"`
}

type SportsmanDB struct {
	Id              bson.ObjectId `bson:"_id"`
	BirthYear       int32         `bson:birthYear`
	Name            string        `bson:name`
	CurrentCategory string        `bson:"grade"`
	DesiredCategory string        `bson:"gradeToClaim"`
	Coaches         string        `bson:"coaches"`
}

var currenctSportsman Sportsman
var sumResult SumResult
var results []Result
var numberCurrentSportsman int32 = 0

//API
func finallyResults(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var sumResult SumResult
	_ = json.NewDecoder(r.Body).Decode(&sumResult)
	readResult(sumResult)
}
func getResults(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
func createUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user User
	_ = json.NewDecoder(r.Body).Decode(&user)
	createUserDB(user)
}

func authUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//params := mux.Vars(r)
	var user User
	_ = json.NewDecoder(r.Body).Decode(&user)
	checkAuthUser(user)
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
func createSportsman(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var sportsman Sportsman
	_ = json.NewDecoder(r.Body).Decode(&sportsman)
	createSportsmanDB(sportsman)
	json.NewEncoder(w).Encode(sportsman)
	//json.NewEncoder(w).Encode(result)
}

func getSportsmans(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	result := getSportsmanDB()
	json.NewEncoder(w).Encode(result)
}
func getCurrentSportsman(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	result := getSportsmanDB()
	currenctSportsman = result[numberCurrentSportsman]
	json.NewEncoder(w).Encode(result[numberCurrentSportsman])
	numberCurrentSportsman++
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
	var sumDFirstWave float32
	var sumDSecondWave float32

	var sumEFirstWave float32
	var sumESecondWave float32

	//params := mux.Vars(r)
	for _, item := range results {
		if item.TypeScore == "d1" || item.TypeScore == "d2" {
			sumDFirstWave += float32(item.Score)

		}
		if item.TypeScore == "d3" || item.TypeScore == "d4" {
			sumDSecondWave += float32(item.Score)
		}
		if item.TypeScore == "e1" || item.TypeScore == "e2" {
			sumEFirstWave += float32(item.Score)
		} else {
			sumESecondWave += float32(item.Score)
		}

	}

	sumResult.SumOfPoint = sumDFirstWave/2 + sumDSecondWave/2 + sumEFirstWave/2 + sumESecondWave/2
	sumResult.FirstType = sumDFirstWave / 2
	sumResult.SecondType = sumDSecondWave / 2
	sumResult.ThirdTType = sumEFirstWave / 2
	sumResult.FourthType = sumESecondWave / 2
	sumResult.NameSportsmans = results[0].Sportsman

	json.NewEncoder(w).Encode(&SumResult{})
	writeResult(sumResult)
}
func getUserOwner(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	var user User
	_ = json.NewDecoder(r.Body).Decode(&user)
	result := userOwnerDB(user)
	json.NewEncoder(w).Encode(result)

}

//DATA BASE
func getSportsmanDB() []Sportsman {
	// открываем соединение
	var sportsmansJson []Sportsman
	session, err := mgo.Dial("mongodb://127.0.0.1")
	if err != nil {
		panic(err)

	}
	defer session.Close()

	// получаем коллекцию
	userCollection := session.DB("sportsmanDatadb").C("sportsmans")
	// критерий выборки

	query := bson.M{}
	// объект для сохранения результата
	sportsmans := []SportsmanDB{}
	userCollection.Find(query).All(&sportsmans)

	for _, u := range sportsmans {

		sportsmansJson = append(sportsmansJson, Sportsman{Id: bson.NewObjectId(),
			NameSportsmans:  u.Name,
			CurrentCategory: u.CurrentCategory,
			DesiredCategory: u.DesiredCategory,
			BirthYear:       u.BirthYear,
			Coaches:         u.Coaches})
	}

	session.Close()
	return sportsmansJson
}
func createSportsmanDB(sportsman Sportsman) {
	// открываем соединение
	session, err := mgo.Dial("mongodb://127.0.0.1")
	if err != nil {
		panic(err)

	}
	defer session.Close()

	// получаем коллекцию
	userCollection := session.DB("sportsmanDatadb").C("sportsmans")

	if err != nil {
		fmt.Println(err)
	}

	u := &SportsmanDB{Id: bson.NewObjectId(),
		Name:            sportsman.NameSportsmans,
		CurrentCategory: sportsman.CurrentCategory,
		DesiredCategory: sportsman.DesiredCategory,
		BirthYear:       sportsman.BirthYear,
		Coaches:         sportsman.Coaches}
	// добавляем объект
	err = userCollection.Insert(u)
	if err != nil {
		fmt.Println(err)
	}
	session.Close()
}
func readResult(sumResult SumResult) []SumResult {
	// открываем соединение
	var sumResultJson []SumResult
	session, err := mgo.Dial("mongodb://127.0.0.1")
	if err != nil {
		panic(err)

	}
	defer session.Close()

	// получаем коллекцию
	userCollection := session.DB("sportsmanDatadb").C("resultsEvent")
	// критерий выборки
	query := bson.M{}
	// объект для сохранения результата
	results := []finalResultsSportsmans{}

	userCollection.Find(query).Sort("sumOfPoint").All(&results)

	for _, u := range results {

		sumResultJson = append(sumResultJson, SumResult{
			BirthYear:       u.BirthYear,
			NameSportsmans:  u.NameSportsmans,
			CurrentCategory: u.CurrentCategory,
			DesiredCategory: u.DesiredCategory,
			Coaches:         u.Coaches,
			FirstType:       u.FirstType,
			SecondType:      u.SecondType,
			ThirdTType:      u.ThirdTType,
			FourthType:      u.FourthType,
			SumOfPoint:      u.SumOfPoint})
	}

	session.Close()
	return sumResultJson
}
func writeResult(sumResult SumResult) {
	// открываем соединение
	session, err := mgo.Dial("mongodb://127.0.0.1")
	if err != nil {
		panic(err)

	}
	defer session.Close()

	// получаем коллекцию
	userCollection := session.DB("sportsmanDatadb").C("resultsEvent")

	if err != nil {
		fmt.Println(err)
	}

	u := &finalResultsSportsmans{Id: bson.NewObjectId(),
		BirthYear:       currenctSportsman.BirthYear,
		NameSportsmans:  sumResult.NameSportsmans,
		CurrentCategory: currenctSportsman.CurrentCategory,
		DesiredCategory: currenctSportsman.DesiredCategory,
		Coaches:         currenctSportsman.Coaches,
		FirstType:       sumResult.FirstType,
		SecondType:      sumResult.SecondType,
		ThirdTType:      sumResult.ThirdTType,
		FourthType:      sumResult.FourthType,
		SumOfPoint:      sumResult.SumOfPoint,
	}
	// добавляем объект
	err = userCollection.Insert(u)
	if err != nil {
		fmt.Println(err)
	}
	session.Close()
}

func userOwnerDB(user User) []User {
	// открываем соединение
	var usersJson []User
	session, err := mgo.Dial("mongodb://127.0.0.1")
	if err != nil {
		panic(err)

	}
	defer session.Close()

	// получаем коллекцию
	userCollection := session.DB("sportsmanDatadb").C("users")
	// критерий выборки
	fmt.Println(user.Login)
	query := bson.M{
		"owner": bson.M{
			"$eq": user.Owner,
		},
	}
	// объект для сохранения результата
	users := []UserDB{}

	userCollection.Find(query).All(&users)

	for _, u := range users {

		usersJson = append(usersJson, User{Login: u.Login, Password: u.Password, Owner: u.Owner})
	}

	session.Close()
	return usersJson
}
func createUserDB(user User) {
	// открываем соединение
	session, err := mgo.Dial("mongodb://127.0.0.1")
	if err != nil {
		panic(err)

	}
	defer session.Close()

	// получаем коллекцию
	userCollection := session.DB("sportsmanDatadb").C("users")

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(user.Login)
	u := &UserDB{Id: bson.NewObjectId(), Login: user.Login, Password: user.Password, Owner: user.Owner}
	// добавляем объект
	err = userCollection.Insert(u)
	if err != nil {
		fmt.Println(err)
	}
	session.Close()
}

func checkAuthUser(user User) bool {
	// открываем соединение
	// открываем соединение
	fmt.Println("nani")
	session, err := mgo.Dial("mongodb://127.0.0.1")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// получаем коллекцию
	userCollection := session.DB("sportsmanDatadb").C("users")
	// критерий выборки
	fmt.Println(user.Login)
	query := bson.M{
		"login": bson.M{
			"$eq": user.Login,
		},
		"password": bson.M{
			"$eq": user.Password,
		},
	}
	// объект для сохранения результата
	users := []UserDB{}
	userCollection.Find(query).All(&users)

	if len(users) == 0 {
		return false
	}

	session.Close()
	return true
}

func main() {
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"}, // All origins

		AllowedMethods: []string{"POST"}, // Allowing only get, just an example
	})
	r := mux.NewRouter()

	results = append(results, Result{ID: "2", Referee: "Анатолий Сумаилов Тигранович", Score: 9, Sportsman: "Дарья Кунцевич", TypeScore: "d1"})
	r.HandleFunc("/auth", authUser).Methods("GET")
	r.HandleFunc("/registration", createUser).Methods("POST")
	r.HandleFunc("/usersOwner", getUserOwner).Methods("GET")
	r.HandleFunc("/sportsmans", getSportsmans).Methods("GET")
	r.HandleFunc("/results/currentSportsmans", getCurrentSportsman).Methods("GET")
	r.HandleFunc("/sportsmans/createSportsmans", createSportsman).Methods("POST")
	r.HandleFunc("/results", getResults).Methods("GET")
	r.HandleFunc("/results/finallyResults", finallyResults).Methods("GET")
	r.HandleFunc("/results/{sportsman}", getResult).Methods("GET")

	if len(results) < 10 {
		r.HandleFunc("/results/newResult", createResult).Methods("POST")
	}
	if len(results) == 10 {
		r.HandleFunc("/results/finalResult", sumarryResult).Methods("GET")
		results = nil

	}
	log.Fatal(http.ListenAndServe(":8000", c.Handler(r)))
}
