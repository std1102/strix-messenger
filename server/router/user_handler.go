package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strix-server/common"
	"strix-server/persistence"
	"strix-server/repository"
	"time"
)

// User
func uploadKey(context *gin.Context) {
	var externalKeyBundle ExternalKeyBundleDto
	err := context.BindJSON(&externalKeyBundle)
	if err != nil {
		handleError(context, 400, fmt.Errorf(err.Error()))
	}
	user := getLoggedInUser(context)
	user.IdentityKey = externalKeyBundle.IdentityKey
	user.PreKey = externalKeyBundle.PreKey
	user.PreKeySignature = externalKeyBundle.PreKeySig

	userRepository := repository.NewUserRepository(persistence.DatabaseContext)
	err = userRepository.Save(user)
	if err != nil {
		handleError(context, 500, fmt.Errorf(err.Error()))
		return
	}

	if externalKeyBundle.OneTimeKeyId != "" {
		userOneTimeKey := persistence.OneTimeKey{
			ID:           common.GetUUIDFromString(externalKeyBundle.OneTimeKeyId),
			UserId:       user.ID,
			Key:          externalKeyBundle.OneTimeKey,
			KeySignature: externalKeyBundle.OneTimeKeySig,
			CreatedAt:    time.Now(),
			Owner:        user,
		}
		oneTimeKeyRepo := repository.NewOneTimeKeyRepository(persistence.DatabaseContext)
		err = oneTimeKeyRepo.Save(&userOneTimeKey)
		if err != nil {
			handleError(context, 500, fmt.Errorf(err.Error()))
			return
		}
	}

	context.JSON(200, gin.H{
		"user":    user.Username,
		"message": "Key uploaded",
	})
}

func getExternalKeyBundle(context *gin.Context) {
	username := context.Param("userName")
	currentUser := getLoggedInUser(context)
	userRepository := repository.NewUserRepository(persistence.DatabaseContext)
	var otherUser persistence.User
	err := userRepository.FindByUserName(username, &otherUser)
	if err != nil {
		handleError(context, 400, fmt.Errorf(err.Error()))
		return
	}
	if currentUser.ID.String() == otherUser.ID.String() {
		handleError(context, 400, fmt.Errorf("Invalid userId"))
		return
	}

	var lastedOneTimeKey persistence.OneTimeKey
	for _, element := range otherUser.OneTimePreKeys {
		lastedOneTimeKey = *element
	}

	result := ExternalKeyBundleDto{
		IdentityKey:   otherUser.IdentityKey,
		PreKey:        otherUser.PreKey,
		PreKeySig:     otherUser.PreKeySignature,
		OneTimeKeyId:  lastedOneTimeKey.ID.String(),
		OneTimeKey:    lastedOneTimeKey.Key,
		OneTimeKeySig: lastedOneTimeKey.KeySignature,
	}

	context.JSON(200, result)
}

func getUserInfo(context *gin.Context) {
	user := getLoggedInUser(context)
	var avt string
	if user.Avatar != nil {
		avt = *user.Avatar
	} else {
		avt = ""
	}
	context.JSON(200, UserDto{
		Id:        user.ID.String(),
		UserName:  user.Username,
		AliasName: user.AliasName,
		Avatar:    avt,
	})
}

func searchUser(context *gin.Context) {
	userKeyWord := context.Query("keyWord")
	var userList []persistence.User
	userRepository := repository.NewUserRepository(persistence.DatabaseContext)
	err := userRepository.SearchUserName(userKeyWord, &userList)
	if err != nil {
		if err.Error() == "empty slice found" {
			context.JSON(200, make([]UserDto, 0))
			return
		}
		handleError(context, 500, fmt.Errorf(err.Error()))
		return
	}
	var result []UserDto
	for i := range userList {
		currentUser := userList[i]
		var avt string
		if currentUser.Avatar != nil {
			avt = *currentUser.Avatar
		} else {
			avt = ""
		}
		userDto := UserDto{
			Id:        currentUser.ID.String(),
			UserName:  currentUser.Username,
			AliasName: currentUser.AliasName,
			Avatar:    avt,
		}
		result = append(result, userDto)
	}
	if result == nil {
		result = make([]UserDto, 0)
	}
	context.JSON(200, result)
}

func getUserInfos(context *gin.Context) {
	var userData UserRequestDto
	err := context.BindJSON(&userData)
	if err != nil {
		handleError(context, 500, fmt.Errorf(err.Error()))
		return
	}
	var userList []persistence.User
	userRepository := repository.NewUserRepository(persistence.DatabaseContext)
	err = userRepository.FindAllByUserNames(userData.UserNames, &userList)
	if err != nil {
		if err.Error() == "empty slice found" {
			context.JSON(200, make([]UserDto, 0))
			return
		}
		handleError(context, 500, fmt.Errorf(err.Error()))
		return
	}
	var result []UserDto
	for i := range userList {
		currentUser := userList[i]
		var avt string
		if currentUser.Avatar != nil {
			avt = *currentUser.Avatar
		} else {
			avt = ""
		}
		result = append(result, UserDto{
			Id:        currentUser.ID.String(),
			UserName:  currentUser.AliasName,
			AliasName: currentUser.AliasName,
			Avatar:    avt,
		})
	}
	if result == nil {
		result = make([]UserDto, 0)
	}
	context.JSON(200, result)
}
