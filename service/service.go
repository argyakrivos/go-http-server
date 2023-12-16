package service

import (
	"go-http-server/db"
	"go-http-server/models"
)

// BookService interface abstract the business logic,
// although in this example, there isn't any logic.
type BookService interface {
	GetAllBooks() ([]models.Book, error)
	GetBook(id string) (models.Book, error)
	CreateBook(newBook models.Book) (models.Book, error)
	UpdateBook(id string, updatedBook models.Book) (models.Book, error)
	DeleteBook(id string) (models.Book, error)
}

// DbService implements the BookService interface
type DbService struct {
	DB db.Database
}

func (s *DbService) GetAllBooks() ([]models.Book, error) {
	return s.DB.GetAllBooks()
}

func (s *DbService) GetBook(id string) (models.Book, error) {
	return s.DB.GetBook(id)
}

func (s *DbService) CreateBook(newBook models.Book) (models.Book, error) {
	return s.DB.CreateBook(newBook)
}

func (s *DbService) UpdateBook(id string, updatedBook models.Book) (models.Book, error) {
	return s.DB.UpdateBook(id, updatedBook)
}

func (s *DbService) DeleteBook(id string) (models.Book, error) {
	return s.DB.DeleteBook(id)
}
