package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Book struct {
	Id     int    `json:id`
	Title  string `json:title`
	Author string `json:author`
	Year   string `json:year`
}

var books []Book

func main() {
	router := mux.NewRouter()

	books = append(books,
		Book{Id: 1, Title: "Golang Intro", Author: "Go and Come", Year: "2020"},
		Book{Id: 2, Title: "Goroutines", Author: "Routine Chowdhury", Year: "2011"},
		Book{Id: 3, Title: "Go Routers", Author: "Router Haque", Year: "2013"},
		Book{Id: 4, Title: "Go Concurency", Author: "Mohammad Concurency", Year: "2015"},
		Book{Id: 5, Title: "Go For Good", Author: "Good Tarafdar", Year: "2017"})

	router.HandleFunc("/books", getBooks).Methods("GET")
	router.HandleFunc("/books/{id}", getBook).Methods("GET")
	router.HandleFunc("/books", addBook).Methods("POST")
	router.HandleFunc("/books", updateBook).Methods("PUT")
	router.HandleFunc("/book/{id}", removeBook).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", router))
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	log.Println("getBooks called")
	json.NewEncoder(w).Encode(books)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	log.Println("getBook called")
}

func addBook(w http.ResponseWriter, r *http.Request) {
	log.Println("addBook called")
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	log.Println("updateBook called")
}

func removeBook(w http.ResponseWriter, r *http.Request) {
	log.Println("removeBook called")
}
