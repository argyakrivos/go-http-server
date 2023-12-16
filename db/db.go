package db

import (
	"context"
	"errors"
	"fmt"
	"go-http-server/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Database interface abstracts the database interactions.
type Database interface {
	GetAllBooks() ([]models.Book, error)
	GetBook(id string) (models.Book, error)
	CreateBook(newBook models.Book) (models.Book, error)
	UpdateBook(id string, updatedBook models.Book) (models.Book, error)
	DeleteBook(id string) (models.Book, error)
}

// MongoDB implements the Database interface.
type MongoDB struct {
	Collection *mongo.Collection
}

func (m *MongoDB) GetAllBooks() ([]models.Book, error) {
	cursor, err := m.Collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, fmt.Errorf("error fetching books from MongoDB: %w", err)
	}
	defer cursor.Close(context.Background())

	var books []models.Book
	if err := cursor.All(context.Background(), &books); err != nil {
		return nil, fmt.Errorf("error decoding books from MongoDB: %w", err)
	}

	return books, nil
}

func (m *MongoDB) GetBook(id string) (models.Book, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Book{}, errors.New("invalid book ID")
	}

	var book models.Book
	err = m.Collection.FindOne(context.Background(), bson.M{"_id": objectId}).Decode(&book)
	if err != nil {
		return models.Book{}, errors.New("book not found")
	}

	return book, nil
}

func (m *MongoDB) CreateBook(newBook models.Book) (models.Book, error) {
	result, err := m.Collection.InsertOne(context.Background(), newBook)
	if err != nil {
		return models.Book{}, fmt.Errorf("error inserting book into MongoDB: %w", err)
	}

	newBook.Id = result.InsertedID.(primitive.ObjectID)
	return newBook, nil
}

func (m *MongoDB) UpdateBook(id string, updatedBook models.Book) (models.Book, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Book{}, errors.New("invalid book ID")
	}

	var existingBook models.Book
	err = m.Collection.FindOne(context.Background(), bson.M{"_id": objectId}).Decode(&existingBook)
	if err != nil {
		return models.Book{}, errors.New("book not found")
	}

	setUpdates := bson.M{}
	if updatedBook.Title != "" && updatedBook.Title != existingBook.Title {
		setUpdates["title"] = updatedBook.Title
	}
	if updatedBook.Author != "" && updatedBook.Author != existingBook.Author {
		setUpdates["author"] = updatedBook.Author
	}

	if len(setUpdates) == 0 {
		return existingBook, nil
	}

	_, err = m.Collection.UpdateOne(context.Background(), bson.M{"_id": objectId}, bson.M{"$set": setUpdates})
	if err != nil {
		return models.Book{}, fmt.Errorf("error updating book in MongoDB: %w", err)
	}

	err = m.Collection.FindOne(context.Background(), bson.M{"_id": objectId}).Decode(&updatedBook)
	if err != nil {
		return models.Book{}, fmt.Errorf("error fetching updated book from MongoDB: %w", err)
	}

	return updatedBook, nil
}

func (m *MongoDB) DeleteBook(id string) (models.Book, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Book{}, errors.New("invalid book ID")
	}

	var book models.Book
	err = m.Collection.FindOneAndDelete(context.Background(), bson.M{"_id": objectId}).Decode(&book)
	if err != nil {
		return models.Book{}, errors.New("book not found")
	}

	return book, nil
}
