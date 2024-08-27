package domain

import (
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Loan struct {
	ID            primitive.ObjectID  `json:"id" bson:"_id,omitempty"`
	Title         string              `json:"title"`
	Description   string              `json:"description"`
	BorrowerID    string              `json:"borrower_id"`
	BorrowerName  string              `json:"borrower_name"`
	Amount        float64             `json:"amount"`
	CreatedAt     primitive.Timestamp `json:"created_at" bson:"createdAt"`
	UpdatedAt     primitive.Timestamp `json:"-" bson:"updatedAt"`
	Status        string              `json:"status"`
}

type LoanResponse struct {
	ID           primitive.ObjectID  `json:"id" bson:"_id,omitempty"`
	Title        string              `json:"title"`
	Description  string              `json:"description"`
	BorrowerID   string              `json:"borrower_id"`
	BorrowerName string              `json:"borrower_name"`
	Amount       float64             `json:"amount"`
	CreatedAt    primitive.Timestamp `json:"created_at" bson:"createdAt"`
	Status       string              `json:"status"`
}

func (l *Loan) MarshalJSON() ([]byte, error) {
	return json.Marshal(&LoanResponse{
		ID:           l.ID,
		Title:        l.Title,
		Description:  l.Description,
		BorrowerID:   l.BorrowerID,
		BorrowerName: l.BorrowerName,
		Amount:       l.Amount,
		CreatedAt:    l.CreatedAt,
		Status:       l.Status,
	})
}


type LoanRepository interface {
	CreateLoan(loan Loan) (Loan, error)
	DeleteLoan(id string) (Loan, error)
	// UpdateLoan(loan Loan, loanID string) (Loan, error)
	GetLoanByID(id string) (Loan, error)
	GetLoans() ([]Loan, error)
	UpdateLoanStatus(newStatus, loanID string) (Loan, error)
	// GetUserLoans(borrowerID string) ([]Loan, error)

	// SearchLoan(loan Loan) (error)
}

type LoanUsecase interface {
	CreateLoan( loan Loan) (Loan, ErrorResponse)
	DeleteLoan( loanID string) (Loan, ErrorResponse)
	UpdateLoanStatus(newStatus, loanID string) (Loan, ErrorResponse)
	GetLoanByID(userID,Role, loanID string) (Loan, ErrorResponse)
	GetLoans() ([]Loan, ErrorResponse)
	// GetUserLoans(borrowerID string) ([]Loan, error)


}
