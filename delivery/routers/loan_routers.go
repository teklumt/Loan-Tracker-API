package routers

import (
	"loan-tracker-api/config/db"
	"loan-tracker-api/delivery/controllers"
	"loan-tracker-api/infrastracture"
	"loan-tracker-api/repository"
	"loan-tracker-api/usecase"

	"github.com/gin-gonic/gin"
)

func setUpLoanRoutes(router *gin.Engine) {
	LoanRepo := repository.NewLoanRepositoryImpl(db.LoanCollection)
	LoanUsecase := usecase.NewLoanUsecase(LoanRepo)
	LoanController := controllers.NewLoanController(LoanUsecase)
	// controllers.NewLoanController(LoanUsecase)

	Loan := router.Group("/loans")
	Loan.Use(infrastracture.AuthMiddleware())
	{
		Loan.POST("/", LoanController.CreateLoan)
		Loan.GET("/:id",LoanController.GetLoanByID)
	}
}