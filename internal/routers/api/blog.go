package api

import (
	"github.com/dsxg666/web-tool/global"
	"github.com/dsxg666/web-tool/internal/model"
	"github.com/dsxg666/web-tool/pkg/result"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

type Blog struct{}

func NewBlog() Blog {
	return Blog{}
}

func (Blog) List(c *gin.Context) {
	var temp Page

	if err := c.ShouldBindJSON(&temp); err != nil {
		global.Logger.Errorf("err: %v", err)
		c.JSON(http.StatusBadRequest, result.InvalidRequestData)
		return
	}

	p := &model.PostsDTO{}
	c.JSON(http.StatusOK, result.SuccessWithData(p.List((temp.CurrentPage-1)*100)))
}

type Page struct {
	CurrentPage int `json:"currentPage"`
}

func (Blog) MyList(c *gin.Context) {
	var temp Page

	if err := c.ShouldBindJSON(&temp); err != nil {
		global.Logger.Errorf("err: %v", err)
		c.JSON(http.StatusBadRequest, result.InvalidRequestData)
		return
	}

	userIdAny, _ := c.Get("userId")
	userId, _ := userIdAny.(string)

	p := &model.PostsDTO{UserId: userId}
	c.JSON(http.StatusOK, result.SuccessWithData(p.MyList((temp.CurrentPage-1)*100)))
}

func (Blog) GetMyListTotalCount(c *gin.Context) {
	userIdAny, _ := c.Get("userId")
	userId, _ := userIdAny.(string)

	p := &model.PostsDTO{UserId: userId}
	c.JSON(http.StatusOK, result.SuccessWithData(p.GetMyListTotalCount()))
}

func (Blog) GetListTotalCount(c *gin.Context) {
	p := &model.PostsDTO{}
	c.JSON(http.StatusOK, result.SuccessWithData(p.GetListTotalCount()))
}

func (Blog) GetFavoritePostList(c *gin.Context) {
	var temp Page

	if err := c.ShouldBindJSON(&temp); err != nil {
		global.Logger.Errorf("err: %v", err)
		c.JSON(http.StatusBadRequest, result.InvalidRequestData)
		return
	}

	userIdAny, _ := c.Get("userId")
	userId, _ := userIdAny.(string)

	f := &model.Favorites{UserId: userId}
	postIdList := f.List((temp.CurrentPage - 1) * 100)

	tempP := &model.PostsDTO{}
	ids := tempP.GetFavoriteAndPublicIds(postIdList)
	if len(postIdList) > 0 {
		c.JSON(http.StatusOK, result.SuccessWithData(tempP.ListByIds(ids)))
	}
}

func (Blog) GetFavoritesTotalCount(c *gin.Context) {
	userIdAny, _ := c.Get("userId")
	userId, _ := userIdAny.(string)

	f := &model.Favorites{UserId: userId}
	c.JSON(200, result.SuccessWithData(f.GetUserCount()))
}

func (Blog) Delete(c *gin.Context) {
	var idDTO IdDTO

	if err := c.ShouldBindJSON(&idDTO); err != nil {
		global.Logger.Errorf("err: %v", err)
		c.JSON(http.StatusBadRequest, result.InvalidRequestData)
		return
	}

	p := &model.PostsDTO{Id: idDTO.Id}
	p.Delete()
	c.JSON(http.StatusOK, result.Success)
}

func (Blog) Edit(c *gin.Context) {
	var p model.PostsDTO

	if err := c.ShouldBindJSON(&p); err != nil {
		global.Logger.Errorf("err: %v", err)
		c.JSON(http.StatusBadRequest, result.InvalidRequestData)
		return
	}

	p.Update()
	c.JSON(http.StatusOK, result.Success)
}

func (Blog) Add(c *gin.Context) {
	var p model.PostsDTO

	if err := c.ShouldBindJSON(&p); err != nil {
		global.Logger.Errorf("err: %v", err)
		c.JSON(http.StatusBadRequest, result.InvalidRequestData)
		return
	}

	p.Add()
	c.JSON(http.StatusOK, result.Success)
}

func (Blog) UploadImg(c *gin.Context) {
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

	filePath := filepath.Join(wd, "storage/blog/"+userId, file.Filename)

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save the file"})
		global.Logger.Errorf("Unable to save the file: %v", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"filename": file.Filename})
}

func (Blog) IsLike(c *gin.Context) {
	var idDTO IdDTO

	if err := c.ShouldBindJSON(&idDTO); err != nil {
		global.Logger.Errorf("err: %v", err)
		c.JSON(http.StatusBadRequest, result.InvalidRequestData)
		return
	}

	userIdAny, _ := c.Get("userId")
	userId, _ := userIdAny.(string)

	l := &model.Likes{PostId: idDTO.Id, UserId: userId}
	if l.IsLike() {
		c.JSON(http.StatusOK, result.SuccessWithData("1"))
	} else {
		c.JSON(http.StatusOK, result.SuccessWithData("0"))
	}
}

func (Blog) IsFavorite(c *gin.Context) {
	var idDTO IdDTO

	if err := c.ShouldBindJSON(&idDTO); err != nil {
		global.Logger.Errorf("err: %v", err)
		c.JSON(http.StatusBadRequest, result.InvalidRequestData)
		return
	}

	userIdAny, _ := c.Get("userId")
	userId, _ := userIdAny.(string)

	f := &model.Favorites{PostId: idDTO.Id, UserId: userId}
	if f.IsFavorite() {
		c.JSON(http.StatusOK, result.SuccessWithData("1"))
	} else {
		c.JSON(http.StatusOK, result.SuccessWithData("0"))
	}
}

func (Blog) Like(c *gin.Context) {
	var idDTO IdDTO

	if err := c.ShouldBindJSON(&idDTO); err != nil {
		global.Logger.Errorf("err: %v", err)
		c.JSON(http.StatusBadRequest, result.InvalidRequestData)
		return
	}

	userIdAny, _ := c.Get("userId")
	userId, _ := userIdAny.(string)

	l := &model.Likes{PostId: idDTO.Id, UserId: userId}
	l.Add()
	c.JSON(200, result.Success)
}

func (Blog) Favorite(c *gin.Context) {
	var idDTO IdDTO

	if err := c.ShouldBindJSON(&idDTO); err != nil {
		global.Logger.Errorf("err: %v", err)
		c.JSON(http.StatusBadRequest, result.InvalidRequestData)
		return
	}

	userIdAny, _ := c.Get("userId")
	userId, _ := userIdAny.(string)

	f := &model.Favorites{PostId: idDTO.Id, UserId: userId}
	f.Add()
	c.JSON(200, result.Success)
}

func (Blog) CancelLike(c *gin.Context) {
	var idDTO IdDTO

	if err := c.ShouldBindJSON(&idDTO); err != nil {
		global.Logger.Errorf("err: %v", err)
		c.JSON(http.StatusBadRequest, result.InvalidRequestData)
		return
	}

	userIdAny, _ := c.Get("userId")
	userId, _ := userIdAny.(string)

	l := &model.Likes{PostId: idDTO.Id, UserId: userId}
	l.Delete()
	c.JSON(200, result.Success)
}

func (Blog) CancelFavorite(c *gin.Context) {
	var idDTO IdDTO

	if err := c.ShouldBindJSON(&idDTO); err != nil {
		global.Logger.Errorf("err: %v", err)
		c.JSON(http.StatusBadRequest, result.InvalidRequestData)
		return
	}

	userIdAny, _ := c.Get("userId")
	userId, _ := userIdAny.(string)

	f := &model.Favorites{PostId: idDTO.Id, UserId: userId}
	f.Delete()
	c.JSON(http.StatusOK, result.Success)
}

func (Blog) GetLikesCount(c *gin.Context) {
	var idDTO IdDTO

	if err := c.ShouldBindJSON(&idDTO); err != nil {
		global.Logger.Errorf("err: %v", err)
		c.JSON(http.StatusBadRequest, result.InvalidRequestData)
		return
	}

	l := &model.Likes{PostId: idDTO.Id}

	c.JSON(http.StatusOK, result.SuccessWithData(strconv.Itoa(l.GetCount())))
}

func (Blog) AddView(c *gin.Context) {
	var idDTO IdDTO

	if err := c.ShouldBindJSON(&idDTO); err != nil {
		global.Logger.Errorf("err: %v", err)
		c.JSON(http.StatusBadRequest, result.InvalidRequestData)
		return
	}

	p := &model.PostsDTO{Id: idDTO.Id}
	p.AddView()
	c.JSON(http.StatusOK, result.Success)
}

func (Blog) GetFavoritesCount(c *gin.Context) {
	var idDTO IdDTO

	if err := c.ShouldBindJSON(&idDTO); err != nil {
		global.Logger.Errorf("err: %v", err)
		c.JSON(http.StatusBadRequest, result.InvalidRequestData)
		return
	}

	f := &model.Favorites{PostId: idDTO.Id}

	c.JSON(http.StatusOK, result.SuccessWithData(strconv.Itoa(f.GetPostCount())))
}

func (Blog) GetById2(c *gin.Context) {
	var idDTO IdDTO

	if err := c.ShouldBindJSON(&idDTO); err != nil {
		global.Logger.Errorf("err: %v", err)
		c.JSON(http.StatusBadRequest, result.InvalidRequestData)
		return
	}

	p := &model.PostsDTO{Id: idDTO.Id}
	if p.IsExistAndValid2() {
		c.JSON(http.StatusOK, result.SuccessWithData(p.GetById()))
	} else {
		c.JSON(200, result.InvalidRequestData)
	}
}

func (Blog) GetById(c *gin.Context) {
	var idDTO IdDTO

	if err := c.ShouldBindJSON(&idDTO); err != nil {
		global.Logger.Errorf("err: %v", err)
		c.JSON(http.StatusBadRequest, result.InvalidRequestData)
		return
	}

	userIdAny, _ := c.Get("userId")
	userId, _ := userIdAny.(string)

	p := &model.PostsDTO{Id: idDTO.Id, UserId: userId}
	if p.IsExistAndValid() {
		c.JSON(http.StatusOK, result.SuccessWithData(p.GetById()))
	} else {
		c.JSON(200, result.InvalidRequestData)
	}
}

func (Blog) ToPublic(c *gin.Context) {
	var idDTO IdDTO

	if err := c.ShouldBindJSON(&idDTO); err != nil {
		global.Logger.Errorf("err: %v", err)
		c.JSON(http.StatusBadRequest, result.InvalidRequestData)
		return
	}

	p := &model.PostsDTO{Id: idDTO.Id}
	p.ToPublic()

	c.JSON(http.StatusOK, result.Success)
}

func (Blog) ToPrivate(c *gin.Context) {
	var idDTO IdDTO

	if err := c.ShouldBindJSON(&idDTO); err != nil {
		global.Logger.Errorf("err: %v", err)
		c.JSON(http.StatusBadRequest, result.InvalidRequestData)
		return
	}

	p := &model.PostsDTO{Id: idDTO.Id}
	p.ToPrivate()

	c.JSON(http.StatusOK, result.Success)
}
