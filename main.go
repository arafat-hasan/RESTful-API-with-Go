package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
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

	filter := bson.D{}

	cursor, err := query(mongoDataStore, "testCollection", filter)

	if err != nil {
		panic(err)
	}

	var results []bson.D

	if err := cursor.All(mongoDataStore.Context, &results); err != nil {

		panic(err)
	}

	json.NewEncoder(w).Encode(results)
	w.WriteHeader(http.StatusOK)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	log.Info("getBook called")
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	mongoDataStore := NewDatastore(cfg, log)

	filter := bson.D{
		{"isbn", params["id"]},
	}
	cursor, err := query(mongoDataStore, "testCollection", filter)

	if err != nil {
		panic(err)
	}

	var results []bson.D

	if err := cursor.All(mongoDataStore.Context, &results); err != nil {
		panic(err)
	}

	json.NewEncoder(w).Encode(results)
	w.WriteHeader(http.StatusOK)
}

var addBookHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	log.Println("addBook called")
	w.Header().Set("Content-Type", "application/json")
	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)
	mongoDataStore := NewDatastore(cfg, log)
	result, err := insertOne(mongoDataStore, "testCollection", book)
	if err != nil {
		panic(err)
	}

	fmt.Println("Result of InsertOne")
	fmt.Println(result.InsertedID)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(book)
})

var updateBookHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	log.Println("updateBook called")
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)

	filter := bson.D{
		{"isbn", params["id"]},
	}

	mongoDataStore := NewDatastore(cfg, log)

	updateBook := bson.M{
		"$set": book,
	}
	result, err := UpdateOne(mongoDataStore, "testCollection", filter, updateBook)

	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
})

var deleteBookHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	log.Println("deleteBook called")
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	mongoDataStore := NewDatastore(cfg, log)
	query := bson.D{
		{"isbn", params["id"]},
	}

	result, err := deleteOne(mongoDataStore, "testCollection", query)
	if err != nil {
		panic(err)
	}

	fmt.Println("No.of rows affected by DeleteOne()")
	fmt.Println(result.DeletedCount)

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
	// router.Handle("/books", authMiddleware(addBookHandler)).Methods("POST")
	// router.Handle("/books/{id}", authMiddleware(updateBookHandler)).Methods("PUT")
	// router.Handle("/books/{id}", authMiddleware(deleteBookHandler)).Methods("DELETE")

	router.HandleFunc("/books", addBookHandler).Methods("POST")
	router.HandleFunc("/books/{id}", updateBookHandler).Methods("PUT")
	router.HandleFunc("/books/{id}", deleteBookHandler).Methods("DELETE")

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

	log.Out = os.Stdout
	log.Info("Server Started!")

	handleRequests()

}
