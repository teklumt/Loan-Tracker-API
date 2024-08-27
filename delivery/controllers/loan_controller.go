package controllers

import (
	"loan-tracker-api/domain"

	"github.com/gin-gonic/gin"
)

type LoanController struct {
	loanUsecase domain.LoanUsecase
}

func NewLoanController(loanUsecase domain.LoanUsecase) *LoanController {
	return &LoanController{
		loanUsecase: loanUsecase,
	}
}

func (l *LoanController) CreateLoan(c *gin.Context) {
	var loan domain.Loan
	userID := c.GetString("user_id")
	usename := c.GetString("username")
	if err := c.ShouldBindJSON(&loan); err != nil {

		c.JSON(400, gin.H{
			"code":    400,
			"message": "All fields are required",
		})
		return
	}
	loan.BorrowerID = userID
	loan.BorrowerName = usename

	createdLoan , err := l.loanUsecase.CreateLoan(loan)
	if err.Message != "" {
		c.JSON(err.StatusCode, err)
		return
	}

	c.JSON(200, gin.H{
		"code":    200,
		"message": "Loan created successfully",
		"data":    createdLoan,
	})
}


func (l *LoanController)GetLoanByID(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("user_id")
	Role := c.GetString("role")
	if id == "" {
		c.JSON(400, gin.H{
			"code":    400,
			"message": "ID is required",
		})
		return
	}

	loan, err := l.loanUsecase.GetLoanByID(userID, Role ,id)
	if err.Message != "" {
		c.JSON(err.StatusCode, err)
		return
	}

	c.JSON(200, gin.H{
		"code":    200,
		"message": "Loan fetched successfully",
		"data":    loan,
	})

	
}

func (l *LoanController) GetLoans(c *gin.Context) {
	loans, err := l.loanUsecase.GetLoans()
	if err.Message != ""  {
		c.JSON(500, gin.H{
			"code":    500,
			"message": "Internal Server Error",
		})
		return
	}

	c.JSON(200, gin.H{
		"code":    200,
		"message": "Loans fetched successfully",
		"data":    loans,
	})
}


func (l *LoanController) UpdateLoanStatus(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(400, gin.H{
			"code":    400,
			"message": "Loan ID is required",
		})
		return
	}

	var updateData struct {
		Status string `json:"status"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(400, gin.H{
			"code":    400,
			"message": "Invalid request body",
		})
		return
	}

	if updateData.Status == "" {
		c.JSON(400, gin.H{
			"code":    400,
			"message": "Status is required",
		})
		return
	}

	

	updatedLoan, err := l.loanUsecase.UpdateLoanStatus(updateData.Status, id)
	if err.Message != "" {
		c.JSON(err.StatusCode, gin.H{
			"code":    err.StatusCode,
			"message": err.Message,
		})
		return
	}

	c.JSON(200, gin.H{
		"code":    200,
		"message": "Loan status updated successfully",
		"data":    updatedLoan,
	})
}


func (l *LoanController) DeleteLoan(c *gin.Context) {
	id := c.Param("id")
	deletedLoan, err := l.loanUsecase.DeleteLoan(id)
	if err.Message != "" {
		c.JSON(err.StatusCode, err)
		return
	}

	c.JSON(200, gin.H{
		"code":    200,
		"message": "Loan deleted successfully",
		"data":    deletedLoan,
	})
}