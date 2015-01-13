package main

import (
	"fmt"
	"net/http"
	"time"
	"github.com/gorilla/mux"
	"encoding/json"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

//ROUTING

func main() {

	mainRouter := mux.NewRouter()

	getSubrouter := mainRouter.Methods("GET").Subrouter()
	postSubrouter := mainRouter.Methods("POST").Headers("key", "").Subrouter()
	putSubrouter := mainRouter.Methods("PUT").Headers("key", "").Subrouter()
	deleteSubrouter := mainRouter.Methods("DELETE").Headers("key", "").Subrouter()
	
	http.Handle("/", mainRouter)
	
	getSubrouter.HandleFunc("/api/v1", describe)
	
	//current
	getSubrouter.HandleFunc("/api/v1/current", searchCurrent) //current
	putSubrouter.HandleFunc("/api/v1/current", updateCurrent) //current
	putSubrouter.HandleFunc("/api/v1/current", notAllowed) //current
	deleteSubrouter.HandleFunc("/api/v1/current", notAllowed) //current
	
	//player
	getSubrouter.HandleFunc("/api/v1/player", searchPlayerData) //player
	postSubrouter.HandleFunc("/api/v1/player", registerNewPlayer) //player
	putSubrouter.HandleFunc("/api/v1/player", insertPlayerData) //player
	deleteSubrouter.HandleFunc("/api/v1/player", notImplemented) //player
	
	//history
	getSubrouter.HandleFunc("/api/v1/history", displayHistory) //history
	postSubrouter.HandleFunc("/api/v1/history", insertHistory) //history
	putSubrouter.HandleFunc("/api/v1/history", notAllowed) //history
	deleteSubrouter.HandleFunc("/api/v1/history", notAllowed) //history
	
	//getSubrouter.HandleFunc("/api/v1/search", search)
	//getSubrouter.HandleFunc("/api/v1/searchPlayerData", searchPlayerData)
	//getSubrouter.HandleFunc("/api/v1/registerNewPlayer", registerNewPlayer)
	//getSubrouter.HandleFunc("/api/v1/display", display)

	//postSubrouter.HandleFunc("/api/v1/insert", insert)
	//postSubrouter.HandleFunc("/api/v1/insertPlayerData", insertPlayerData)
	
	postSubrouter.HandleFunc("/api/v1/testauth", testAuth) //test

	//putSubrouter.HandleFunc("/api/v1/update", update)

	err := http.ListenAndServe(":8081", nil)

	if err != nil {
		panic(err)
	}
}

//ACTUAL FUNCTIONALITY

//PLAYER RELATED

func registerNewPlayer(w http.ResponseWriter, r *http.Request) {

	type searchOutput struct {
		Name   string `bson:"n"`
		UserID string `bson:"id"`
		Score  string `bson:"sc"`
		Time   string `bson:"t"`
		Level  string `bson:"l"`
	}

	urlValues := r.URL.Query()
	name := urlValues.Get("name")

	query := bson.M{
		"n": name,
	}

	s, err := mgo.Dial("localhost:27017")

	if err != nil {
		panic(err)
	}

	c := s.DB("game").C("player")

	results := []searchOutput{}
	c.Find(query).Sort("n").All(&results)

	if len(results) != 0 {
		fmt.Printf("Player with username %s already exists", name)
		out, _ := json.MarshalIndent(results, " ", "  ")
		w.Write(out)
	} else {
		id := "foo"
		newPlayer := bson.M{
			"n":  name,
			"id": id,
			"sc": "0",
			"t":  "0",
			"l":  "0",
		}

		err = c.Insert(newPlayer)

		if err != nil {
			panic(err)
		}
		out, _ := json.MarshalIndent(newPlayer, " ", "  ")
		w.Write(out)
	}
	s.Close()

}

func searchPlayerData(w http.ResponseWriter, r *http.Request) {

	type searchOutput struct {
		Name   string `bson:"n"`
		UserID string `bson:"id"`
		Score  string `bson:"sc"`
		Time   string `bson:"t"`
		Level  string `bson:"l"`
	}

	urlValues := r.URL.Query()
	name := urlValues.Get("name")

	query := bson.M{
		"n": name,
	}

	s, err := mgo.Dial("localhost:27017")

	if err != nil {
		panic(err)
	}

	c := s.DB("game").C("player")

	results := []searchOutput{}

	c.Find(query).Sort("n").All(&results)
	out, _ := json.MarshalIndent(results, " ", "  ")
	w.Write(out)
	s.Close()
}

func insertPlayerData(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()

	if err != nil {
		panic(err)
	}

	urlValues := r.Form

	name := urlValues.Get("name")
	userid := urlValues.Get("userid")
	score := urlValues.Get("score")
	time := urlValues.Get("time")
	level := urlValues.Get("level")

	query := bson.M{
		"n":  name,
		"id": userid,
		"sc": score,
		"t":  time,
		"l":  level,
	}

	s, err := mgo.Dial("localhost:27017")

	if err != nil {
		panic(err)
	}

	c := s.DB("game").C("player")

	err = c.Insert(query)

	if err != nil {
		panic(err)
	}

	s.Close()

}

//CURRENT RELATED

func updateCurrent(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()

	if err != nil {
		panic(err)
	}

	urlValues := r.Form

	name := urlValues.Get("name")
	userid := urlValues.Get("userid")
	perc := urlValues.Get("perc")
	fperc := urlValues.Get("fperc")
	color := urlValues.Get("color")
	fcolor := urlValues.Get("fcolor")

	query := bson.M{
		"n": name,
	}

	s, err := mgo.Dial("localhost:27017")

	if err != nil {
		panic(err)
	}

	c := s.DB("game").C("current")

	change := bson.M{"$set": bson.M{"id": userid, "p": perc, "fp": fperc, "c": color, "fc": fcolor}}
	_, err = c.Upsert(query, change)

	if err != nil {
		panic(err)
	}

	s.Close()

}

func searchCurrent(w http.ResponseWriter, r *http.Request) {

	type searchOutput struct {
		Name   string `bson:"n"`
		UserID string `bson:"id"`
		Perc   string `bson:"p"`
		FPerc  string `bson:"fp"`
		Color  string `bson:"c"`
		FColor string `bson:"fc"`
	}

	urlValues := r.URL.Query()
	name := urlValues.Get("name")

	query := bson.M{
		"n": name,
	}

	s, err := mgo.Dial("localhost:27017")

	if err != nil {
		panic(err)
	}

	c := s.DB("game").C("current")

	results := []searchOutput{}
	c.Find(query).Sort("n").All(&results)

	if len(results) != 0 {
		out, _ := json.MarshalIndent(results, " ", "  ")
		w.Write(out)
	} else {
		fmt.Println(len(results))
		results := searchOutput{name, "-1", "-1", "-1", "-1", "-1"}
		out, _ := json.MarshalIndent(results, " ", "  ")
		w.Write(out)
	}

	s.Close()
}

//HISTORY RELATED

func displayHistory(w http.ResponseWriter, r *http.Request) {

	type searchOutput struct {
		Name      string `bson:"n"`
		UserID    string `bson:"id"`
		Perc      string `bson:"p"`
		FPerc     string `bson:"fp"`
		Color     string `bson:"c"`
		FColor    string `bson:"fc"`
		TimeStamp string `bson:"ts"`
	}

	urlValues := r.URL.Query()
	name := urlValues.Get("name")

	query := bson.M{
		"n": name,
	}

	s, err := mgo.Dial("localhost:27017")

	if err != nil {
		panic(err)
	}

	c := s.DB("game").C("history")

	results := []searchOutput{}

	c.Find(query).Sort("ts").All(&results)
	out, _ := json.MarshalIndent(results, " ", "  ")
	w.Write(out)
	s.Close()
}

func insertHistory(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()

	if err != nil {
		panic(err)
	}

	urlValues := r.Form

	name := urlValues.Get("name")
	userid := urlValues.Get("userid")
	perc := urlValues.Get("perc")
	fperc := urlValues.Get("fperc")
	color := urlValues.Get("color")
	fcolor := urlValues.Get("fcolor")

	query := bson.M{
		"n":  name,
		"id": userid,
		"p":  perc,
		"fp": fperc,
		"c":  color,
		"fc": fcolor,
		"ts": time.Now().Format("2006-01-02 15:04:05"),
	}

	s, err := mgo.Dial("localhost:27017")

	if err != nil {
		panic(err)
	}

	c := s.DB("game").C("history")

	err = c.Insert(query)

	if err != nil {
		panic(err)
	}

	s.Close()

}

//DEFAULT BEHAVIOR

func describe(w http.ResponseWriter, r *http.Request) {

	desc := `Mongodb`

	// ========USAGE=========

	//: /api/v1/
	//: /api/v1/
	//: /api/v1/
	//: /api/v1/

	w.Write([]byte(desc))
}


func notAllowed(w http.ResponseWriter, r *http.Request) {
	msg := "Method Not Allowed"

	w.WriteHeader(405)
	w.Write([]byte(msg))
	
}

func notImplemented(w http.ResponseWriter, r *http.Request) {
	msg := "Not Implemented"

	w.WriteHeader(501)
	w.Write([]byte(msg))
	
}

func testAuth(w http.ResponseWriter, r *http.Request){
	
	if authenticate(r.Header) {
		w.Write([]byte("authenticated"))
	} else{
		w.Write([]byte("skata na fas"))
	}
	
}

func authenticate(h http.Header) bool {
	
	type Auth struct {
		ApiKey string `bson:"apiKey"`
	}

	s, err := mgo.Dial("localhost:27017")

	if err != nil {
		panic(err)
	}

	c := s.DB("game").C("authentication")
	
	query := bson.M{
		"apiKey": h.Get("key"),
	}

	results := []Auth{}

 	c.Find(query).All(&results)
	
	if err != nil {
		return false
	}

	if len(results) > 0 {
		return true
	}
	return false
}

//func newPlayer() {
//}
//func newGameData(){
//}