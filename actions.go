package main

import (
	"log"
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
	)

	func getSession() *mgo.Session {
		session, err := mgo.Dial("mongodb://localhost")
	
		if err != nil {
			panic(err)
		}
	
		return session
	}
	var collection = getSession().DB("curso_go").C("movies");

	func responseMovie(w http.ResponseWriter, status int, results Movie){
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(results)
	}
	func responseMovies(w http.ResponseWriter, status int, results []Movie){
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(results)
	}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "El servidor  esta  corriendo en http:localhost:8080")
}

func MoviesList(w http.ResponseWriter, r *http.Request) {
	var results []Movie
	err := collection.Find(nil).Sort("-_id").All(&results)

	if err != nil {
		log.Fatal(err)
	}else {
		fmt.Println("Resultados: ", results)
	}
	responseMovies(w, 200, results)	
}
func MovieShow(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	movie_id := params["id"]

	if !bson.IsObjectIdHex(movie_id) {
		w.WriteHeader(404)
		return
	}
	oid := bson.ObjectIdHex(movie_id)					
	
	results := Movie{}
	err := collection.FindId(oid).One(&results)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	responseMovie(w, 200, results)	
	
}
func MovieAdd(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	
	var movie_data Movie
	err := decoder.Decode(&movie_data)

	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	err = collection.Insert(movie_data)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	responseMovie(w, 200, movie_data)	
}

func MovieUpdate(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	movie_id := params["id"]

	if !bson.IsObjectIdHex(movie_id) {
		w.WriteHeader(404)
		return
	}
	oid := bson.ObjectIdHex(movie_id)					
	decoder := json.NewDecoder(r.Body)
	
	var movie_data Movie
	err := decoder.Decode(&movie_data)

	if err != nil {
		w.WriteHeader(500)
		return
	}

	defer r.Body.Close()
	
	document := bson.M{"_id": oid}
	change := bson.M{"$set": movie_data}
	err = collection.Update(document, change)
	
	if err != nil {
		w.WriteHeader(404)
		return
	}
	responseMovie(w, 200, movie_data)	
	
}

type Message struct {
	Status string `json:"status"`
	Message string `json:"message"`
}

func (this *Message) setStatus(data string)  {
	this.Status = data
}

func (this *Message) setMessage(data string)  {
	this.Message = data
}

func MovieRemove(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	movie_id := params["id"]

	if !bson.IsObjectIdHex(movie_id) {
		w.WriteHeader(404)
		return
	}
	oid := bson.ObjectIdHex(movie_id)					
	
	err := collection.RemoveId(oid)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	message := new(Message)

	message.setStatus("success")
	message.setMessage("La pelicula con id" +movie_id+"ha sido borrada correctamente")
	results := message
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(results)
	
}