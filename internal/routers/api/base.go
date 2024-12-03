package api

import (
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
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

var (
	oauthConfig = &oauth2.Config{
		ClientID:     "Ov23liKsjwU6r7hmHYar",                     // 替换为你的 GitHub Client ID
		ClientSecret: "111bf2c77917b4b0afbcb14a9989ce4e5ca4e5c7", // 替换为你的 GitHub Client Secret
		RedirectURL:  "http://localhost:8000/api/base/callback",
		Scopes:       []string{"read:user"},
		Endpoint:     github.Endpoint,
	}
	state = "random_state_string" // 防止 CSRF，可以随机生成并验证
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

// LoginByGithub 跳转到 GitHub 授权页面
func (Base) LoginByGithub(c *gin.Context) {
	url := oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
	c.Redirect(http.StatusFound, url)
}

// Callback 处理 GitHub 回调
func (Base) Callback(c *gin.Context) {
	// 验证 state 是否匹配
	if c.Query("state") != state {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid state"})
		return
	}

	// 使用授权码获取访问令牌
	code := c.Query("code")
	token, err := oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		global.Logger.Errorf("Exchange error: %v", err)
		c.Redirect(http.StatusFound, global.ServerSetting.FrontendHost)
		return
	}

	// 使用访问令牌获取用户信息
	client := oauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		global.Logger.Errorf("Get user info error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user info"})
		return
	}
	defer resp.Body.Close()

	var userInfo map[string]interface{}
	if err = json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		global.Logger.Errorf("Decode error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse user info"})
		return
	}
	emailI := userInfo["email"]
	emailStr := emailI.(string)
	usernameI := userInfo["name"]
	username := usernameI.(string)
	pathI := userInfo["login"]
	path := pathI.(string)
	// 判断是否存在，登陆或者注册
	userRegisterDTO := &model.UserRegisterDTO{Username: username, Email: emailStr, Password: util.RandomString(8), Path: path}
	if userRegisterDTO.IsEmailExist() {
		// 返回 Token 进行登陆
		user := userRegisterDTO.GetUser()
		claims := jwt.NewClaims(user.Id, user.Username, user.Path)
		myToken, err := jwt.NewJwtToken(claims)
		if err != nil {
			global.Logger.Error("Error creating token: ", err)
			c.JSON(http.StatusUnauthorized, result.InvalidJwtToken)
		} else {
			ipAddress := c.ClientIP()
			dua := &model.Dau{UserId: user.Id, UserIp: ipAddress}
			dua.Add()
			c.HTML(200, "main/index.html", gin.H{"token": myToken})
		}
	} else {
		// 注册后再登陆
		// Set encrypt password
		hashPass, err := encrypt.HashPassword(userRegisterDTO.Password)
		if err != nil {
			global.Logger.Errorf("err: %v", err)
		}
		userRegisterDTO.Password = hashPass

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

		// 返回 Token 进行登陆
		claims := jwt.NewClaims(lastId, username, path)
		myToken, err := jwt.NewJwtToken(claims)
		if err != nil {
			global.Logger.Error("Error creating token: ", err)
			c.JSON(http.StatusUnauthorized, result.InvalidJwtToken)
		} else {
			ipAddress := c.ClientIP()
			dua := &model.Dau{UserId: lastId, UserIp: ipAddress}
			dua.Add()
			c.HTML(200, "main/index.html", gin.H{"token": myToken})
		}
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
