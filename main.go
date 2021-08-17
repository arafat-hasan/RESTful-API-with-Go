package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Book struct {
	Id     string `json:id`
	Title  string `json:title`
	Author string `json:author`
	Year   string `json:year`
}

var books []Book

func main() {
	router := mux.NewRouter()

	books = append(books,
		Book{Id: "1", Title: "Golang Intro", Author: "Go and Come", Year: "2020"},
		Book{Id: "2", Title: "Goroutines", Author: "Routine Chowdhury", Year: "2011"},
		Book{Id: "3", Title: "Go Routers", Author: "Router Haque", Year: "2013"},
		Book{Id: "4", Title: "Go Concurency", Author: "Mohammad Concurency", Year: "2015"},
		Book{Id: "5", Title: "Go For Good", Author: "Good Tarafdar", Year: "2017"})

	router.HandleFunc("/books", getBooks).Methods("GET")
	router.HandleFunc("/books/{id}", getBook).Methods("GET")
	router.HandleFunc("/books", addBook).Methods("POST")
	router.HandleFunc("/books", updateBook).Methods("PUT")
	router.HandleFunc("/book/{id}", removeBook).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", router))
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	log.Println("getBooks called")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	log.Println("getBook called")
	params := mux.Vars(r) // Get Params
	// Iterate through books and find with id
	for _, item := range books {
		if item.Id == params["id"] {
			json.NewEncoder(w).Encode((item))
			return
		}
	}
	json.NewEncoder(w).Encode(&Book{})
	// w.WriteHeader(http.StatusNotFound)
}

func addBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	log.Println("addBook called")
	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)
	book.Id = strconv.Itoa(rand.Intn(1000000)) // mock data - not safe
	books = append(books, book)
	json.NewEncoder(w).Encode(book)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	log.Println("updateBook called")
}

func removeBook(w http.ResponseWriter, r *http.Request) {
	log.Println("removeBook called")
}
