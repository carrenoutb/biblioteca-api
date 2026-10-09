package service

import (
	"errors"
	"time"

	"biblioteca-api/internal/models"
	"biblioteca-api/internal/repository"
)

var (
	ErrBookNotAvailable = errors.New("no hay copias disponibles de este libro")
	ErrAlreadyBorrowed  = errors.New("ya tienes un préstamo activo de este libro")
	ErrLoanNotFound     = errors.New("préstamo no encontrado")
	ErrAlreadyReturned  = errors.New("este préstamo ya fue devuelto")
)

const loanDurationDays = 14

type LoanService struct {
	loanRepo *repository.LoanRepository
	bookRepo *repository.BookRepository
}

func NewLoanService(loanRepo *repository.LoanRepository, bookRepo *repository.BookRepository) *LoanService {
	return &LoanService{loanRepo: loanRepo, bookRepo: bookRepo}
}

func (s *LoanService) Borrow(userID int, req models.LoanRequest) (*models.Loan, error) {
	book, err := s.bookRepo.FindByID(req.BookID)
	if err != nil {
		return nil, err
	}
	if book.AvailableCopies <= 0 {
		return nil, ErrBookNotAvailable
	}

	alreadyBorrowed, err := s.loanRepo.HasActiveLoan(userID, req.BookID)
	if err != nil {
		return nil, err
	}
	if alreadyBorrowed {
		return nil, ErrAlreadyBorrowed
	}

	dueAt := time.Now().AddDate(0, 0, loanDurationDays)
	loan, err := s.loanRepo.Create(userID, req.BookID, dueAt)
	if err != nil {
		return nil, err
	}

	if err := s.bookRepo.DecrementAvailable(req.BookID); err != nil {
		return nil, err
	}

	return loan, nil
}

func (s *LoanService) Return(loanID int) error {
	loan, err := s.loanRepo.FindByID(loanID)
	if err != nil {
		return ErrLoanNotFound
	}
	if loan.Status == "returned" {
		return ErrAlreadyReturned
	}

	if err := s.loanRepo.MarkReturned(loanID); err != nil {
		return err
	}
	return s.bookRepo.IncrementAvailable(loan.BookID)
}

func (s *LoanService) ListByUser(userID int) ([]models.LoanDetail, error) {
	return s.loanRepo.FindByUser(userID)
}

func (s *LoanService) ListAll() ([]models.LoanDetail, error) {
	return s.loanRepo.FindAll()
}
