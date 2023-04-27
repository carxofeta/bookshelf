package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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

func getCollection() (*mongo.Collection, error) {
	client, err := connect()
	if err != nil {
		return nil, err
	}
	collection := client.Database("library").Collection("books")
	return collection, nil
}

func addBook(bookJSON Book) {
	collection, err := getCollection()
	if err != nil {
		fmt.Println(err)
	}
	// Insertar el libro en la base de datos MongoDB usando UpdateOne por si ya existe
	filter := bson.M{"title": bookJSON.Title, "isbn": bookJSON.ISBN}
	update := bson.M{"$setOnInsert": bookJSON}
	opts := options.Update().SetUpsert(true)
	a, err := collection.UpdateOne(context.Background(), filter, update, opts)
	if err != nil {
		log.Fatal(err)
	}

	if mongo.UpdateResult(*a).MatchedCount == 1 {
		fmt.Println("El libro especificado ya existe")
	}

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

func insertFromFile(filePath string) {
	collection, err := getCollection()
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	var books []Book
	if err := json.Unmarshal(file, &books); err != nil {
		log.Fatal(err)
	}

	// Crear una lista de operaciones "bulk write" para insertar los libros
	var ops []mongo.WriteModel
	for _, book := range books {
		filter := bson.M{"title": book.Title, "isbn": book.ISBN}
		update := bson.M{"$setOnInsert": book}
		op := mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update).SetUpsert(true)
		ops = append(ops, op)
	}

	// Ejecutar las operaciones "bulk write"
	result, err := collection.BulkWrite(context.Background(), ops) //, options.BulkWrite().SetOrdered(true))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Insertados %d libros\n", result.InsertedCount)
	fmt.Printf("Upsertados %d libros\n", result.UpsertedCount)
}

func deleteBook(book string) {
	collection, err := getCollection()
	if err != nil {
		fmt.Println(err)
	}
	filter := bson.M{"title": book}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		fmt.Println(err)
	}
	if result.DeletedCount == 0 {
		fmt.Println("El libro solicitado no se encuentra en la colección")
	} else {
		fmt.Printf("El título \"%s\" ha sido eliminado de la colección\n", book)
	}
}

//TODO: Update book
