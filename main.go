package main

//.exe file used to run our server with ./build
//the formatter has print commands that we can print to the console
import (
	"encoding/json"
	"log"      //to log errors
	"net/http" //to work with http https://golang.org/pkg/net/http/ & https://www.integralist.co.uk/posts/understanding-golangs-func-type/

	//when we creat a book, need to add id to be generated as random
	//string converter
	"math/rand"
	"strconv"

	"github.com/gorilla/mux" //router
)

//Book create our book structs whis is our MODEL. Like a class. it can have properties and methods
type Book struct {
	ID     string  `json:"id"` //reads the response body in the client after making the request
	Isbn   string  `json:"isbn"`
	Title  string  `json:"title"`
	Author *Author `json:"author"` //have to create a struct for it
}

//Author struct
type Author struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

//init books var as a slice book struct (slice variable length arrary)
var books []Book //slice because its a collection of books

//every func we create has to have these two params.  (request, response)
func getBooks(w http.ResponseWriter, r *http.Request) {
	//set the header of content type to applications/json
	w.Header().Set("Content-Type", "applications/json") //set content type to json. or else its going to be served as text
	json.NewEncoder(w).Encode(books)

}

//get single book
func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "applications/json")
	params := mux.Vars(r)        //Get params
	for _, item := range books { //item is iterator; _ ignores the key part of key value
		if item.ID == params["id"] { //the id that is in the URL
			json.NewEncoder(w).Encode(item)
			return
		}
	} //call this method to write JSON to the server, output the item
	//Loop through books and find with ID
	json.NewEncoder(w).Encode(&Book{})
}

func createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)
	book.ID = strconv.Itoa(rand.Intn(10000000)) //will convert and int to a string. It's a mock ID
	books = append(books, book)
	json.NewEncoder(w).Encode(book)
	// add this to postman using post request{
	// 	"isbn":"123234",
	// 	"title":"book 3",
	// 	"author":{
	// 		"firstname:"carol";
	// 		"lastname":"williams"
	// 		}
	// }
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range books {
		if item.ID == params["id"] {
			books = append(books[:index], books[index+1:]...)
			var book Book
			_ = json.NewDecoder(r.Body).Decode(&book)
			book.ID = params["id"] //will convert and int to a string. It's a mock ID
			books = append(books, book)
			json.NewEncoder(w).Encode(book)
			return
		}
	}
	json.NewEncoder(w).Encode(books)

}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range books {
		if item.ID == params["id"] {
			books = append(books[:index], books[index+1:]...) // appends evrything up to index minus index
			break
		}
	}
	json.NewEncoder(w).Encode(books)
}

func main() {
	//initialize the mux router
	r := mux.NewRouter()

	//mock data @todo-implement DB
	books = append(books, Book{ID: "1", Isbn: "12234", Title: "Book 1",
		Author: &Author{Firstname: "John", Lastname: "Matt"}}) //we want to append to books

	books = append(books, Book{ID: "2", Isbn: "1223433", Title: "Book 2",
		Author: &Author{Firstname: "Johny", Lastname: "Matty"}}) //we want to append to books

	books = append(books, Book{ID: "3", Isbn: "122344", Title: "Book 3",
		Author: &Author{Firstname: "Johnyy", Lastname: "Mattyy"}}) //we want to append to books

	// Create route handlers which will establish our endpoints for our api
	//getBooks is the func we want to run  .Methods = which http method to use
	r.HandleFunc("/api/books", getBooks).Methods("GET")
	r.HandleFunc("/api/books/{id}", getBook).Methods("GET")
	r.HandleFunc("/api/books", createBook).Methods("POST")
	r.HandleFunc("/api/books/{id}", updateBook).Methods("PUT") //update will be put request
	r.HandleFunc("/api/books/{id}", deleteBook).Methods("DELETE")

	//TO run the server, use http and the method listen and serve, wrap it in log.fatal in case it fails, we get an error
	log.Fatal(http.ListenAndServe(":8000", r))
}
