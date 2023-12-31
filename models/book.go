package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Book represents a book entity.
type Book struct {
	Id     primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title  string             `json:"title,omitempty" bson:"title,omitempty"`
	Author string             `json:"author,omitempty" bson:"author,omitempty"`
}
