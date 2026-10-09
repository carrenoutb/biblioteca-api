package service

import (
	"biblioteca-api/internal/models"
	"biblioteca-api/internal/repository"
)

type BookService struct {
	bookRepo *repository.BookRepository
}

func NewBookService(bookRepo *repository.BookRepository) *BookService {
	return &BookService{bookRepo: bookRepo}
}

func (s *BookService) List(category, search string) ([]models.Book, error) {
	return s.bookRepo.FindAll(category, search)
}

func (s *BookService) Get(id int) (*models.Book, error) {
	return s.bookRepo.FindByID(id)
}

func (s *BookService) Create(req models.BookRequest) (*models.Book, error) {
	book := &models.Book{
		Title:           req.Title,
		Author:          req.Author,
		ISBN:            req.ISBN,
		Category:        req.Category,
		CoverURL:        req.CoverURL,
		TotalCopies:     req.TotalCopies,
		AvailableCopies: req.TotalCopies,
	}
	if err := s.bookRepo.Create(book); err != nil {
		return nil, err
	}
	return book, nil
}

func (s *BookService) Update(id int, req models.BookRequest) (*models.Book, error) {
	book, err := s.bookRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	book.Title = req.Title
	book.Author = req.Author
	book.ISBN = req.ISBN
	book.Category = req.Category
	book.CoverURL = req.CoverURL
	book.TotalCopies = req.TotalCopies

	if err := s.bookRepo.Update(book); err != nil {
		return nil, err
	}
	return book, nil
}

func (s *BookService) Delete(id int) error {
	return s.bookRepo.Delete(id)
}
