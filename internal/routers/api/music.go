package api

import (
	"github.com/dsxg666/web-tool/global"
	"github.com/dsxg666/web-tool/internal/model"
	"github.com/dsxg666/web-tool/pkg/result"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Music struct{}

func NewMusic() Music {
	return Music{}
}

func (Music) List(c *gin.Context) {
	s := &model.Songs{}
	c.JSON(http.StatusOK, result.SuccessWithData(s.List()))
}

func (Music) ListByFavorite(c *gin.Context) {
	userIdAny, _ := c.Get("userId")
	userId, _ := userIdAny.(string)
	sf := &model.SongFavorites{UserId: userId}
	ids := sf.List()
	s := &model.Songs{}
	c.JSON(http.StatusOK, result.SuccessWithData(s.ListByFavorite(ids)))
}

func (Music) FavoriteList(c *gin.Context) {
	userIdAny, _ := c.Get("userId")
	userId, _ := userIdAny.(string)
	sf := &model.SongFavorites{UserId: userId}
	c.JSON(http.StatusOK, result.SuccessWithData(sf.List()))
}

func (Music) Favorite(c *gin.Context) {
	var idDTO IdDTO

	if err := c.ShouldBindJSON(&idDTO); err != nil {
		global.Logger.Errorf("err: %v", err)
		c.JSON(http.StatusBadRequest, result.InvalidRequestData)
		return
	}

	userIdAny, _ := c.Get("userId")
	userId, _ := userIdAny.(string)

	sf := &model.SongFavorites{UserId: userId, SongId: idDTO.Id}
	sf.Add()
	c.JSON(http.StatusOK, result.Success)
}

func (Music) CancelFavorite(c *gin.Context) {
	var idDTO IdDTO

	if err := c.ShouldBindJSON(&idDTO); err != nil {
		global.Logger.Errorf("err: %v", err)
		c.JSON(http.StatusBadRequest, result.InvalidRequestData)
		return
	}

	userIdAny, _ := c.Get("userId")
	userId, _ := userIdAny.(string)

	sf := &model.SongFavorites{UserId: userId, SongId: idDTO.Id}
	sf.Delete()
	c.JSON(http.StatusOK, result.Success)
}
