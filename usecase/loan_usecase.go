package usecase

import "loan-tracker-api/domain"

type LoanUsecaseImpl struct {
	loanRepo domain.LoanRepository
}

func NewLoanUsecase(loanRepo domain.LoanRepository) domain.LoanUsecase {
	return &LoanUsecaseImpl{
		loanRepo: loanRepo,
	}
}

func (l *LoanUsecaseImpl) CreateLoan(loan domain.Loan) (domain.Loan, domain.ErrorResponse) {

	if loan.Title == "" || loan.Description == "" || loan.BorrowerID == "" || loan.BorrowerName == "" || loan.Amount == 0 {
		return domain.Loan{}, domain.ErrorResponse{
			StatusCode: 400,
			Message:    "All fields are required",
		}
		
	}

	loan.Status = "pending"

	createdLoan, err := l.loanRepo.CreateLoan(loan)
	if err != nil {
		return domain.Loan{}, domain.ErrorResponse{
			StatusCode: 500,
			Message:    "Internal Server Error",
		}
	}

	return createdLoan, domain.ErrorResponse{}
	
}


func (l *LoanUsecaseImpl) GetLoanByID( userID,Role ,id string) (domain.Loan, domain.ErrorResponse) {
	if id == "" {
		return domain.Loan{}, domain.ErrorResponse{
			StatusCode: 400,
			Message:    "ID is required",
		}
	}

	loan, err := l.loanRepo.GetLoanByID(id)
	if err != nil {
		return domain.Loan{}, domain.ErrorResponse{
			StatusCode: 500,
			Message:    "Internal Server Error",
		}
	}
	if loan.BorrowerID != userID &&  Role != "admin" {
		return domain.Loan{}, domain.ErrorResponse{
			StatusCode: 401,
			Message:    "Unauthorized",
		}
	}

	return loan, domain.ErrorResponse{}
}

func (l *LoanUsecaseImpl) GetLoans() ([]domain.Loan, domain.ErrorResponse) {
	loans, err := l.loanRepo.GetLoans()
	if err != nil {
		return []domain.Loan{}, domain.ErrorResponse{
			StatusCode: 500,
			Message:    "Internal Server Error",
		}
	}

	return loans, domain.ErrorResponse{}
}


func (l *LoanUsecaseImpl) UpdateLoanStatus(newStatus, loanID string) (domain.Loan, domain.ErrorResponse) {
	if newStatus == "" || loanID == "" {
		return domain.Loan{}, domain.ErrorResponse{
			StatusCode: 400,
			Message:    "All fields are required",
		}
	}

	if newStatus != "approved" && newStatus != "rejected" && newStatus != "pending" {
		return domain.Loan{}, domain.ErrorResponse{
			StatusCode: 400,
			Message:    "Invalid status",
		}
	}

	_, err := l.loanRepo.GetLoanByID(loanID)
	if err != nil {
		return domain.Loan{}, domain.ErrorResponse{
			StatusCode: 500,
			Message:    "Internal Server Error",
		}
	}

	
	_, err = l.loanRepo.UpdateLoanStatus(newStatus, loanID)
	if err != nil {
		return domain.Loan{}, domain.ErrorResponse{
			StatusCode: 500,
			Message:    "Internal Server Error",
		}
	}
	loan, err := l.loanRepo.GetLoanByID(loanID)
	if err != nil {
		return domain.Loan{}, domain.ErrorResponse{
			StatusCode: 500,
			Message:    "Internal Server Error",
		}
	}

	return loan, domain.ErrorResponse{}
}

func (l *LoanUsecaseImpl) DeleteLoan(loanID string) (domain.Loan, domain.ErrorResponse) {
	if loanID == "" {
		return domain.Loan{}, domain.ErrorResponse{
			StatusCode: 400,
			Message:    "ID is required",
		}
	}

	loan, err := l.loanRepo.DeleteLoan(loanID)
	if err != nil {
		return domain.Loan{}, domain.ErrorResponse{
			StatusCode: 500,
			Message:    "Internal Server Error",
		}
	}

	return loan, domain.ErrorResponse{}
}