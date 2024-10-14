package api

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/dsxg666/web-tool/global"
	"github.com/dsxg666/web-tool/internal/model"
	"github.com/dsxg666/web-tool/pkg/result"
	"github.com/gin-gonic/gin"
)

type ListDTO struct {
	Date string `json:"date"`
}

type DeleteDTO struct {
	Id string `json:"id"`
}

type TodoList struct{}

func NewTodoList() TodoList {
	return TodoList{}
}

func (TodoList) UploadImg(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file is received"})
		return
	}

	wd, err := os.Getwd()
	if err != nil {
		global.Logger.Errorf("Failed to get working directory err: %v", err)
	}

	userIdAny, _ := c.Get("userId")
	userId, _ := userIdAny.(string)

	filePath := filepath.Join(wd, "storage/todolist/"+userId, file.Filename)

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save the file"})
		global.Logger.Errorf("Unable to save the file: %v", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"filename": file.Filename})
}

func (TodoList) Edit(c *gin.Context) {
	var todoEditDTO model.TodoEditDTO

	if err := c.ShouldBindJSON(&todoEditDTO); err != nil {
		global.Logger.Errorf("err: %v", err)
		c.JSON(http.StatusBadRequest, result.InvalidRequestData)
		return
	}

	todoEditDTO.Edit()
	c.JSON(http.StatusOK, result.Success)
}

func (TodoList) GetById(c *gin.Context) {
	userIdAny, _ := c.Get("userId")
	userId, _ := userIdAny.(string)

	id := c.Query("id")
	t := &model.TodoEditDTO{Id: id, UserID: userId}
	if t.IsExistAndValid() {
		todoIdDTO := &model.TodoIdDTO{Id: id}
		c.JSON(http.StatusOK, result.SuccessWithData(todoIdDTO.GetById()))
	} else {
		c.JSON(200, result.InvalidRequestData)
	}
}

func (TodoList) Delete(c *gin.Context) {
	var deleteDTO DeleteDTO

	if err := c.ShouldBindJSON(&deleteDTO); err != nil {
		global.Logger.Errorf("err: %v", err)
		c.JSON(http.StatusBadRequest, result.InvalidRequestData)
		return
	}

	todoIdDTO := &model.TodoIdDTO{Id: deleteDTO.Id}
	todoIdDTO.Delete()
	c.JSON(http.StatusOK, result.Success)
}

func (TodoList) Add(c *gin.Context) {
	var todoAddDTO model.TodoAddDTO

	if err := c.ShouldBindJSON(&todoAddDTO); err != nil {
		global.Logger.Errorf("err: %v", err)
		c.JSON(http.StatusBadRequest, result.InvalidRequestData)
		return
	}

	todoAddDTO.Add()
	c.JSON(http.StatusOK, result.Success)
}

func (TodoList) List(c *gin.Context) {
	var listDTO ListDTO

	if err := c.ShouldBindJSON(&listDTO); err != nil {
		global.Logger.Errorf("err: %v", err)
		c.JSON(http.StatusBadRequest, result.InvalidRequestData)
		return
	}

	userIdAny, _ := c.Get("userId")
	userId, _ := userIdAny.(string)
	todoListDTO := &model.TodoListDTO{UserID: userId, CreatedAt: listDTO.Date}
	list := todoListDTO.List()
	c.JSON(http.StatusOK, result.SuccessWithData(list))
}
