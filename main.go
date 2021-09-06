package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/yaml.v2"
)

type Book struct {
	ID     string  `json:"id" bson:"id"`
	Isbn   string  `json:"isbn" bson:"isbn"`
	Title  string  `json:"title" bson:"title"`
	Author *Author `json:"author" bson:"author"`
}

type Author struct {
	Firstname string `json:"firstname" bson:"firstname"`
	Lastname  string `json:"lastname" bson:"lastname"`
}

var books []Book
var log = logrus.New()
var cfg Config

func homePage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Welcome to the HomePage!")
	log.Println("Endpoint Hit: homePage")
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	log.Println("getBooks called")
	w.Header().Set("Content-Type", "application/json")
	mongoDataStore := NewDatastore(cfg, log)
	var filter, option interface{}

	// filter  gets all document,
	// with maths field greater that 70
	filter = bson.D{}

	//  option remove id field from all documents
	option = bson.D{}

	cursor, err := query(mongoDataStore, "testCollection", filter, option)

	// handle the errors.
	if err != nil {
		panic(err)
	}

	fmt.Println(cursor)

	//var results []bson.D

	// to get bson object  from cursor,
	// returns error if any.
	//if err := cursor.All(mongoDataStore.Context, &results); err != nil {

	// handle the error
	//  panic(err)
	//}

	// printing the result of query.
	//fmt.Println("Query Reult")
	//for _, doc := range results {
	//  fmt.Println(doc)
	//}

	//json.NewEncoder(w).Encode(results)
	//w.WriteHeader(http.StatusOK)
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

type Config struct {
	Database struct {
		Name string `yaml:"name"`
		Host string `yaml:"host"`
	} `yaml:"database"`
}

func processError(err error) {
	fmt.Println(err)
	os.Exit(2)
}

func main() {

	// Hardcoded data - @todo: add database
	//book1 := Book{ID: "1", Isbn: "438227", Title: "Book One", Author: &Author{Firstname: "John", Lastname: "Doe"}}
	//book2 := Book{ID: "2", Isbn: "454555", Title: "Book Two", Author: &Author{Firstname: "Steve", Lastname: "Smith"}}

	f, err := os.Open(".config.yml")
	if err != nil {
		processError(err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		processError(err)
	}

	//fmt.Println(cfg)

	//logger := logrus.Logger{
	//  Out: os.Stdout,
	//}
	log.Out = os.Stdout
	handleRequests()

	//log.Printf("Log message")

	//fmt.Println(mongoDataStore)

	//defer close(client, ctx, cancel)

	//_, err = insertOne(mongoDataStore, "testCollection", book1)
	//if err != nil {
	//  panic(err)
	//}
	//println("book2 inserted")

	//_, err = insertOne(mongoDataStore, "testCollection", book2)
	//if err != nil {
	//  panic(err)
	//}
	//println("book2 inserted")

}
