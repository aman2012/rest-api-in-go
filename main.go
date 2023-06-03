package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// data model
type Book struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Price  int    `json:"price"`
}

// temporary database
var DB *sql.DB

// middleware
func (c *Book) IsEmpty() bool {
	return c.Title == ""
}

// server entry point
func main() {
	fmt.Println("Api")
	//establishing db connection
	db, err := sql.Open("mysql", "root:aman2012@tcp(localhost:3306)/sakila")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	// Test the connection
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	DB = db
	r := mux.NewRouter()
	//seeding

	//routing
	r.HandleFunc("/", serveHome).Methods("GET")
	r.HandleFunc("/books", getAllBooks).Methods("GET")
	r.HandleFunc("/books/{id}", getOneBook).Methods("GET")
	r.HandleFunc("/books", createBook).Methods("POST")
	r.HandleFunc("/books/{id}", updateBook).Methods("PUT")
	r.HandleFunc("/books/{id}", deleteOneBook).Methods("DELETE")

	//listen to a port

	log.Fatal(http.ListenAndServe(":9090", r))

}

// homepage
func serveHome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1> Hello to my home page</h1>"))
}

// fetch all books
func getAllBooks(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get all the list of books")
	w.Header().Set("Content-Type", "application/json")

	rows, err := DB.Query("SELECT * FROM books")
	if err != nil {
		json.NewEncoder(w).Encode(err)
		return
	}

	var books []Book
	for rows.Next() {
		var book Book
		if err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Price); err != nil {
			json.NewEncoder(w).Encode(err)
			return
		}
		books = append(books, book)
	}

	json.NewEncoder(w).Encode(books)
}

// fetch one book
func getBookByID(db *sql.DB, bookID string) (Book, error) {
	query := "SELECT * FROM books WHERE id = ?"
	row := DB.QueryRow(query, bookID)

	var book Book
	err := row.Scan(&book.ID, &book.Title, &book.Author, &book.Price)

	if err != nil {
		return Book{}, err
	}

	return book, nil
}

func getOneBook(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get one book with help of the id")
	w.Header().Set("Content-Type", "application/json")

	param := mux.Vars(r)

	// Fetch a book by ID
	bookID := param["id"]
	book, err := getBookByID(DB, bookID)
	if err != nil {
		log.Fatal(err)
		return
	}

	json.NewEncoder(w).Encode(book)
}

// create a book
func createBook(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Create one book")

	w.Header().Set("Content-Type", "application/json")
	if r.Body == nil {
		json.NewEncoder(w).Encode("Please pass some value")
		return
	}

	var book Book

	_ = json.NewDecoder(r.Body).Decode(&book)

	if book.IsEmpty() {
		json.NewEncoder(w).Encode("No data inside")
		return
	}

	rand.Seed(time.Now().Unix())
	book.ID = strconv.Itoa(rand.Intn(100))
	query := "INSERT INTO books (id,title, author,price) VALUES (?, ?,?,?)"
	_, err := DB.Exec(query, book.ID, book.Title, book.Author, book.Price)
	if err != nil {
		json.NewEncoder(w).Encode(err)
		return
	}

	json.NewEncoder(w).Encode(book)

	return
}

// update a book
func updateBook(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Update one book")
	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)
	w.Header().Set("Content-Type", "application/json")
	param := mux.Vars(r)

	query := "UPDATE books SET title = ? , author = ? , price = ? WHERE id = ?"
	_, err := DB.Exec(query, book.Title, book.Author, book.Price, param["id"])
	if err != nil {
		json.NewEncoder(w).Encode(err)
		return
	}

	json.NewEncoder(w).Encode(book)
}

// delete a book
func deleteOneBook(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Update one book")
	w.Header().Set("Content-Type", "application/json")

	param := mux.Vars(r)

	query := "DELETE FROM books WHERE id = ?"
	_, err := DB.Exec(query, param["id"])

	if err != nil {
		json.NewEncoder(w).Encode(err)
		return
	}

	json.NewEncoder(w).Encode("Deleted")

}
