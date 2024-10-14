package api

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/dsxg666/web-tool/global"
	"github.com/dsxg666/web-tool/internal/model"
	"github.com/dsxg666/web-tool/pkg/encrypt"
	"github.com/dsxg666/web-tool/pkg/result"
	"github.com/gin-gonic/gin"
)

type User struct{}

func NewUser() User {
	return User{}
}

func (User) GetById(c *gin.Context) {
	var idDTO IdDTO

	if err := c.ShouldBindJSON(&idDTO); err != nil {
		global.Logger.Errorf("err: %v", err)
		c.JSON(http.StatusBadRequest, result.InvalidRequestData)
		return
	}

	userIdDTO := &model.UserIdDTO{Id: idDTO.Id}
	c.JSON(http.StatusOK, result.SuccessWithData(userIdDTO.GetById()))
}

func (User) ModifyEmail(c *gin.Context) {
	var userModifyDTO model.UserModifyDTO

	if err := c.ShouldBindJSON(&userModifyDTO); err != nil {
		global.Logger.Errorf("err: %v", err)
		c.JSON(http.StatusBadRequest, result.InvalidRequestData)
		return
	}

	if userModifyDTO.IsEmailExist() {
		c.JSON(http.StatusOK, result.OperateError("The email address has already exists.", ""))
	} else {
		userModifyDTO.ModifyEmail()
		c.JSON(http.StatusOK, result.Success)
	}
}

func (User) ModifyPassword(c *gin.Context) {
	var userModifyDTO model.UserModifyDTO

	if err := c.ShouldBindJSON(&userModifyDTO); err != nil {
		global.Logger.Errorf("err: %v", err)
		c.JSON(http.StatusBadRequest, result.InvalidRequestData)
		return
	}

	hashPass, err := encrypt.HashPassword(userModifyDTO.Password)
	if err != nil {
		global.Logger.Errorf("err: %v", err)
	}
	userModifyDTO.Password = hashPass
	userModifyDTO.ModifyPassword()
	c.JSON(http.StatusOK, result.Success)
}

func (User) ModifyUsername(c *gin.Context) {
	var userModifyDTO model.UserModifyDTO

	if err := c.ShouldBindJSON(&userModifyDTO); err != nil {
		global.Logger.Errorf("err: %v", err)
		c.JSON(http.StatusBadRequest, result.InvalidRequestData)
		return
	}

	userModifyDTO.ModifyUsername()
	c.JSON(http.StatusOK, result.Success)
}

func (User) UploadAvatar(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file is received"})
		return
	}

	if file.Size > 500*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File too large"})
		return
	}

	ext := filepath.Ext(file.Filename)
	if ext != ".png" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only png files are allowed"})
		return
	}

	wd, err := os.Getwd()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get working directory"})
		return
	}
	userIdAny, _ := c.Get("userId")
	userId, _ := userIdAny.(string)

	userModifyDTO := &model.UserModifyDTO{Id: userId, Avatar: userId + ".png"}
	userModifyDTO.ModifyAvatar()

	filePath := filepath.Join(wd, "storage/avatars", userId+".png")

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save the file"})
		return
	}

	c.JSON(http.StatusOK, result.Success)
}
