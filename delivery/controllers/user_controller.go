package controllers

import (
	"loan-tracker-api/domain"
	"loan-tracker-api/infrastracture"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	UserUsecase domain.UserUsecase
}

func NewUserController(userUsecase domain.UserUsecase) *UserController {
	return &UserController{UserUsecase: userUsecase}
}

func (u *UserController) Register(c *gin.Context) {
	var user domain.User
	c.BindJSON(&user)

	err := u.UserUsecase.Register(user)
	if err.Message != "" {
		c.JSON(err.StatusCode, err)
		return
	}

	c.JSON(200, gin.H{
		"code":    200,
		"message": "Your account has been created successfully see your email to activate your account",
	})	
}

func (u *UserController) ActivateAccount(c *gin.Context) {
	Email := c.Query("email")
	Token := c.Query("token")

	err := u.UserUsecase.AccountActivation(Email, Token)
	if err.Message != "" {
		c.JSON(err.StatusCode, err)
		return
	}

	c.JSON(200, gin.H{
		"code":    200,
		"message": "Your account has been activated successfully",
	})
	
}



func (u *UserController) Login(c *gin.Context) {
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ipAddress := c.ClientIP()
	userAgent := c.Request.UserAgent()
	deviceFingerprint := infrastracture.GenerateDeviceFingerprint(ipAddress, userAgent)

	response, err := u.UserUsecase.Login(&user, deviceFingerprint)
	if err.Message != "" {
		c.JSON(err.StatusCode, err)
		return
	}

	c.JSON(200, gin.H{
		"code":    200,
		"message": "Login successful",
		"data":    response,
	})
	
}

func (u *UserController) RefreshToken(c *gin.Context) {
	var refreshRequest domain.RefreshTokenRequest

	if err := c.ShouldBindJSON(&refreshRequest); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	}
	ipAddress := c.ClientIP()
	userAgent := c.Request.UserAgent()
	deviceFingerprint := infrastracture.GenerateDeviceFingerprint(ipAddress, userAgent)

	refreshResponse, uerr := u.UserUsecase.RefreshToken(refreshRequest.UserID, deviceFingerprint, refreshRequest.Token)
	if uerr.Message != "" {
		c.JSON(uerr.StatusCode, gin.H{
			"code":uerr.StatusCode,
			"error": uerr.Message,
		})
		return
	}

	c.JSON(200, gin.H{
		"code":200,
		"tokens": refreshResponse,
	})
}


func (u *UserController)GetMyProfile(c *gin.Context) {
	userID := c.GetString("user_id")
	user, err := u.UserUsecase.GetMyProfile(userID)

	if err.Message != "" {
		c.JSON(err.StatusCode, err)
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"data": user,
	})
}

func (u *UserController) GetUsers(c *gin.Context) {
	byName := c.Query("name")
	limit := c.Query("limit") 
	page := c.Query("page")
	if limit == "" {
		limit = "10"

	}
	if page == "" {
		page = "1"
	}




	users, err := u.UserUsecase.GetUsers( byName, limit, page)
	if err.Message != "" {
		c.JSON(err.StatusCode, err)
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"message": "Users fetched successfully",
		"data": users,

		"quantity": strconv.Itoa(len(users)) + "/" + limit,
		"current_page": page,
	})
}


func (u *UserController) DeleteUser(c *gin.Context) {
	userID := c.Param("id")
	user, err := u.UserUsecase.DeleteUser(userID)
	if err.Message != "" {
		c.JSON(err.StatusCode, err)
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"message": "User deleted successfully",
		"data": user,
	})
}



func (uc *UserController) SendPasswordResetLink(c *gin.Context) {
	var req domain.ResetPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	uerr := uc.UserUsecase.SendPasswordResetLink(req.Email)
	if uerr.Message != "" {
		c.JSON(uerr.StatusCode, gin.H{
			"code": uerr.StatusCode,
			"error": uerr.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "Password reset link sent"})
}

func (uc *UserController) ResetPassword(c *gin.Context) {
	var req struct {
		Password string `json:"password"`
	}
	token := c.Param("token")
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	uerr := uc.UserUsecase.ResetPassword(token, req.Password)
	if uerr.Message != "" {
		c.JSON(uerr.StatusCode, gin.H{
			"code": uerr.StatusCode,
			
			"error": uerr.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "Password has been reset"})
}
