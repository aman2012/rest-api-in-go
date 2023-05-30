package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type book struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Author string  `json:"author"`
	Price  float64 `json:"price"`
}

var books = []book{
	{ID: "1", Title: "Harry Potter", Author: "JK Rowling", Price: 700},
	{ID: "2", Title: "Percy Jackson", Author: "Rick Riordan", Price: 500},
	{ID: "3", Title: "Secret Seven", Author: "Enid Blyton", Price: 300},
	{ID: "4", Title: "The Room on the Roof", Author: "Ruskin Bond", Price: 250},
	{ID: "5", Title: "Tom Sawyer", Author: "Mark Twain", Price: 675},
}

func getBook(context *gin.Context) {
	id := context.Param("id")
	book, err := getBookById(id)

	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book Not Found"})
		return

	}
	context.IndentedJSON(http.StatusOK, book)

}

func getBooks(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, books)

}

func addBook(context *gin.Context) {
	var newBook book

	if err := context.BindJSON(&newBook); err != nil {
		return
	}

	books = append(books, newBook)

	context.IndentedJSON(http.StatusCreated, newBook)

}

func getBookById(id string) (*book, error) {
	for i, j := range books {
		if j.ID == id {
			return &books[i], nil
		}
	}
	return nil, errors.New("book not found")
}

func updatePrice(context *gin.Context) {
	id := context.Param("id")
	book, err := getBookById(id)

	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book Not Found"})
		return

	}

	book.Price = book.Price * 0.75
	context.IndentedJSON(http.StatusCreated, book)

}

func deleteBook(context *gin.Context) {
	id := context.Param("id")
	book, err := getBookById(id)

	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book Not Found"})
		return

	}
	for i, j := range books {
		if j.ID == id {
			books = append(books[:i], books[i+1:]...)

		}
	}
	context.IndentedJSON(http.StatusOK, book)

}

func main() {
	router := gin.Default()
	router.GET("/books", getBooks)
	router.GET("/books/:id", getBook)
	router.PATCH("/books/:id", updatePrice)
	router.POST("/books", addBook)
	router.DELETE("/books/:id", deleteBook)
	router.Run("localhost:9090")

}
