package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"strix-server/common"
	"strix-server/persistence"
	"strix-server/repository"
	"strix-server/system"
	"time"
)

// Auth
func register(context *gin.Context) {
	var dto RegisterDto
	err := context.BindJSON(&dto)
	if err != nil {
		handleError(context, 400, fmt.Errorf(err.Error()))
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(common.StringToByte(dto.Password), BCRYPT_COST)
	if err != nil {
		handleError(context, 500, fmt.Errorf(err.Error()))
	}

	hashedPasswordString := common.EncodeToString(hashedPassword)

	var userRepository = repository.NewUserRepository(persistence.DatabaseContext)

	var user persistence.User

	err = userRepository.FindByUserName(dto.Username, &user)
	if err == nil {
		handleError(context, 400, fmt.Errorf("User existed"))
		return
	}

	userId, _ := uuid.NewUUID()

	user = persistence.User{
		ID:                userId,
		Username:          dto.Username,
		Password:          hashedPasswordString,
		AliasName:         dto.AliasName,
		Email:             dto.Email,
		Avatar:            nil,
		IdentityKey:       "",
		PreKeyCreatedTime: nil,
		PreKeys:           nil,
		Devices:           nil,
		CreatedAt:         time.Now(),
	}

	err = userRepository.Save(&user)
	if err != nil {
		handleError(context, 500, fmt.Errorf(err.Error()))
		return
	}

	context.JSON(200, gin.H{
		"user":    user.Username,
		"message": "User created",
	})
}

func login(context *gin.Context) {
	var loginDto LoginDto
	err := context.BindJSON(&loginDto)
	if err != nil {
		handleError(context, 400, fmt.Errorf("Invalid request"))
		return
	}
	if loginDto.LoginType == "" {
		handleError(context, 400, fmt.Errorf("Invalid request"))
		return
	}
	if loginDto.LoginType == "password" {
		// Validate something
		if loginDto.Username == "" || loginDto.Password == "" {
			handleError(context, 400, fmt.Errorf("Invalid request"))
			return
		}
		// Retrieve user
		var userRepository = repository.NewUserRepository(persistence.DatabaseContext)
		var user persistence.User
		err = userRepository.FindByUserName(loginDto.Username, &user)
		if err != nil {
			handleError(context, 401, fmt.Errorf(err.Error()))
			return
		}
		// Check password
		err := bcrypt.CompareHashAndPassword(common.DecodeToByte(user.Password), common.StringToByte(loginDto.Password))
		if err != nil {
			handleError(context, 401, fmt.Errorf(err.Error()))
			return
		}
		// Generate token
		currentTime := time.Now()
		accessToken, err := generateToken(&currentTime, system.SystemConfig.Auth.AccessTokenExpireTime, user.ID.String())
		if err != nil {
			handleError(context, 500, fmt.Errorf(err.Error()))
			return
		}
		// Refresh token
		var refreshToken string
		if loginDto.RememberMe {
			refreshToken, err = generateToken(&currentTime, system.SystemConfig.Auth.RefreshTokenExpireTime, user.ID.String())
			if err != nil {
				handleError(context, 500, fmt.Errorf(err.Error()))
				return
			}
		}
		context.JSON(200, LoginResponseDto{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			LoggedInAt:   common.FormatTime(&currentTime),
		})
		return
	}
	if loginDto.LoginType == "refresh_token" {
		if loginDto.RefreshToken == "" {
			if err != nil {
				handleError(context, 400, fmt.Errorf(err.Error()))
				return
			}
		}
		var claimMap jwt.MapClaims
		_, err := jwt.ParseWithClaims(loginDto.RefreshToken, &claimMap, func(token *jwt.Token) (any, error) {
			return common.DecodeToByte(system.SystemConfig.JwtKey), nil
		})
		if err != nil {
			handleError(context, 401, fmt.Errorf(err.Error()))
			return
		}
		userId := claimMap["userId"].(string)
		expiredAt := claimMap["exp"].(float64)

		currentTime := time.Now()

		if float64(currentTime.UnixMilli()) > expiredAt {
			handleError(context, 401, fmt.Errorf("Unauthorized"))
			return
		}

		var userRepository = repository.NewUserRepository(persistence.DatabaseContext)
		var user persistence.User
		err = userRepository.FindById(userId, &user)
		if err != nil {
			handleError(context, 401, fmt.Errorf(err.Error()))
			return
		}
		accessToken, err := generateToken(&currentTime, system.SystemConfig.Auth.AccessTokenExpireTime, user.ID.String())
		if err != nil {
			handleError(context, 500, fmt.Errorf(err.Error()))
			return
		}
		refreshToken, err := generateToken(&currentTime, system.SystemConfig.Auth.RefreshTokenExpireTime, user.ID.String())
		if err != nil {
			handleError(context, 500, fmt.Errorf(err.Error()))
			return
		}
		context.JSON(200, LoginResponseDto{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			LoggedInAt:   common.FormatTime(&currentTime),
		})
		return
	}
	handleError(context, 401, fmt.Errorf(err.Error()))
	return
}
