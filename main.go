package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type book struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Quantity int    `json:"quantity"`
}

var books = []book{
	{ID: "1", Title: "In Search of Lost Time", Author: "Marcel Proust", Quantity: 2},
	{ID: "2", Title: "The Great Gatsby", Author: "F. Scott Fitzgerald", Quantity: 5},
	{ID: "3", Title: "War and Peace", Author: "Leo Tolstoy", Quantity: 6},
}

func getBooks(context *gin.Context) {
	fmt.Println("Getting Books")
	context.IndentedJSON(http.StatusOK, books)
}

func bookById(context *gin.Context) {
	id := context.Param("id")
	book, err := getBookById(id)

	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found."})
		return
	}

	context.IndentedJSON(http.StatusOK, book)
}

func checkoutBook(context *gin.Context) {
	id, ok := context.GetQuery("id")

	if !ok {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing id query parameter."})

		return
	}

	book, err := getBookById(id)

	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found."})

		return
	}

	if book.Quantity <= 0 {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Book not found."})

		return
	}

	book.Quantity -= 1
	context.IndentedJSON(http.StatusOK, book)
}

func returnBook(context *gin.Context) {
	id, ok := context.GetQuery("id")

	if !ok {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing id query parameter."})

		return
	}

	book, err := getBookById(id)

	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found."})

		return
	}

	if book.Quantity <= 0 {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Book not found."})

		return
	}

	book.Quantity += 1
	context.IndentedJSON(http.StatusOK, book)
}

func getBookById(id string) (*book, error) {
	for i, b := range books {
		if b.ID == id {
			return &books[i], nil
		}
	}

	return nil, errors.New("book not found")
}

func createBook(context *gin.Context) {
	var newBook book

	if err := context.BindJSON(&newBook); err != nil {
		return
	}

	books = append(books, newBook)
	context.IndentedJSON(http.StatusCreated, newBook)
}

func main() {
	router := gin.Default()
	router.GET("/books", getBooks)
	router.GET("books/:id", bookById)
	router.POST("/books", createBook)
	router.PATCH("/checkout", checkoutBook)
	router.PATCH("/return", returnBook)
	router.Run("localhost:8080")
}
