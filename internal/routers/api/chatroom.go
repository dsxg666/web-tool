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

type Chatroom struct{}

func NewChatroom() Chatroom {
	return Chatroom{}
}

func (Chatroom) SendMessage(c *gin.Context) {
	var message model.Messages

	if err := c.ShouldBindJSON(&message); err != nil {
		global.Logger.Errorf("err: %v", err)
		c.JSON(http.StatusBadRequest, result.InvalidRequestData)
		return
	}

	c.JSON(http.StatusOK, result.SuccessWithData(message.Add()))
}

func (Chatroom) SendGroupMessage(c *gin.Context) {
	var message model.GroupMessages

	if err := c.ShouldBindJSON(&message); err != nil {
		global.Logger.Errorf("err: %v", err)
		c.JSON(http.StatusBadRequest, result.InvalidRequestData)
		return
	}

	c.JSON(http.StatusOK, result.SuccessWithData(message.Add()))
}

type SpecialGroupMessage struct {
	GroupId string `json:"groupId"`
	Count   int    `json:"count"`
}

func (Chatroom) GetGroupMessage(c *gin.Context) {
	var temp SpecialGroupMessage

	if err := c.ShouldBindJSON(&temp); err != nil {
		global.Logger.Errorf("err: %v", err)
		c.JSON(http.StatusBadRequest, result.InvalidRequestData)
		return
	}

	groupMessage := &model.GroupMessages{GroupId: temp.GroupId}
	c.JSON(http.StatusOK, result.SuccessWithData(groupMessage.GetMessage(temp.Count)))
}

type SpecialMessage struct {
	SenderId   string `json:"senderId"`
	ReceiverId string `json:"receiverId"`
	Count      int    `json:"count"`
}

func (Chatroom) GetMessage(c *gin.Context) {
	var temp SpecialMessage

	if err := c.ShouldBindJSON(&temp); err != nil {
		global.Logger.Errorf("err: %v", err)
		c.JSON(http.StatusBadRequest, result.InvalidRequestData)
		return
	}

	message := model.Messages{SenderId: temp.SenderId, ReceiverId: temp.ReceiverId}

	c.JSON(http.StatusOK, result.SuccessWithData(message.GetMessage(temp.Count)))
}

func (Chatroom) BelongGroupUserList(c *gin.Context) {
	var idDTO IdDTO

	if err := c.ShouldBindJSON(&idDTO); err != nil {
		global.Logger.Errorf("err: %v", err)
		c.JSON(http.StatusBadRequest, result.InvalidRequestData)
		return
	}

	groupMembers := &model.GroupMembers{GroupId: idDTO.Id}
	c.JSON(http.StatusOK, result.SuccessWithData(groupMembers.BelongGroupList()))
}

func (Chatroom) BelongUserGroupList(c *gin.Context) {
	userIdAny, _ := c.Get("userId")
	userId, _ := userIdAny.(string)
	groupMembers := &model.GroupMembers{UserId: userId}
	c.JSON(http.StatusOK, result.SuccessWithData(groupMembers.BelongUserList()))
}

func (Chatroom) GetFriendList(c *gin.Context) {
	userIdAny, _ := c.Get("userId")
	userId, _ := userIdAny.(string)
	f := &model.Friends{SelfId: userId}
	c.JSON(http.StatusOK, result.SuccessWithData(f.FriendList()))
}

func (Chatroom) GetRequestList(c *gin.Context) {
	userIdAny, _ := c.Get("userId")
	userId, _ := userIdAny.(string)
	rs := &model.Requests{ReceiverId: userId}
	c.JSON(http.StatusOK, result.SuccessWithData(rs.List()))
}

func (Chatroom) HandleRequest(c *gin.Context) {
	var rs model.Requests

	if err := c.ShouldBindJSON(&rs); err != nil {
		global.Logger.Errorf("err: %v", err)
		c.JSON(http.StatusBadRequest, result.InvalidRequestData)
		return
	}

	rs.HandleRequest()
	c.JSON(http.StatusOK, result.Success)
}

func (Chatroom) IsGroupExist(c *gin.Context) {
	var idDTO IdDTO

	if err := c.ShouldBindJSON(&idDTO); err != nil {
		global.Logger.Errorf("err: %v", err)
		c.JSON(http.StatusBadRequest, result.InvalidRequestData)
		return
	}

	g := &model.Groups{Id: idDTO.Id}
	if g.IsExist() {
		c.JSON(http.StatusOK, result.SuccessWithData("1"))
	} else {
		c.JSON(http.StatusOK, result.SuccessWithData("0"))
	}
}

func (Chatroom) IsInGroup(c *gin.Context) {
	var gm model.GroupMembers

	if err := c.ShouldBindJSON(&gm); err != nil {
		global.Logger.Errorf("err: %v", err)
		c.JSON(http.StatusBadRequest, result.InvalidRequestData)
		return
	}

	if gm.IsInGroup() {
		c.JSON(http.StatusOK, result.SuccessWithData("1"))
	} else {
		c.JSON(http.StatusOK, result.SuccessWithData("0"))
	}
}

func (Chatroom) IsOnline(c *gin.Context) {
	var idDTO IdDTO

	if err := c.ShouldBindJSON(&idDTO); err != nil {
		global.Logger.Errorf("err: %v", err)
		c.JSON(http.StatusBadRequest, result.InvalidRequestData)
		return
	}

	if contains(global.OnlineUser, idDTO.Id) {
		c.JSON(http.StatusOK, result.SuccessWithData("1"))
	} else {
		c.JSON(http.StatusOK, result.SuccessWithData("0"))
	}
}

func contains(slice []string, element string) bool {
	for _, v := range slice {
		if v == element {
			return true
		}
	}
	return false
}

func (Chatroom) IsFriend(c *gin.Context) {
	var f model.Friends

	if err := c.ShouldBindJSON(&f); err != nil {
		global.Logger.Errorf("err: %v", err)
		c.JSON(http.StatusBadRequest, result.InvalidRequestData)
		return
	}

	if f.IsFriend() {
		c.JSON(http.StatusOK, result.SuccessWithData("1"))
	} else {
		c.JSON(http.StatusOK, result.SuccessWithData("0"))
	}
}

func (Chatroom) UploadAvatar(c *gin.Context) {
	file, err := c.FormFile("file")

	groupId := c.Query("id")

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

	g := &model.Groups{Id: groupId, Avatar: groupId + ".png"}
	g.ModifyAvatar()

	filePath := filepath.Join(wd, "storage/chatroom/avatars", groupId+".png")

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save the file"})
		return
	}

	c.JSON(http.StatusOK, result.Success)
}

func (Chatroom) ModifyGroupName(c *gin.Context) {
	var g model.Groups

	if err := c.ShouldBindJSON(&g); err != nil {
		global.Logger.Errorf("err: %v", err)
		c.JSON(http.StatusBadRequest, result.InvalidRequestData)
		return
	}

	g.ModifyGroupName()
	c.JSON(http.StatusOK, result.Success)
}

func (Chatroom) DeleteMember(c *gin.Context) {
	var gm model.GroupMembers

	if err := c.ShouldBindJSON(&gm); err != nil {
		global.Logger.Errorf("err: %v", err)
		c.JSON(http.StatusBadRequest, result.InvalidRequestData)
		return
	}

	if gm.IsInGroup() {
		gm.DeleteGroupMember()
	}
	c.JSON(http.StatusOK, result.Success)
}

func (Chatroom) DeleteFriend(c *gin.Context) {
	var idDTO IdDTO

	if err := c.ShouldBindJSON(&idDTO); err != nil {
		global.Logger.Errorf("err: %v", err)
		c.JSON(http.StatusBadRequest, result.InvalidRequestData)
		return
	}

	userIdAny, _ := c.Get("userId")
	userId, _ := userIdAny.(string)

	fs := &model.Friends{SelfId: userId, FriendId: idDTO.Id}
	if fs.IsFriend() {
		fs.DeleteFriend()
		ms := &model.Messages{SenderId: userId, ReceiverId: idDTO.Id}
		ms.DeleteFriendMessage()
	}
	c.JSON(http.StatusOK, result.Success)
}

func (Chatroom) DeleteGroup(c *gin.Context) {
	var idDTO IdDTO

	if err := c.ShouldBindJSON(&idDTO); err != nil {
		global.Logger.Errorf("err: %v", err)
		c.JSON(http.StatusBadRequest, result.InvalidRequestData)
		return
	}

	g := &model.Groups{Id: idDTO.Id}
	g.Delete()
	gm := &model.GroupMembers{GroupId: idDTO.Id}
	gm.DeleteGroupMembers()
	gm2 := &model.GroupMessages{GroupId: idDTO.Id}
	gm2.DeleteGroupMessages()
	c.JSON(http.StatusOK, result.Success)
}

func (Chatroom) AddGroup(c *gin.Context) {
	var gs model.Groups

	if err := c.ShouldBindJSON(&gs); err != nil {
		global.Logger.Errorf("err: %v", err)
		c.JSON(http.StatusBadRequest, result.InvalidRequestData)
		return
	}

	userIdAny, _ := c.Get("userId")
	userId, _ := userIdAny.(string)

	lastId := gs.Add()
	ms := &model.GroupMembers{GroupId: lastId, UserId: userId, Status: "1"}
	ms.Add()
	c.JSON(http.StatusOK, result.Success)
}

func (Chatroom) EnterGroup(c *gin.Context) {
	var idDTO IdDTO

	if err := c.ShouldBindJSON(&idDTO); err != nil {
		global.Logger.Errorf("err: %v", err)
		c.JSON(http.StatusBadRequest, result.InvalidRequestData)
		return
	}

	userIdAny, _ := c.Get("userId")
	userId, _ := userIdAny.(string)

	ms := &model.GroupMembers{GroupId: idDTO.Id, UserId: userId, Status: "0"}
	ms.Add()
	c.JSON(http.StatusOK, result.Success)
}

func (Chatroom) AddFriend(c *gin.Context) {
	var idDTO IdDTO

	if err := c.ShouldBindJSON(&idDTO); err != nil {
		global.Logger.Errorf("err: %v", err)
		c.JSON(http.StatusBadRequest, result.InvalidRequestData)
		return
	}

	userIdAny, _ := c.Get("userId")
	userId, _ := userIdAny.(string)
	f := &model.Friends{SelfId: userId, FriendId: idDTO.Id}
	if f.IsFriend() {
		c.JSON(http.StatusOK, result.OperateError("You're already friends!", ""))
	} else {
		f.Add()
		c.JSON(http.StatusOK, result.Success)
	}
}

func (Chatroom) EnterGroupRequest(c *gin.Context) {
	var rs model.Requests

	if err := c.ShouldBindJSON(&rs); err != nil {
		global.Logger.Errorf("err: %v", err)
		c.JSON(http.StatusBadRequest, result.InvalidRequestData)
		return
	}

	userIdDTO := &model.UserIdDTO{Id: rs.ReceiverId}
	if userIdDTO.IsIdExist() {
		gs := &model.GroupMembers{GroupId: rs.GroupId, UserId: rs.ReceiverId}
		if gs.IsInGroup() {
			c.JSON(http.StatusOK, result.OperateError("The user already in the group chat!", ""))
		} else {
			if rs.IsGroupRequestExist() {
				c.JSON(http.StatusOK, result.OperateError("The request has been sent!", ""))
			} else {
				rs.AddGroupRequest()
				c.JSON(http.StatusOK, result.Success)
			}
		}
	} else {
		c.JSON(http.StatusOK, result.OperateError("The UserId does not exist!", ""))
	}
}

type RequestTemp struct {
	Id     string `json:"id"`
	Remark string `json:"remark"`
}

func (Chatroom) AddFriendRequest(c *gin.Context) {
	var requestTemp RequestTemp

	if err := c.ShouldBindJSON(&requestTemp); err != nil {
		global.Logger.Errorf("err: %v", err)
		c.JSON(http.StatusBadRequest, result.InvalidRequestData)
		return
	}

	userIdAny, _ := c.Get("userId")
	userId, _ := userIdAny.(string)
	userIdDTO := &model.UserIdDTO{Id: requestTemp.Id}
	if userIdDTO.IsIdExist() {
		if userId == requestTemp.Id {
			c.JSON(http.StatusOK, result.OperateError("Can't add yourself!", ""))
		} else {
			f := &model.Friends{SelfId: userId, FriendId: requestTemp.Id}
			if f.IsFriend() {
				c.JSON(http.StatusOK, result.OperateError("You're already friends!", ""))
			} else {
				rs := &model.Requests{UserId: userId, ReceiverId: requestTemp.Id, Remark: requestTemp.Remark}
				if rs.IsUserRequestExist() {
					c.JSON(http.StatusOK, result.OperateError("The request has been sent!", ""))
				} else {
					rs.AddUserRequest()
					c.JSON(http.StatusOK, result.Success)
				}
			}
		}
	} else {
		c.JSON(http.StatusOK, result.OperateError("The UserId does not exist!", ""))
	}
}

func (Chatroom) GetGroupById(c *gin.Context) {
	var idDTO IdDTO

	if err := c.ShouldBindJSON(&idDTO); err != nil {
		global.Logger.Errorf("err: %v", err)
		c.JSON(http.StatusBadRequest, result.InvalidRequestData)
		return
	}

	groups := &model.Groups{Id: idDTO.Id}
	c.JSON(http.StatusOK, result.SuccessWithData(groups.GetById()))
}
