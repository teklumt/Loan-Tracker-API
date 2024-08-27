package routers

import (
	"loan-tracker-api/config/db"
	"loan-tracker-api/delivery/controllers"
	"loan-tracker-api/infrastracture"
	"loan-tracker-api/repository"
	"loan-tracker-api/usecase"

	"github.com/gin-gonic/gin"
)

func setUpAuthRoutes(router *gin.Engine) {
	// Initialize repository with database collection
	userRepo := repository.NewUserRepositoryImpl(db.UserCollection)
	logRepo := repository.NewLogRepositoryImpl(db.LogCollection)

	// Initialize token generator and password service
	tokenGen := infrastracture.NewTokenGenerator() 
	passwordSvc := infrastracture.NewPasswordService()

	// Initialize usecase with dependencies
	userUsecase := usecase.NewUserUsecase(userRepo, tokenGen, passwordSvc, logRepo)

	// Initialize controller with usecase
	authController := controllers.NewUserController(userUsecase)
	

	auth := router.Group("/users")
	{
		auth.POST("/register", authController.Register)
		auth.GET("/verify-email", authController.ActivateAccount)
		auth.POST("/login", authController.Login)
		auth.POST("/token/refresh", authController.RefreshToken)
		


		auth.POST("/reset-password", authController.SendPasswordResetLink)
		auth.POST("/reset-password/:token", authController.ResetPassword)
		
	}
	
}
