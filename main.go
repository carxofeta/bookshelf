package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// func connection() {

// 	const uri = "mongodb://localhost:27017"

// 	// Use the SetServerAPIOptions() method to set the Stable API version to 1
// 	serverAPI := options.ServerAPI(options.ServerAPIVersion1)

// 	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

// 	// Create a new client and connect to the server
// 	client, err := mongo.Connect(context.TODO(), opts)
// 	if err != nil {
// 		panic(err)
// 	}

// 	defer func() {
// 		if err = client.Disconnect(context.TODO()); err != nil {
// 			panic(err)
// 		}
// 	}()

// 	// Send a ping to confirm a successful connection
// 	var result bson.M
// 	if err := client.Database("library").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Decode(&result); err != nil {
// 		panic(err)
// 	}
// 	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")
//}

func connect() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}
	return client, nil
}

type Book struct {
	Title      string `json:"title"`
	Author     string `json:"author"`
	Publishing string `json:"publishing"`
	Edition    int    `json:"edition,omitempty"`
	Year       int    `json:"year,omitempty"`
	ISBN       string `json:"ISBN,omitempty"`
	Rating     int    `json:"rating,omitempty"`
	Read       bool   `json:"read"`
}

func main() {

	//var books = make(map[string]Book)

	book := Book{
		Title:      "La caída del Imperio Romano. Las causas militares",
		Author:     "Arther Ferrill",
		Publishing: "Biblioteca Edaf",
		Edition:    1,
		Year:       1998,
		ISBN:       "8441403988",
		Read:       false,
	}

	//books[book.Title] = book

	fmt.Println("La libreria de Java...")
	//connection()
	//getCollection()
	//createBook(book)
	getBooks()
}

func getCollection() (*mongo.Collection, error) {
	client, err := connect()
	if err != nil {
		return nil, err
	}
	collection := client.Database("library").Collection("books")
	return collection, nil
}

func createBook(book Book) error {
	collection, err := getCollection()
	if err != nil {
		return err
	}
	_, err = collection.InsertOne(context.Background(), book)
	if err != nil {
		return err
	}
	return nil
}

func getBooks() ([]Book, error) {
	collection, err := getCollection()
	if err != nil {
		return nil, err
	}
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	var books []Book
	for cursor.Next(context.Background()) {
		var book Book
		err := cursor.Decode(&book)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return books, nil
}
