package model

import (
	"github.com/dsxg666/web-tool/global"
	"github.com/dsxg666/web-tool/pkg/util"
	"strconv"
)

type Requests struct {
	Id         string `json:"id"`
	UserId     string `json:"userId"`
	GroupId    string `json:"groupId"`
	ReceiverId string `json:"receiverId"`
	Remark     string `json:"remark"`
	Type       string `json:"type"`
	Finish     string `json:"finish"`
	Result     string `json:"result"`
	CreatedAt  string `json:"createdAt"`
}

func (r *Requests) HandleRequest() {
	sql := "UPDATE `requests` SET `finish` = '1', `result` = ? WHERE `id` = ?;"
	_, err := global.Database.DbHandle.Exec(sql, r.Result, r.Id)
	if err != nil {
		global.Logger.Errorf("update error: %v", err)
	}
}

func (r *Requests) IsUserRequestExist() bool {
	sql := "SELECT COUNT(*) AS count FROM `requests` WHERE `user_id` = ? AND `receiver_id` = ? AND `finish` = '0';"
	var count int
	err := global.Database.DbHandle.QueryRow(sql, r.UserId, r.ReceiverId).Scan(&count)
	if err != nil {
		global.Logger.Errorf("err: %v", err)
	}
	return count > 0
}

func (r *Requests) IsGroupRequestExist() bool {
	sql := "SELECT COUNT(*) AS count FROM `requests` WHERE `group_id` = ? AND `receiver_id` = ? AND `finish` = '0';"
	var count int
	err := global.Database.DbHandle.QueryRow(sql, r.GroupId, r.ReceiverId).Scan(&count)
	if err != nil {
		global.Logger.Errorf("err: %v", err)
	}
	return count > 0
}

func (r *Requests) AddGroupRequest() {
	sql := "INSERT INTO `requests` (`group_id`, `receiver_id`, `type`) VALUES (?, ?, '1')"
	_, err := global.Database.DbHandle.Exec(sql, r.GroupId, r.ReceiverId)
	if err != nil {
		global.Logger.Errorf("add error: %v", err)
	}
}

func (r *Requests) AddUserRequest() {
	sql := "INSERT INTO `requests` (`user_id`, `receiver_id`, `remark`, `type`) VALUES (?, ?, ?, '0')"
	_, err := global.Database.DbHandle.Exec(sql, r.UserId, r.ReceiverId, r.Remark)
	if err != nil {
		global.Logger.Errorf("add error: %v", err)
	}
}

func (r *Requests) List() []*Requests {
	sql := "SELECT * FROM `requests` WHERE `receiver_id` = ? ORDER BY `created_at` DESC;"
	var rs []*Requests
	rows, err := global.Database.DbHandle.Query(sql, r.ReceiverId)
	if err != nil {
		global.Logger.Errorf("requests.List err: %v", err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var temp Requests
		err = rows.Scan(&temp.Id, &temp.UserId, &temp.GroupId, &temp.ReceiverId, &temp.Remark, &temp.Type, &temp.Finish, &temp.Result, &temp.CreatedAt)
		if err != nil {
			global.Logger.Errorf("err: %v", err)
		}
		rs = append(rs, &temp)
	}

	return rs
}

type Friends struct {
	Id        string `json:"id"`
	SelfId    string `json:"selfId"`
	FriendId  string `json:"friendId"`
	CreatedAt string `json:"createdAt"`
}

func (f *Friends) IsFriend() bool {
	sql := "SELECT COUNT(*) AS count FROM `friends` WHERE `self_id` = ? AND `friend_id` = ?;"
	var count int
	err := global.Database.DbHandle.QueryRow(sql, f.SelfId, f.FriendId).Scan(&count)
	if err != nil {
		global.Logger.Errorf("err: %v", err)
	}
	return count > 0
}

func (f *Friends) DeleteFriend() {
	sql := "DELETE FROM `friends` WHERE `self_id` = ? And `friend_id` = ?;"
	_, err := global.Database.DbHandle.Exec(sql, f.SelfId, f.FriendId)
	_, err2 := global.Database.DbHandle.Exec(sql, f.FriendId, f.SelfId)
	if err != nil || err2 != nil {
		global.Logger.Errorf("add error: %v and %v", err, err2)
	}
}

func (f *Friends) Add() {
	sql := "INSERT INTO `friends` (`self_id`, `friend_id`) VALUES (?, ?)"
	_, err := global.Database.DbHandle.Exec(sql, f.SelfId, f.FriendId)
	_, err2 := global.Database.DbHandle.Exec(sql, f.FriendId, f.SelfId)
	if err != nil || err2 != nil {
		global.Logger.Errorf("add error: %v and %v", err, err2)
	}
}

func (f *Friends) FriendList() []*Friends {
	sql := "SELECT * FROM `friends` WHERE `self_id` = ?;"
	var fs []*Friends
	rows, err := global.Database.DbHandle.Query(sql, f.SelfId)
	if err != nil {
		global.Logger.Errorf("err: %v", err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var temp Friends
		err = rows.Scan(&temp.Id, &temp.SelfId, &temp.FriendId, &temp.CreatedAt)
		if err != nil {
			global.Logger.Errorf("err: %v", err)
		}
		fs = append(fs, &temp)
	}

	return fs
}

type Messages struct {
	Id         string `json:"id"`
	SenderId   string `json:"senderId"`
	ReceiverId string `json:"receiverId"`
	Message    string `json:"message"`
	CreatedAt  string `json:"createdAt"`
}

func (m *Messages) GetPastWeekData(arr []string) []int {
	datas := make([]int, 7)
	sql := "SELECT DATE(created_at) AS date, COUNT(`id`) AS message_num FROM `messages` WHERE `created_at` >= NOW() - INTERVAL 7 DAY GROUP BY DATE(`created_at`) ORDER BY DATE(`created_at`) DESC;"
	rows, err := global.Database.DbHandle.Query(sql)
	if err != nil {
		global.Logger.Errorf("err: %v", err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var temp int
		var dateTemp string
		err = rows.Scan(&dateTemp, &temp)
		if err != nil {
			global.Logger.Errorf("err: %v", err)
		}

		for i, v := range arr {
			if v == util.StrToFormatDate(dateTemp) {
				datas[i] = temp
			}
		}
	}

	return datas
}

func (m *Messages) GetYesterdayMessageNum() int {
	sql := "SELECT COUNT(`id`) AS message_num FROM `messages` WHERE DATE(`created_at`) = DATE_SUB(CURDATE(), INTERVAL 1 DAY);"
	var count int
	err := global.Database.DbHandle.QueryRow(sql).Scan(&count)
	if err != nil {
		global.Logger.Errorf(" GetTodayMessageNum error: %v", err)
	}
	return count
}

func (m *Messages) GetTodayMessageNum() int {
	sql := "SELECT COUNT(`id`) AS message_num FROM `messages` WHERE DATE(`created_at`) = CURDATE();"
	var count int
	err := global.Database.DbHandle.QueryRow(sql).Scan(&count)
	if err != nil {
		global.Logger.Errorf(" GetTodayMessageNum error: %v", err)
	}
	return count
}

func (m *Messages) DeleteFriendMessage() {
	sql := "DELETE FROM `messages` WHERE `sender_id` = ? And `receiver_id` = ?"
	_, err := global.Database.DbHandle.Exec(sql, m.SenderId, m.ReceiverId)
	_, err2 := global.Database.DbHandle.Exec(sql, m.ReceiverId, m.SenderId)
	if err != nil || err2 != nil {
		global.Logger.Errorf("add error: %v and %v", err, err2)
	}
}

func (m *Messages) Add() string {
	sqlStat := "INSERT INTO `messages` (`sender_id`, `receiver_id`, `message`) VALUES (?, ?, ?)"
	res, err := global.Database.DbHandle.Exec(sqlStat, m.SenderId, m.ReceiverId, m.Message)
	if err != nil {
		global.Logger.Errorf("add error: %v", err)
	}
	id, _ := res.LastInsertId()
	return strconv.Itoa(int(id))
}

func (m *Messages) GetMessage(time int) []*Messages {
	sql := "SELECT * FROM (SELECT * FROM `messages` WHERE `sender_id` = ? and `receiver_id` = ? or `sender_id` = ? and `receiver_id` = ? ORDER BY `created_at` DESC LIMIT ?) AS recent_messages ORDER BY `created_at`;"
	var gms []*Messages
	rows, err := global.Database.DbHandle.Query(sql, m.SenderId, m.ReceiverId, m.ReceiverId, m.SenderId, 20*time)
	if err != nil {
		global.Logger.Errorf("GetMessage err: %v", err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var temp Messages
		err = rows.Scan(&temp.Id, &temp.SenderId, &temp.ReceiverId, &temp.Message, &temp.CreatedAt)
		if err != nil {
			global.Logger.Errorf("GetMessage err: %v", err)
		}
		gms = append(gms, &temp)
	}

	return gms
}

type GroupMessages struct {
	Id        string `json:"id"`
	SenderId  string `json:"senderId"`
	GroupId   string `json:"groupId"`
	Message   string `json:"message"`
	CreatedAt string `json:"createdAt"`
}

func (g *GroupMessages) GetPastWeekData(arr []string) []int {
	datas := make([]int, 7)
	sql := "SELECT DATE(created_at) AS date, COUNT(`id`) AS message_num FROM `group_messages` WHERE `created_at` >= NOW() - INTERVAL 7 DAY GROUP BY DATE(`created_at`) ORDER BY DATE(`created_at`) DESC;"
	rows, err := global.Database.DbHandle.Query(sql)
	if err != nil {
		global.Logger.Errorf("err: %v", err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var temp int
		var dateTemp string
		err = rows.Scan(&dateTemp, &temp)
		if err != nil {
			global.Logger.Errorf("err: %v", err)
		}

		for i, v := range arr {
			if v == util.StrToFormatDate(dateTemp) {
				datas[i] = temp
			}
		}
	}

	return datas
}

func (g *GroupMessages) GetYesterdayGroupMessageNum() int {
	sql := "SELECT COUNT(`id`) AS group_message_num FROM `group_messages` WHERE DATE(`created_at`) = DATE_SUB(CURDATE(), INTERVAL 1 DAY);"
	var count int
	err := global.Database.DbHandle.QueryRow(sql).Scan(&count)
	if err != nil {
		global.Logger.Errorf(" GetTodayGroupMessageNum error: %v", err)
	}
	return count
}

func (g *GroupMessages) GetTodayGroupMessageNum() int {
	sql := "SELECT COUNT(`id`) AS group_message_num FROM `group_messages` WHERE DATE(`created_at`) = CURDATE();"
	var count int
	err := global.Database.DbHandle.QueryRow(sql).Scan(&count)
	if err != nil {
		global.Logger.Errorf(" GetTodayGroupMessageNum error: %v", err)
	}
	return count
}

func (g *GroupMessages) DeleteGroupMessages() {
	sql := "DELETE FROM `group_messages` WHERE `group_id` = ?;"
	_, err := global.Database.DbHandle.Exec(sql, g.GroupId)
	if err != nil {
		global.Logger.Errorf("delete error: %v", err)
	}
}

func (g *GroupMessages) Add() string {
	sqlStat := "INSERT INTO group_messages (`sender_id`, `group_id`, `message`) VALUES (?, ?, ?);"
	res, err := global.Database.DbHandle.Exec(sqlStat, g.SenderId, g.GroupId, g.Message)
	if err != nil {
		global.Logger.Errorf("add error: %v", err)
	}
	id, _ := res.LastInsertId()
	return strconv.Itoa(int(id))
}

func (g *GroupMessages) GetMessage(time int) []*GroupMessages {
	sql := "SELECT * FROM (SELECT * FROM `group_messages` WHERE `group_id` = ? ORDER BY `created_at` DESC LIMIT ?) AS recent_messages ORDER BY `created_at`;"
	var gms []*GroupMessages
	rows, err := global.Database.DbHandle.Query(sql, g.GroupId, 20*time)
	if err != nil {
		global.Logger.Errorf("[GetMessage] %v", err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var temp GroupMessages
		err = rows.Scan(&temp.Id, &temp.SenderId, &temp.GroupId, &temp.Message, &temp.CreatedAt)
		if err != nil {
			global.Logger.Errorf("[GetMessage] %v", err)
		}
		gms = append(gms, &temp)
	}

	return gms
}

type Groups struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Avatar    string `json:"avatar"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

func (g *Groups) IsExist() bool {
	sql := "SELECT COUNT(*) AS count FROM `groups` WHERE `id` = ?;"
	var count int
	err := global.Database.DbHandle.QueryRow(sql, g.Id).Scan(&count)
	if err != nil {
		global.Logger.Errorf("err: %v", err)
	}
	return count > 0
}

func (g *Groups) ModifyAvatar() {
	sql := "UPDATE `groups` SET `avatar` = ? WHERE `id` = ?;"
	_, err := global.Database.DbHandle.Exec(sql, g.Avatar, g.Id)
	if err != nil {
		global.Logger.Errorf("modify avatar error: %v", err)
	}
}

func (g *Groups) ModifyGroupName() {
	sql := "UPDATE `groups` SET `name` = ? WHERE `id` = ?;"
	_, err := global.Database.DbHandle.Exec(sql, g.Name, g.Id)
	if err != nil {
		global.Logger.Errorf("modify name error: %v", err)
	}
}

func (g *Groups) Delete() {
	sql := "DELETE FROM `groups` WHERE `id`=?"
	_, err := global.Database.DbHandle.Exec(sql, g.Id)
	if err != nil {
		global.Logger.Errorf("delete error: %v", err)
	}
}

func (g *Groups) GetById() *Groups {
	sql := "SELECT * FROM `groups` WHERE `id`=?"
	var temp Groups
	err := global.Database.DbHandle.QueryRow(sql, g.Id).Scan(&temp.Id, &temp.Name, &temp.Avatar, &temp.CreatedAt, &temp.UpdatedAt)
	if err != nil {
		global.Logger.Errorf("err: %v", err)
	}
	return &temp
}

func (g *Groups) Add() string {
	sql := "INSERT INTO `groups` (`name`, `avatar`) VALUES (?, ?);"
	row, err := global.Database.DbHandle.Exec(sql, g.Name, g.Avatar)
	if err != nil {
		global.Logger.Errorf("err: %v", err)
	}
	id, _ := row.LastInsertId()
	return strconv.Itoa(int(id))
}

type GroupMembers struct {
	Id      string `json:"id"`
	GroupId string `json:"groupId"`
	UserId  string `json:"userId"`
	Status  string `json:"status"`
}

func (g *GroupMembers) IsInGroup() bool {
	sql := "SELECT COUNT(*) AS count FROM `group_members` WHERE `group_id` = ? AND `user_id` = ?;"
	var count int
	err := global.Database.DbHandle.QueryRow(sql, g.GroupId, g.UserId).Scan(&count)
	if err != nil {
		global.Logger.Errorf("err: %v", err)
	}
	return count > 0
}

func (g *GroupMembers) DeleteGroupMembers() {
	sql := "DELETE FROM `group_members` WHERE `group_id`=?"
	_, err := global.Database.DbHandle.Exec(sql, g.GroupId)
	if err != nil {
		global.Logger.Errorf("delete error: %v", err)
	}
}

func (g *GroupMembers) DeleteGroupMember() {
	sql := "DELETE FROM `group_members` WHERE `group_id` = ? AND `user_id` = ?;"
	_, err := global.Database.DbHandle.Exec(sql, g.GroupId, g.UserId)
	if err != nil {
		global.Logger.Errorf("delete error: %v", err)
	}
}

func (g *GroupMembers) Add() {
	sql := "INSERT INTO group_members (`group_id`, `user_id`, `status`) VALUES (?, ?, ?);"
	_, err := global.Database.DbHandle.Exec(sql, g.GroupId, g.UserId, g.Status)
	if err != nil {
		global.Logger.Errorf("err: %v", err)
	}
}

func (g *GroupMembers) BelongUserList() []*GroupMembers {
	sql := "SELECT * FROM group_members WHERE `user_id` = ?;"
	var gms []*GroupMembers
	rows, err := global.Database.DbHandle.Query(sql, g.UserId)
	if err != nil {
		global.Logger.Errorf("[TodoListDTO.List] %v", err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var temp GroupMembers
		err = rows.Scan(&temp.Id, &temp.GroupId, &temp.UserId, &temp.Status)
		if err != nil {
			global.Logger.Errorf("[TodoListDTO.Scan] %v", err)
		}
		gms = append(gms, &temp)
	}

	return gms
}

func (g *GroupMembers) BelongGroupList() []*GroupMembers {
	sql := "SELECT * FROM group_members WHERE `group_id` = ?;"
	var gms []*GroupMembers
	rows, err := global.Database.DbHandle.Query(sql, g.GroupId)
	if err != nil {
		global.Logger.Errorf("[TodoListDTO.List] %v", err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var temp GroupMembers
		err = rows.Scan(&temp.Id, &temp.GroupId, &temp.UserId, &temp.Status)
		if err != nil {
			global.Logger.Errorf("[TodoListDTO.Scan] %v", err)
		}
		gms = append(gms, &temp)
	}

	return gms
}
