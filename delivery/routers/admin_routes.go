package routers

import (
	"loan-tracker-api/config/db"
	"loan-tracker-api/delivery/controllers"
	"loan-tracker-api/infrastracture"
	"loan-tracker-api/repository"
	"loan-tracker-api/usecase"

	"github.com/gin-gonic/gin"
)

func setUpAuthAdminRoutes(router *gin.Engine) {
	// Initialize repository with database collection
	userRepo := repository.NewUserRepositoryImpl(db.UserCollection)

	// Initialize token generator and password service
	tokenGen := infrastracture.NewTokenGenerator()
	passwordSvc := infrastracture.NewPasswordService()

	// Initialize usecase with dependencies
	userUsecase := usecase.NewUserUsecase(userRepo, tokenGen, passwordSvc)

	// Initialize controller with usecase
	userController := controllers.NewUserController(userUsecase)

	Admin := router.Group("/admin")
	Admin.Use(infrastracture.AuthMiddleware(), infrastracture.RoleMiddleware("admin"))
	{
		Admin.GET("/users", userController.GetUsers)
		Admin.DELETE("/users/:id", userController.DeleteUser)
	}
}