package api

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dsxg666/web-tool/global"
	"github.com/dsxg666/web-tool/internal/model"
	"github.com/dsxg666/web-tool/pkg/email"
	"github.com/dsxg666/web-tool/pkg/encrypt"
	"github.com/dsxg666/web-tool/pkg/jwt"
	"github.com/dsxg666/web-tool/pkg/result"
	"github.com/dsxg666/web-tool/pkg/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Base struct{}

func NewBase() Base {
	return Base{}
}

type TokenVerify struct {
	Token string `json:"token"`
}

func (Base) LoginByPassword(c *gin.Context) {
	var userLoginByPasswordDTO model.UserLoginByPasswordDTO

	if err := c.ShouldBindJSON(&userLoginByPasswordDTO); err != nil {
		global.Logger.Errorf("err: %v", err)
		c.JSON(http.StatusBadRequest, result.InvalidRequestData)
		return
	}

	if userLoginByPasswordDTO.IsEmailExist() {
		if userLoginByPasswordDTO.IsCorrectPassword() {
			user := userLoginByPasswordDTO.GetUser()
			claims := jwt.NewClaims(user.Id, user.Username, user.Path)
			token, err := jwt.NewJwtToken(claims)
			if err != nil {
				global.Logger.Error("Error creating token: ", err)
				c.JSON(http.StatusUnauthorized, result.InvalidJwtToken)
			} else {
				ipAddress := c.ClientIP()
				dua := &model.Dau{UserId: user.Id, UserIp: ipAddress}
				dua.Add()
				c.JSON(http.StatusOK, result.SuccessWithMessage("Login successfully.", gin.H{
					"token": token,
				}))
			}
		} else {
			c.JSON(http.StatusOK, result.OperateError("Wrong password", ""))
		}
	} else {
		c.JSON(http.StatusOK, result.OperateError("The email address could not be recognized.", ""))
	}
}

func (Base) LoginByCode(c *gin.Context) {
	var userLoginByCodeDTO model.UserLoginByCodeDTO

	if err := c.ShouldBindJSON(&userLoginByCodeDTO); err != nil {
		global.Logger.Errorf("err: %v", err)
		c.JSON(http.StatusBadRequest, result.InvalidRequestData)
		return
	}

	if userLoginByCodeDTO.IsEmailExist() {
		codeDTO := &model.CodeDTO{Email: userLoginByCodeDTO.Email, Code: userLoginByCodeDTO.Code}
		if codeDTO.IsValidVerificationCode() {
			user := userLoginByCodeDTO.GetUser()
			claims := jwt.NewClaims(user.Id, user.Username, user.Path)
			token, err := jwt.NewJwtToken(claims)
			if err != nil {
				global.Logger.Error("Error creating token: ", err)
				c.JSON(http.StatusUnauthorized, result.InvalidJwtToken)
			} else {
				ipAddress := c.ClientIP()
				dua := &model.Dau{UserId: user.Id, UserIp: ipAddress}
				dua.Add()
				c.JSON(http.StatusOK, result.SuccessWithMessage("Login successfully.", gin.H{
					"token": token,
				}))
			}
		} else {
			c.JSON(http.StatusOK, result.OperateError("Invalid verification code.", ""))
		}
	} else {
		c.JSON(http.StatusOK, result.OperateError("The email address could not be recognized.", ""))
	}
}

func (Base) Register(c *gin.Context) {
	var userRegisterDTO model.UserRegisterDTO

	if err := c.ShouldBindJSON(&userRegisterDTO); err != nil {
		global.Logger.Errorf("err: %v", err)
		c.JSON(http.StatusBadRequest, result.InvalidRequestData)
		return
	}

	if userRegisterDTO.IsEmailExist() {
		c.JSON(http.StatusOK, result.OperateError("This email address has already been registered.", ""))
	} else {
		codeDTO := &model.CodeDTO{Email: userRegisterDTO.Email, Code: userRegisterDTO.Code}
		if codeDTO.IsValidRegisterCode() {
			// Set encrypt password
			hashPass, err := encrypt.HashPassword(userRegisterDTO.Password)
			if err != nil {
				global.Logger.Errorf("err: %v", err)
			}
			userRegisterDTO.Password = hashPass

			// Add random path
			userRegisterDTO.Path = util.RandomString(6)

			// Add User
			lastId := userRegisterDTO.Add()

			// Dir work
			wd, err := os.Getwd()
			if err != nil {
				global.Logger.Errorf("Failed to get working directory err: %v", err)
			}

			folderPath := filepath.Join(wd, "storage/todolist/"+lastId)
			folderPath2 := filepath.Join(wd, "storage/blog/"+lastId)
			err = os.MkdirAll(folderPath, 0755)
			err2 := os.MkdirAll(folderPath2, 0755)
			if err != nil || err2 != nil {
				global.Logger.Errorf("Error creating directories: %v", err)
			}

			// Add to WorldGroup
			groupMembers := &model.GroupMembers{GroupId: "1", UserId: lastId, Status: "0"}
			groupMembers.Add()

			c.JSON(http.StatusOK, result.SuccessWithMessage("Registered successfully.", ""))
		} else {
			c.JSON(http.StatusOK, result.OperateError("Invalid register code.", ""))
		}
	}
}

func (Base) GetRegisterCode(c *gin.Context) {
	e := c.Query("email")
	code := &model.CodeDTO{Email: e, Code: util.GetSixRandomCode()}
	userRegisterDTO := &model.UserRegisterDTO{Email: e}
	if userRegisterDTO.IsEmailExist() {
		c.JSON(http.StatusOK, result.OperateError("This email address has already been registered.", ""))
	} else {
		if code.IsTooQuickRegister() {
			c.JSON(http.StatusOK, result.OperateError("Your code is sent too often, please try again later.", ""))
		} else {
			code.AddRegisterCode()
			remind := "Please use the register code as soon as possible, the register code will be invalid after 5 minutes."
			email.SendEmail(code.Email, "Register code", fmt.Sprintf("Your register code is: %s. %s", code.Code, remind))
			c.JSON(http.StatusOK, result.SuccessWithMessage("The register code is sent successfully.", ""))
		}
	}
}

func (Base) GetVerificationCode(c *gin.Context) {
	userLoginByCodeDTO := &model.UserLoginByCodeDTO{Email: c.Query("email")}
	if userLoginByCodeDTO.IsEmailExist() {
		code := &model.CodeDTO{Email: c.Query("email"), Code: util.GetSixRandomCode()}
		if code.IsTooQuickVerification() {
			c.JSON(http.StatusOK, result.OperateError("Your code is sent too often, please try again later.", ""))
		} else {
			code.AddVerificationCode()
			remind := "Please use the verification code as soon as possible, the verification code will be invalid after 5 minutes."
			email.SendEmail(code.Email, "Verification code", fmt.Sprintf("Your verification code is: %s. %s", code.Code, remind))
			c.JSON(http.StatusOK, result.SuccessWithMessage("The verification code is sent successfully.", ""))
		}
	} else {
		c.JSON(http.StatusOK, result.OperateError("The email address could not be recognized.", ""))
	}
}

func (Base) VerifyToken(c *gin.Context) {
	var tokenVerify TokenVerify

	if err := c.ShouldBindJSON(&tokenVerify); err != nil {
		global.Logger.Errorf("err: %v", err)
		c.JSON(http.StatusOK, result.InvalidRequestData)
		return
	}

	claim, ok, err := jwt.ParseJwtToken(tokenVerify.Token)
	if ok {
		c.Set("userId", claim.UserId)
		c.JSON(http.StatusOK, result.SuccessWithMessage("Token is valid.", ""))
	} else {
		global.Logger.Error("parse token err: ", err)
		c.JSON(http.StatusOK, result.InvalidJwtToken)
	}
}

func (Base) GetBlogImg(c *gin.Context) {
	id := c.Query("id")

	filename := c.Param("filename")

	wd, err := os.Getwd()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get working directory"})
		return
	}

	filePath := filepath.Join(wd, "storage/blog/"+id, filename)

	if fileInfo, err := os.Stat(filePath); err != nil || fileInfo.IsDir() {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		global.Logger.Errorf("FileNotFound: %v", err)
		return
	}

	c.File(filePath)
}

func (Base) GetTodolistImg(c *gin.Context) {
	id := c.Query("id")

	filename := c.Param("filename")

	wd, err := os.Getwd()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get working directory"})
		return
	}

	filePath := filepath.Join(wd, "storage/todolist/"+id, filename)

	if fileInfo, err := os.Stat(filePath); err != nil || fileInfo.IsDir() {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		global.Logger.Errorf("FileNotFound: %v", err)
		return
	}

	c.File(filePath)
}

func (Base) GetGroupAvatar(c *gin.Context) {
	filename := c.Param("filename")

	wd, err := os.Getwd()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get working directory"})
		return
	}

	filePath := filepath.Join(wd, "storage/chatroom/avatars", filename)

	if fileInfo, err := os.Stat(filePath); err != nil || fileInfo.IsDir() {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		global.Logger.Errorf("FileNotFound: %v", err)
		return
	}

	c.File(filePath)
}

func (Base) GetAvatar(c *gin.Context) {
	userId := c.Param("userId")

	userAvatarDTO := &model.UserAvatarDTO{Id: userId}

	wd, err := os.Getwd()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get working directory"})
		return
	}

	filePath := filepath.Join(wd, "storage/avatars", userAvatarDTO.GetUser().Avatar)

	if fileInfo, err := os.Stat(filePath); err != nil || fileInfo.IsDir() {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		global.Logger.Errorf("FileNotFound: %v", err)
		return
	}

	c.File(filePath)
}
