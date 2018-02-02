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
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(results)
}
func MovieShow(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	movie_id := params["id"]

	if !bson.IsObjectIdHex(movie_id) {
		w.WriteHeader(404)
		return
	}
	oid := bson.ObjectIdHex(movie_id)					
	fmt.Println(movie_id)
	fmt.Println(oid)
	results := Movie{}
	err := collection.Find(oid).One(&results)
	fmt.Println(results)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(results)
	
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
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(movie_data)

}