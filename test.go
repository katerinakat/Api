package main

import (
	//"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"encoding/json"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"time"
)

func main() {

	mainRouter := mux.NewRouter()

	getSubrouter := mainRouter.Methods("GET").Subrouter()
	postSubrouter := mainRouter.Methods("POST").Subrouter()
	putSubrouter := mainRouter.Methods("PUT").Subrouter()
	
	http.Handle("/", mainRouter)
	
	getSubrouter.HandleFunc("/api/v1", describe)
	
	getSubrouter.HandleFunc("/api/v1/search", search)
	getSubrouter.HandleFunc("/api/v1/searchGameData", searchGameData)
	getSubrouter.HandleFunc("/api/v1/display", display)

	postSubrouter.HandleFunc("/api/v1/insert", insert)
	postSubrouter.HandleFunc("/api/v1/insertGameData", insertGameData)

	putSubrouter.HandleFunc("/api/v1/update", update)

	err := http.ListenAndServe(":8081", nil)

	if err != nil {
		panic(err)
	}
}

func describe(w http.ResponseWriter, r *http.Request) {

 	desc := `Mongodb`
   
   // ========USAGE=========

	//: /api/v1/
	//: /api/v1/
	//: /api/v1/
	//: /api/v1/

 	w.Write([]byte(desc))
 }

func search(w http.ResponseWriter, r *http.Request) {
	
	type searchOutput struct {	
		Name      string        `bson:"n"`
		UserID    string        `bson:"id"`
		Perc	  string        `bson:"p"`
		FPerc	  string        `bson:"fp"`
		Color	  string        `bson:"c"`
		FColor	  string        `bson:"fc"`
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
	
	results:= []searchOutput{}
	
	c.Find(query).Sort("n").All(&results)
	out, _  := json.MarshalIndent(results," ","  ")
	w.Write(out)
	s.Close()
}

func searchGameData(w http.ResponseWriter, r *http.Request) {
	
	type searchOutput struct {	
		Name      string        `bson:"n"`
		UserID    string        `bson:"id"`
		Score	  string        `bson:"sc"`
		Time	  string        `bson:"t"`
		Level	  string        `bson:"l"`
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
	
	results:= []searchOutput{}
	
	c.Find(query).Sort("n").All(&results)
	out, _  := json.MarshalIndent(results," ","  ")
	w.Write(out)
	s.Close()
}

func display(w http.ResponseWriter, r *http.Request) {
	
	type searchOutput struct {	
		Name      string        `bson:"n"`
		UserID    string        `bson:"id"`
		Perc	  string        `bson:"p"`
		FPerc	  string        `bson:"fp"`
		Color	  string        `bson:"c"`
		FColor	  string        `bson:"fc"`
		TimeStamp string		`bson:"ts"`
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
	
	results:= []searchOutput{}
	
	c.Find(query).Sort("ts").All(&results)
	out, _  := json.MarshalIndent(results," ","  ")
	w.Write(out)
	s.Close()
}



func insert(w http.ResponseWriter, r *http.Request) {
	
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
	
	query := bson.M {
		"n": name,
		"id": userid,
		"p": perc,
		"fp": fperc,
		"c": color,
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

func insertGameData(w http.ResponseWriter, r *http.Request) {
	
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
		"n": name,
		"id": userid,
		"sc": score,
		"t": time,
		"l": level,
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

func update(w http.ResponseWriter, r *http.Request) {


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


//func newPlayer() {
//}
//func newGameData(){
//}