package service

import (
	"go-microservice/db"
	"go-microservice/models"
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

// BookServiceImpl implements the BookService interface
type BookServiceImpl struct {
	DB db.BookRepository
}

func NewBookServiceImpl(database db.BookRepository) BookService {
	return &BookServiceImpl{DB: database}
}

func (s *BookServiceImpl) GetAllBooks() ([]models.Book, error) {
	return s.DB.GetAllBooks()
}

func (s *BookServiceImpl) GetBook(id string) (models.Book, error) {
	return s.DB.GetBook(id)
}

func (s *BookServiceImpl) CreateBook(newBook models.Book) (models.Book, error) {
	return s.DB.CreateBook(newBook)
}

func (s *BookServiceImpl) UpdateBook(id string, updatedBook models.Book) (models.Book, error) {
	return s.DB.UpdateBook(id, updatedBook)
}

func (s *BookServiceImpl) DeleteBook(id string) (models.Book, error) {
	return s.DB.DeleteBook(id)
}
