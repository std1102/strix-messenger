package router

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"mime/multipart"
	"strconv"
	"strings"
	"strix-server/persistence"
	"strix-server/repository"
	"strix-server/system"
	"time"
)

// File
func uploadFile(c *gin.Context) {
	currentUser := getLoggedInUser(c)
	file, err := c.FormFile("upload")
	if err != nil {
		handleError(c, 500, err)
		return
	}
	fileBuffer, err := file.Open()
	if err != nil {
		handleError(c, 500, err)
		return
	}
	defer func(fileBuffer multipart.File) {
		err := fileBuffer.Close()
		if err != nil {
			system.Logger.Error(err)
		}
	}(fileBuffer)
	fileSize := file.Size
	var fileType string
	if c.PostForm("fileType") != "" {
		fileType = c.GetHeader("fileType")
	}
	fileNameAndExtension := strings.Split(file.Filename, ".")
	if len(fileNameAndExtension) > 1 {
		fileType = fileNameAndExtension[1]
	}
	newId, err := uuid.NewUUID()
	if err != nil {
		handleError(c, 500, err)
		return
	}
	bucketName := system.SystemConfig.Binary.Bucket
	ctx := context.Background()
	_, err = persistence.MinioClient.PutObject(ctx, bucketName, newId.String(), fileBuffer, fileSize, minio.PutObjectOptions{ContentType: fileType})
	if err != nil {
		handleError(c, 500, err)
		return
	}
	fileRepo := repository.NewFileRepository(persistence.DatabaseContext)
	fileInfo := persistence.UploadedFile{
		ID:        newId,
		Type:      fileType,
		Size:      uint64(fileSize),
		CreatedAt: time.Now(),
		OwnerId:   currentUser.ID,
		Owner:     currentUser,
	}
	err = fileRepo.Save(&fileInfo)
	if err != nil {
		handleError(c, 500, err)
		return
	}
	c.JSON(200, gin.H{
		"filePath": newId.String(),
	})
}

func uploadAvatar(c *gin.Context) {
	currentUser := getLoggedInUser(c)
	file, err := c.FormFile("upload")
	if err != nil {
		handleError(c, 500, err)
		return
	}
	fileBuffer, err := file.Open()
	if err != nil {
		handleError(c, 500, err)
		return
	}
	defer func(fileBuffer multipart.File) {
		err := fileBuffer.Close()
		if err != nil {
			system.Logger.Error(err)
		}
	}(fileBuffer)
	fileSize := file.Size
	var fileType string
	if c.PostForm("fileType") != "" {
		fileType = c.GetHeader("fileType")
	}
	fileNameAndExtension := strings.Split(file.Filename, ".")
	if len(fileNameAndExtension) > 1 {
		fileType = fileNameAndExtension[1]
	}
	newId, err := uuid.NewUUID()
	if err != nil {
		handleError(c, 500, err)
		return
	}
	bucketName := system.SystemConfig.Binary.Bucket
	ctx := context.Background()
	_, err = persistence.MinioClient.PutObject(ctx, bucketName, newId.String(), fileBuffer, fileSize, minio.PutObjectOptions{ContentType: fileType})
	if err != nil {
		handleError(c, 500, err)
		return
	}
	fileRepo := repository.NewFileRepository(persistence.DatabaseContext)
	fileInfo := persistence.UploadedFile{
		ID:        newId,
		Type:      fileType,
		Size:      uint64(fileSize),
		CreatedAt: time.Now(),
		OwnerId:   currentUser.ID,
		Owner:     currentUser,
	}
	err = fileRepo.Save(&fileInfo)
	newFileIDString := newId.String()
	currentUser.Avatar = &newFileIDString
	userRepo := repository.NewUserRepository(persistence.DatabaseContext)
	err = userRepo.Save(currentUser)
	if err != nil {
		handleError(c, 500, err)
		return
	}
	if err != nil {
		handleError(c, 500, err)
		return
	}
	c.JSON(200, gin.H{
		"filePath": newId.String(),
	})
}

func getFile(c *gin.Context) {
	fileId := c.Query("fileId")
	fileRepo := repository.NewFileRepository(persistence.DatabaseContext)
	var storedFile persistence.UploadedFile
	err := fileRepo.FindById(fileId, &storedFile)
	if err != nil {
		handleError(c, 500, err)
		return
	}
	bucketName := system.SystemConfig.Binary.Bucket
	ctx := context.Background()
	object, err := persistence.MinioClient.GetObject(ctx, bucketName, fileId, minio.GetObjectOptions{})
	if err != nil {
		if err.Error() != "EOF" {
			handleError(c, 500, err)
			return
		}
	}
	fileData := make([]byte, storedFile.Size)
	_, err = object.Read(fileData)
	if err != nil {
		if err.Error() != "EOF" {
			handleError(c, 500, err)
			return
		}
	}
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename="+fileId+"."+storedFile.Type)
	c.Header("Content-Length", strconv.FormatUint(storedFile.Size, 10))
	c.Data(200, "application/octet-stream", fileData)
}
