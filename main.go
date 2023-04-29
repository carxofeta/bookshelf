package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type Book struct {
	Title      string `bson:"title" json:"title" validate:"required"`
	Author     string `bson:"author" json:"author" validate:"required"`
	Publishing string `bson:"publishing" json:"publishing" validate:"required"`
	Edition    int    `bson:"edition,omitempty" json:"edition,omitempty"`
	Year       int    `bson:"year,omitempty" json:"year,omitempty"`
	ISBN       int64  `bson:"isbn,omitempty" json:"isbn,omitempty" validate:"required"`
	Rating     int    `bson:"rating,omitempty" json:"rating,omitempty"`
	Read       bool   `bson:"read" json:"read"`
}

func main() {

	//var books = make(map[string]Book)

	book := Book{
		Title:      "El primer hombre de Roma",
		Author:     "Colleen McCullough",
		Publishing: "Editorial Planeta S.A.",
		Edition:    1,
		Year:       2000,
		ISBN:       8408024000,
		Rating:     4,
		Read:       true,
	}

	//books[book.Title] = book

	fmt.Println("La libreria de Java...")
	//insertBook(book)
	addBook(book)

	insertFromFile("./books.json")

	//deleteBook("El primer hombre de Roma")

	libros, err := getBooks()
	if err != nil {
		log.Fatal(err)
	}
	for _, libro := range libros {
		bookJSON, err := json.Marshal(libro)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(bookJSON))
	}

}
