package usecase

import (
	"loan-tracker-api/domain"
	"loan-tracker-api/infrastracture"
	"log"
	"time"
)

type UserUsecase struct {
	UserRepo       domain.UserRepository
	TokenGen       domain.TokenGenerator
	PasswordSvc    domain.PasswordService
	LogRepo        domain.LogRepository
}

func NewUserUsecase(userRepo domain.UserRepository, tokenGen domain.TokenGenerator, passwordSvc domain.PasswordService,LogRepo domain.LogRepository ) domain.UserUsecase {
	return &UserUsecase{
		UserRepo:    userRepo,
		TokenGen:    tokenGen,
		PasswordSvc: passwordSvc,
		LogRepo:     LogRepo,


	}
}

func (u *UserUsecase) Register(user domain.User) domain.ErrorResponse {
	if user.Username == "" || user.Email == "" || user.Password == "" {
		return domain.ErrorResponse{StatusCode: 400, Message: "All fields are required"}
	}

	if !infrastracture.IsValidEmail(user.Email) {
		return domain.ErrorResponse{StatusCode: 400, Message: "Invalid email address"}
	}

	if !infrastracture.IsValidPassword(user.Password) {
		return domain.ErrorResponse{StatusCode: 400, Message: "Password must be at least 8 characters long"}
	}

	_, err := u.UserRepo.GetUserByUsernameOrEmail(user.Username, user.Email)
	if err == nil {
		return domain.ErrorResponse{StatusCode: 400, Message: "User already exists"}
	}

	user.Role = "user"

	// Hash password
	hashedPassword, err := u.PasswordSvc.HashPassword(user.Password)
	if err != nil {
		return domain.ErrorResponse{StatusCode: 500, Message: "Failed to hash password"}
	}

	token, err := infrastracture.GenerateActivationToken()
	if err != nil {
		return domain.ErrorResponse{StatusCode: 500, Message: "Failed to generate activation token"}
	}

	user.Password = hashedPassword
	user.ActivationToken = token
	user.TokenCreatedAt = time.Now()

	// Create user account in the database
	err = u.UserRepo.Register(user)
	if err != nil {
		return domain.ErrorResponse{StatusCode: 500, Message: "Failed to create user account"}
	}

	// Send activation email or link to the user
	err = infrastracture.SendActivationEmail(user.Email, token)
	if err != nil {
		return domain.ErrorResponse{StatusCode: 500, Message: "Failed to send activation email"}
	}

	return domain.ErrorResponse{}
}


func (u *UserUsecase) AccountActivation(email, token string) domain.ErrorResponse {
	user, err := u.UserRepo.GetUserByEmail(email)
	if err != nil {
		return domain.ErrorResponse{StatusCode: 400, Message: "User not found"}
	}

	if user.ActivationToken != token {
		return domain.ErrorResponse{StatusCode: 400, Message: "Invalid activation token"}
	}

	if time.Since(user.TokenCreatedAt).Minutes() > 30 {
		return domain.ErrorResponse{StatusCode: 400, Message: "Activation token has expired"}
	}

	err = u.UserRepo.AccountActivation(email)
	if err != nil {
		return domain.ErrorResponse{StatusCode: 500, Message: "Failed to activate account"}
	}

	return domain.ErrorResponse{}
}

func (u *UserUsecase) Login(user *domain.User, deviceID string) (domain.LogInResponse, domain.ErrorResponse) {
    if u.UserRepo == nil || u.PasswordSvc == nil || u.TokenGen == nil {
        log.Fatal("Necessary services are nil")
        return domain.LogInResponse{}, domain.ErrorResponse{StatusCode: 500, Message: "Internal server error"}
    }

    existingUser, err := u.UserRepo.Login(user)
    if err != nil {
        // Log failed login attempt
        u.LogRepo.CreateLog(domain.SystemLog{
            Timestamp: time.Now().String(),
            Event:     "Login Attempt",
            Details:   "Failed login attempt for user " + user.Email,
        })
        return domain.LogInResponse{}, domain.ErrorResponse{StatusCode: 400, Message: "Invalid credentials"}
    }

    if !u.PasswordSvc.CheckPasswordHash(user.Password, existingUser.Password) {
        // Log failed login attempt
        u.LogRepo.CreateLog(domain.SystemLog{
            Timestamp: time.Now().String(),
            Event:     "Login Attempt",
            Details:   "Failed login attempt for user " + user.Email,
        })
        return domain.LogInResponse{}, domain.ErrorResponse{StatusCode: 400, Message: "Invalid credentials"}
    }

    // Log successful login attempt
    u.LogRepo.CreateLog(domain.SystemLog{
        Timestamp: time.Now().String(),
        Event:     "Login Attempt",
        Details:   "User " + user.Email + " logged in successfully",
    })

    if !existingUser.IsActive {
        token, err := infrastracture.GenerateActivationToken()
        if err != nil {
            return domain.LogInResponse{}, domain.ErrorResponse{StatusCode: 500, Message: "Failed to generate activation token"}
        }

        existingUser, err := u.UserRepo.GetUserByEmail(user.Email)
        if err != nil {
            return domain.LogInResponse{}, domain.ErrorResponse{StatusCode: 400, Message: "User not found"}
        }

        existingUser.ActivationToken = token
        existingUser.TokenCreatedAt = time.Now()

        err = u.UserRepo.UpdateUser(&existingUser)
        if err != nil {
            return domain.LogInResponse{}, domain.ErrorResponse{StatusCode: 500, Message: "Failed to create user account"}
        }

        // Send activation email or link to the user
        err = infrastracture.SendActivationEmail(user.Email, token)
        if err != nil {
            return domain.LogInResponse{}, domain.ErrorResponse{StatusCode: 500, Message: "Failed to send activation email"}
        }

        return domain.LogInResponse{}, domain.ErrorResponse{StatusCode: 400, Message: "Account is not activated yet. Please check your email for activation link."}
    }

    refreshToken, err := u.TokenGen.GenerateRefreshToken(*existingUser)
    if err != nil {
        return domain.LogInResponse{}, domain.ErrorResponse{StatusCode: 500, Message: "Failed to generate refresh token"}
    }

    newRefreshToken := domain.RefreshToken{
        Token:     refreshToken,
        DeviceID:  deviceID,
        CreatedAt: time.Now(),
    }

    for i, rt := range existingUser.RefreshTokens {
        if rt.DeviceID == deviceID {
            existingUser.RefreshTokens = append(existingUser.RefreshTokens[:i], existingUser.RefreshTokens[i+1:]...)
            break
        }
    }

    existingUser.RefreshTokens = append(existingUser.RefreshTokens, newRefreshToken)

    err = u.UserRepo.UpdateUser(existingUser)
    if err != nil {
        return domain.LogInResponse{}, domain.ErrorResponse{StatusCode: 500, Message: "Failed to update user"}
    }

    accessToken, err := u.TokenGen.GenerateToken(*existingUser)
    if err != nil {
        return domain.LogInResponse{}, domain.ErrorResponse{StatusCode: 500, Message: "Failed to generate access token"}
    }

    return domain.LogInResponse{
        AccessToken:  accessToken,
        RefreshToken: newRefreshToken.Token,
        User: domain.ReturnUser{
            ID:        existingUser.ID,
            Username:  existingUser.Username,
            Email:     existingUser.Email,
            Role:      existingUser.Role,
            CreatedAt: existingUser.CreatedAt,
        },
    }, domain.ErrorResponse{}
}

func (u *UserUsecase) RefreshToken(userID, deviceID, token string) (domain.RefreshTokenResponse, domain.ErrorResponse) {
	user, err := u.UserRepo.GetUserByID(userID)
	if err != nil {
		return domain.RefreshTokenResponse{}, domain.ErrorResponse{StatusCode: 400, Message: "Invalid refresh token"}
	}

	for i, rt := range user.RefreshTokens {
		if rt.Token == token && rt.DeviceID == deviceID {
			user.RefreshTokens = append(user.RefreshTokens[:i], user.RefreshTokens[i+1:]...)

			refreshToken, err := u.TokenGen.GenerateRefreshToken(user)
			if err != nil {
				return domain.RefreshTokenResponse{}, domain.ErrorResponse{ }
			}

			newRefreshToken := domain.RefreshToken{
				Token:     refreshToken,
				DeviceID:  deviceID,
				CreatedAt: time.Now(),
			}

			user.RefreshTokens = append(user.RefreshTokens, newRefreshToken)
			err = u.UserRepo.UpdateUser(&user)
			if err != nil {
				return domain.RefreshTokenResponse{}, domain.ErrorResponse{ StatusCode: 500 ,Message :"Failed to update user"}
			}

			accessToken, err := u.TokenGen.GenerateToken(user)
			if err != nil {
				return domain.RefreshTokenResponse{}, domain.ErrorResponse{StatusCode: 500, Message: "Failed to Generate Token"}
			}

			return domain.RefreshTokenResponse{
				AccessToken:  accessToken,
				
			}, domain.ErrorResponse{}
		}
	}

	return domain.RefreshTokenResponse{}, domain.ErrorResponse{ StatusCode: 500, Message: "You are not logged in."}
}


func(u *UserUsecase)GetMyProfile (userID string) (domain.ReturnUser, domain.ErrorResponse) {
	user, err := u.UserRepo.GetUserByID(userID)
	if err != nil {
		return domain.ReturnUser{}, domain.ErrorResponse{StatusCode: 400, Message: "User not found"}
	}

	return domain.ReturnUser{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
		CreatedAt: user.CreatedAt,
	}, domain.ErrorResponse{}
}


func (u *UserUsecase) GetUsers(byName, limit , page string) ([]domain.ReturnUser, domain.ErrorResponse) {
	users, err := u.UserRepo.GetUsers( byName, limit , page)
	if err != nil {
		return []domain.ReturnUser{}, domain.ErrorResponse{StatusCode: 500, Message: "Failed to get users"}
	}

	var returnUsers []domain.ReturnUser
	for _, user := range users {
		returnUsers = append(returnUsers, domain.ReturnUser{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Role:     user.Role,
			CreatedAt: user.CreatedAt,
		})
	}

	return returnUsers, domain.ErrorResponse{}
}


func (u *UserUsecase) DeleteUser(id string) (domain.ReturnUser ,domain.ErrorResponse) {
	user, err := u.UserRepo.GetUserByID(id)
	if err != nil {
		return domain.ReturnUser{},domain.ErrorResponse{StatusCode: 400, Message: "User not found"}
	}

	_ ,err = u.UserRepo.DeleteUser(id)
	if err != nil {
		return domain.ReturnUser{},domain.ErrorResponse{StatusCode: 500, Message: "Failed to delete user"}
	}

	return domain.ReturnUser{
		ID: 	 user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
		CreatedAt: user.CreatedAt,
	}, domain.ErrorResponse{}
}


func (u *UserUsecase) SendPasswordResetLink(email string) domain.ErrorResponse{
	user, err := u.UserRepo.GetUserByEmail(email)
	if err != nil {
		return domain.ErrorResponse{StatusCode: 400, Message: "User not found"}
	}

	resetToken, err := infrastracture.GenerateActivationToken()
	if err != nil {
		return domain.ErrorResponse{StatusCode: 500, Message: "Failed to generate reset token"}
	}
	user.PasswordResetToken = resetToken
	user.TokenCreatedAt = time.Now()

	err = u.UserRepo.UpdateUser(&user)
	if err != nil {
		return  domain.ErrorResponse{StatusCode: 500, Message: "Failed to update user"}
	}

	err = infrastracture.SendResetLink(user.Email, resetToken)
	if err != nil {
		return domain.ErrorResponse{StatusCode: 500, Message: "Failed to send reset link"}
	}

	return domain.ErrorResponse{}
}



func (u *UserUsecase) ResetPassword(token, newPassword string) domain.ErrorResponse {
	user, err := u.UserRepo.GetUserByResetToken(token)
	if err != nil {

		u.LogRepo.CreateLog(domain.SystemLog{
			Timestamp: time.Now().String(),
			Event:     "Password Reset",
			Details: "Failed password reset attempt",
		})



		return  domain.ErrorResponse{StatusCode: 400, Message: "Invalid reset token"}
	}

	hashedPassword, err := u.PasswordSvc.HashPassword(newPassword)
	if err != nil {

		u.LogRepo.CreateLog(domain.SystemLog{
			Timestamp: time.Now().String(),
			Event:     "Password Reset",
			Details:  user.Email + "Failed password reset attempt",
		})
		return domain.ErrorResponse{StatusCode: 500, Message: "Failed to hash password"}
	}

	user.Password = hashedPassword
	user.PasswordResetToken = ""
	user.TokenCreatedAt = time.Time{}

	err = u.UserRepo.UpdateUser(&user)
	if err != nil {

		return domain.ErrorResponse{StatusCode: 500, Message: "Failed to update user"}
	}

	u.LogRepo.CreateLog(domain.SystemLog{
		Timestamp: time.Now().String(),
		Event:     "Password Reset",
		Details: user.Email +  "Password reset successful",
	})

	return domain.ErrorResponse{}
}