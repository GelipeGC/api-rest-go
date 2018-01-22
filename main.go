package main

import (
	"fmt"
	"net/http"
	"log"
	)

	func main() {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
			fmt.Fprintf(w, "holamundo desde mi servidor con go")
		})

		server := http.ListenAndServe(":8080", nil)
		log.Fatal(server)
		fmt.Println("El servidor  esta  corriendo en http:localhost:8080")
	}