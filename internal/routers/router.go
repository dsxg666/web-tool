package routers

import (
	_ "github.com/dsxg666/web-tool/docs"
	"github.com/dsxg666/web-tool/internal/middleware"
	"github.com/dsxg666/web-tool/internal/routers/api"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	r.Use(middleware.CorsMiddleware())

	authorizedGroup := r.Group("/api/auth")
	authorizedGroup.Use(middleware.AuthMiddleware())
	{
		blog := api.NewBlog()
		blogGroup := authorizedGroup.Group("blog")
		{
			blogGroup.POST("/list", blog.List)
			blogGroup.POST("/myList", blog.MyList)
			blogGroup.POST("/add", blog.Add)
			blogGroup.POST("/edit", blog.Edit)
			blogGroup.POST("/delete", blog.Delete)
			blogGroup.POST("/toPublic", blog.ToPublic)
			blogGroup.POST("/toPrivate", blog.ToPrivate)
			blogGroup.POST("/uploadImg", blog.UploadImg)
			blogGroup.POST("/getById", blog.GetById)
			blogGroup.POST("/getById2", blog.GetById2)
			blogGroup.POST("/getLikesCount", blog.GetLikesCount)
			blogGroup.POST("/getFavoritesCount", blog.GetFavoritesCount)
			blogGroup.POST("/isLike", blog.IsLike)
			blogGroup.POST("/isFavorite", blog.IsFavorite)
			blogGroup.POST("/like", blog.Like)
			blogGroup.POST("/favorite", blog.Favorite)
			blogGroup.POST("/cancelLike", blog.CancelLike)
			blogGroup.POST("/cancelFavorite", blog.CancelFavorite)
			blogGroup.POST("/addView", blog.AddView)
			blogGroup.POST("/getMyListTotalCount", blog.GetMyListTotalCount)
			blogGroup.POST("/getListTotalCount", blog.GetListTotalCount)
			blogGroup.POST("/getFavoritesTotalCount", blog.GetFavoritesTotalCount)
			blogGroup.POST("/getFavoritePostList", blog.GetFavoritePostList)
		}

		statistic := api.NewStatistic()
		statisticGroup := authorizedGroup.Group("/statistic")
		{
			statisticGroup.POST("/getDailyDauData", statistic.GetDailyDauData)
			statisticGroup.POST("/getDailyMessageData", statistic.GetDailyMessageData)
			statisticGroup.POST("/getDailyTodolistData", statistic.GetDailyTodolistData)
			statisticGroup.POST("/getChartDauData", statistic.GetChartDauData)
			statisticGroup.POST("/getChartMessageData", statistic.GetChartMessageData)
			statisticGroup.POST("/getChartTodolistData", statistic.GetChartTodolistData)
		}

		user := api.NewUser()
		userGroup := authorizedGroup.Group("/user")
		{
			userGroup.POST("/getUserById", user.GetById)
			userGroup.POST("/modifyEmail", user.ModifyEmail)
			userGroup.POST("/modifyPassword", user.ModifyPassword)
			userGroup.POST("/modifyUsername", user.ModifyUsername)
			userGroup.POST("/uploadAvatar", user.UploadAvatar)
		}

		todolist := api.NewTodoList()
		todolistGroup := authorizedGroup.Group("/todolist")
		{
			todolistGroup.POST("/list", todolist.List)
			todolistGroup.POST("/add", todolist.Add)
			todolistGroup.POST("/delete", todolist.Delete)
			todolistGroup.GET("/getById", todolist.GetById)
			todolistGroup.POST("/edit", todolist.Edit)
			todolistGroup.POST("/uploadImg", todolist.UploadImg)
		}

		chatroom := api.NewChatroom()
		chatroomGroup := authorizedGroup.Group("/chatroom")
		{
			chatroomGroup.POST("/belongUserGroupList", chatroom.BelongUserGroupList)
			chatroomGroup.POST("/belongGroupUserList", chatroom.BelongGroupUserList)
			chatroomGroup.POST("/getMessage", chatroom.GetMessage)
			chatroomGroup.POST("/getGroupMessage", chatroom.GetGroupMessage)
			chatroomGroup.POST("/getGroupById", chatroom.GetGroupById)
			chatroomGroup.POST("/sendMessage", chatroom.SendMessage)
			chatroomGroup.POST("/sendGroupMessage", chatroom.SendGroupMessage)
			chatroomGroup.POST("/getFriendList", chatroom.GetFriendList)
			chatroomGroup.POST("/addFriendRequest", chatroom.AddFriendRequest)
			chatroomGroup.POST("/enterGroupRequest", chatroom.EnterGroupRequest)
			chatroomGroup.POST("/getRequestList", chatroom.GetRequestList)
			chatroomGroup.POST("/handleRequest", chatroom.HandleRequest)
			chatroomGroup.POST("/addFriend", chatroom.AddFriend)
			chatroomGroup.POST("/enterGroup", chatroom.EnterGroup)
			chatroomGroup.POST("/addGroup", chatroom.AddGroup)
			chatroomGroup.POST("/deleteGroup", chatroom.DeleteGroup)
			chatroomGroup.POST("/deleteFriend", chatroom.DeleteFriend)
			chatroomGroup.POST("/deleteMember", chatroom.DeleteMember)
			chatroomGroup.POST("/modifyGroupName", chatroom.ModifyGroupName)
			chatroomGroup.POST("/uploadAvatar", chatroom.UploadAvatar)
			chatroomGroup.POST("/isGroupExist", chatroom.IsGroupExist)
			chatroomGroup.POST("/isInGroup", chatroom.IsInGroup)
			chatroomGroup.POST("/isFriend", chatroom.IsFriend)
			chatroomGroup.POST("/isOnline", chatroom.IsOnline)
		}
	}

	base := api.NewBase()
	baseGroup := r.Group("/api/base")
	{
		baseGroup.POST("/loginByPassword", base.LoginByPassword)
		baseGroup.POST("/loginByCode", base.LoginByCode)
		baseGroup.POST("/register", base.Register)
		baseGroup.POST("/verifyToken", base.VerifyToken)
		baseGroup.GET("/getRegisterCode", base.GetRegisterCode)
		baseGroup.GET("/getVerificationCode", base.GetVerificationCode)
		baseGroup.GET("/getAvatar/:userId", base.GetAvatar)
		baseGroup.GET("/getGroupAvatar/:filename", base.GetGroupAvatar)
		baseGroup.GET("/getTodolistImg/:filename", base.GetTodolistImg)
		baseGroup.GET("/getBlogImg/:filename", base.GetBlogImg)
	}

	return r
}
