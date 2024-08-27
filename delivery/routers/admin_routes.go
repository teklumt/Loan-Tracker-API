// delivery/routers/admin_routes.go
package routers

import (
	"loan-tracker-api/config/db"
	"loan-tracker-api/delivery/controllers"
	"loan-tracker-api/infrastracture"
	"loan-tracker-api/repository"
	"loan-tracker-api/usecase"

	"github.com/gin-gonic/gin"
)

func setUpAdminRoutes(router *gin.Engine) {
    // Initialize repository with database collection
    userRepo := repository.NewUserRepositoryImpl(db.UserCollection)
    loanRepo := repository.NewLoanRepositoryImpl(db.LoanCollection)
    logRepo := repository.NewLogRepositoryImpl(db.LogCollection)

    // Initialize token generator and password service
    tokenGen := infrastracture.NewTokenGenerator()
    passwordSvc := infrastracture.NewPasswordService()

    // Initialize usecase with dependencies
    userUsecase := usecase.NewUserUsecase(userRepo, tokenGen, passwordSvc, logRepo)
    loanUsecase := usecase.NewLoanUsecase(loanRepo)
    logUsecase := usecase.NewLogUsecase(logRepo)

    // Initialize controller with usecase
    userController := controllers.NewUserController(userUsecase)
    loanController := controllers.NewLoanController(loanUsecase)
    logController := controllers.NewLogController(logUsecase)

    Admin := router.Group("/admin")
    Admin.Use(infrastracture.AuthMiddleware(), infrastracture.RoleMiddleware("admin"))
    {
        Admin.GET("/users", userController.GetUsers)
        Admin.DELETE("/users/:id", userController.DeleteUser)

        Admin.GET("/loans", loanController.GetLoans)
        Admin.PATCH("/loans/:id/status", loanController.UpdateLoanStatus)
        Admin.DELETE("/loans/:id", loanController.DeleteLoan)

        Admin.GET("/logs", logController.GetLogs) 
    }
}



// package routers

// import (
// 	"loan-tracker-api/config/db"
// 	"loan-tracker-api/delivery/controllers"
// 	"loan-tracker-api/infrastracture"
// 	"loan-tracker-api/repository"
// 	"loan-tracker-api/usecase"

// 	"github.com/gin-gonic/gin"
// )

// func setUpAdminRoutes(router *gin.Engine) {
// 	// Initialize repository with database collection
// 	userRepo := repository.NewUserRepositoryImpl(db.UserCollection)
// 	loanRepo := repository.NewLoanRepositoryImpl(db.LoanCollection)

// 	// Initialize token generator and password service
// 	tokenGen := infrastracture.NewTokenGenerator()
// 	passwordSvc := infrastracture.NewPasswordService()

// 	// Initialize usecase with dependencies
// 	userUsecase := usecase.NewUserUsecase(userRepo, tokenGen, passwordSvc)
// 	loanUsecase := usecase.NewLoanUsecase(loanRepo)

// 	// Initialize controller with usecase
// 	userController := controllers.NewUserController(userUsecase)
// 	loanController := controllers.NewLoanController(loanUsecase)

// 	Admin := router.Group("/admin")
// 	Admin.Use(infrastracture.AuthMiddleware(), infrastracture.RoleMiddleware("admin"))
// 	{
// 		Admin.GET("/users", userController.GetUsers)
// 		Admin.DELETE("/users/:id", userController.DeleteUser)

// 		Admin.GET("/loans", loanController.GetLoans)
// 		Admin.PATCH("/loans/:id/status", loanController.UpdateLoanStatus)
// 		Admin.DELETE("/loans/:id", loanController.DeleteLoan)
// 	}
// }