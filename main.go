package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// Book struct (Model)
type Book struct {
	ID     string  `json:"id"`
	Isbn   string  `json:"isbn"`
	Title  string  `json:"title"`
	Author *Author `json:"author"`
}

// Author struct
type Author struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

var books []Book

func homePage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Welcome to the HomePage!")
	log.Println("Endpoint Hit: homePage")
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	log.Println("getBooks called")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
	w.WriteHeader(http.StatusOK)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	log.Println("getBook called")
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // Get Params
	// Iterate through books and find with id
	for _, item := range books {
		if item.ID == params["id"] {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode((item))
			return
		}
	}
	// json.NewEncoder(w).Encode(&Book{})
	w.WriteHeader(http.StatusNotFound)
}

var addBookHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	log.Println("addBook called")
	w.Header().Set("Content-Type", "application/json")
	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)
	book.ID = strconv.Itoa(rand.Intn(1000000)) // mock data - not safe
	books = append(books, book)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(book)
})

var updateBookHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	log.Println("updateBook called")
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range books {
		if item.ID == params["id"] {
			books = append(books[:index], books[index+1:]...)
			var book Book
			_ = json.NewDecoder(r.Body).Decode(&book)
			book.ID = params["id"]
			books = append(books, book)
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(book)
			return
		}
	}
})

var deleteBookHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	log.Println("deleteBook called")
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range books {
		if item.ID == params["id"] {
			books = append(books[:index], books[index+1:]...)
			break
		}
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(books)
})

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)

	router.PathPrefix("/static").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	router.HandleFunc("/authenticate", authenticate)

	router.HandleFunc("/", homePage)
	router.HandleFunc("/books", getBooks).Methods("GET")
	router.HandleFunc("/books/{id}", getBook).Methods("GET")
	router.Handle("/books", authMiddleware(addBookHandler)).Methods("POST")
	router.Handle("/books/{id}", authMiddleware(updateBookHandler)).Methods("PUT")
	router.Handle("/books/{id}", authMiddleware(deleteBookHandler)).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", handlers.LoggingHandler(os.Stdout, router)))
}

func main() {

	// Hardcoded data - @todo: add database
	books = append(books, Book{ID: "1", Isbn: "438227", Title: "Book One", Author: &Author{Firstname: "John", Lastname: "Doe"}})
	books = append(books, Book{ID: "2", Isbn: "454555", Title: "Book Two", Author: &Author{Firstname: "Steve", Lastname: "Smith"}})
	handleRequests()
}
