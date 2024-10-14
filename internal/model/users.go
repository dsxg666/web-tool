package model

import (
	"github.com/dsxg666/web-tool/global"
	"github.com/dsxg666/web-tool/pkg/encrypt"
	"strconv"
)

type User struct {
	Id            string `json:"id"`
	Username      string `json:"username"`
	Password      string `json:"password"`
	Email         string `json:"email"`
	Avatar        string `json:"avatar"`
	EmailUpdateAt string `json:"emailUpdateAt"`
	CreatedAt     string `json:"createdAt"`
	UpdatedAt     string `json:"updatedAt"`
}

type UserRegisterDTO struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Code     string `json:"code"`
}

type UserLoginByPasswordDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserLoginByCodeDTO struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

type UserIdDTO struct {
	Id string `json:"id"`
}

type UserAvatarDTO struct {
	Id     string `json:"id"`
	Avatar string `json:"avatar"`
}

type UserModifyDTO struct {
	Id            string `json:"id"`
	Email         string `json:"email"`
	EmailUpdateAt string `json:"emailUpdateAt"`
	Username      string `json:"username"`
	Password      string `json:"password"`
	Avatar        string `json:"avatar"`
}

func (u *UserModifyDTO) ModifyAvatar() {
	sql := "UPDATE `users` SET `avatar` = ? WHERE `id` = ?;"
	_, err := global.Database.DbHandle.Exec(sql, u.Avatar, u.Id)
	if err != nil {
		global.Logger.Errorf("[user] modify avatar error: %v", err)
	}
}

func (u *UserModifyDTO) ModifyEmail() {
	sql := "UPDATE `users` SET `email` = ?, `email_update_at` = ? WHERE `id` = ?;"
	_, err := global.Database.DbHandle.Exec(sql, u.Email, u.EmailUpdateAt, u.Id)
	if err != nil {
		global.Logger.Errorf("[user] modify email error: %v", err)
	}
}

func (u *UserModifyDTO) ModifyUsername() {
	sql := "UPDATE `users` SET `username` = ? WHERE `id` = ?;"
	_, err := global.Database.DbHandle.Exec(sql, u.Username, u.Id)
	if err != nil {
		global.Logger.Errorf("[user] modify username error: %v", err)
	}
}

func (u *UserModifyDTO) ModifyPassword() {
	sql := "UPDATE `users` SET `password` = ? WHERE `id` = ?;"
	_, err := global.Database.DbHandle.Exec(sql, u.Password, u.Id)
	if err != nil {
		global.Logger.Errorf("[user] modify password error: %v", err)
	}
}

func (u *UserModifyDTO) IsEmailExist() bool {
	sql := "SELECT COUNT(*) AS count FROM `users` WHERE `email` = ?;"
	var count int
	err := global.Database.DbHandle.QueryRow(sql, u.Email).Scan(&count)
	if err != nil {
		global.Logger.Errorf("err: %v", err)
	}
	return count > 0
}

func (u *UserIdDTO) GetById() *User {
	sql := "SELECT * FROM `users` WHERE `id` = ?;"
	var user User
	err := global.Database.DbHandle.QueryRow(sql, u.Id).Scan(&user.Id, &user.Username, &user.Password, &user.Email,
		&user.Avatar, &user.EmailUpdateAt, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		global.Logger.Errorf("err: %v", err)
	}
	return &user
}

func (u *UserIdDTO) IsIdExist() bool {
	sql := "SELECT COUNT(*) AS count FROM `users` WHERE `id` = ?;"
	var count int
	err := global.Database.DbHandle.QueryRow(sql, u.Id).Scan(&count)
	if err != nil {
		global.Logger.Errorf("err: %v", err)
	}
	return count > 0
}

func (u *UserAvatarDTO) GetUser() *User {
	sql := "SELECT `avatar` FROM `users` WHERE `id` = ?;"
	var user User
	err := global.Database.DbHandle.QueryRow(sql, u.Id).Scan(&user.Avatar)
	if err != nil {
		global.Logger.Errorf("err: %v", err)
	}
	return &user
}

func (u *UserRegisterDTO) Add() string {
	sql := "INSERT INTO `users` (`username`, `email`, `password`) VALUES (?, ?, ?);"
	res, err := global.Database.DbHandle.Exec(sql, u.Username, u.Email, u.Password)
	if err != nil {
		global.Logger.Errorf("err: %v", err)
	}
	id, _ := res.LastInsertId()
	return strconv.Itoa(int(id))
}

func (u *UserRegisterDTO) IsEmailExist() bool {
	sql := "SELECT COUNT(*) AS count FROM `users` WHERE `email` = ?;"
	var count int
	err := global.Database.DbHandle.QueryRow(sql, u.Email).Scan(&count)
	if err != nil {
		global.Logger.Errorf("err: %v", err)
	}
	return count > 0
}

func (u *UserLoginByCodeDTO) IsEmailExist() bool {
	sql := "SELECT COUNT(*) AS count FROM `users` WHERE `email` = ?;"
	var count int
	err := global.Database.DbHandle.QueryRow(sql, u.Email).Scan(&count)
	if err != nil {
		global.Logger.Errorf("err: %v", err)
	}
	return count > 0
}

func (u *UserLoginByPasswordDTO) IsCorrectPassword() bool {
	sql := "SELECT `password` FROM `users` WHERE `email` = ?;"
	var pass string
	err := global.Database.DbHandle.QueryRow(sql, u.Email).Scan(&pass)
	if err != nil {
		global.Logger.Errorf("err: %v", err)
	}
	return encrypt.CheckPasswordHash(u.Password, pass)
}

func (u *UserLoginByPasswordDTO) IsEmailExist() bool {
	sql := "SELECT COUNT(*) AS count FROM `users` WHERE `email` = ?;"
	var count int
	err := global.Database.DbHandle.QueryRow(sql, u.Email).Scan(&count)
	if err != nil {
		global.Logger.Errorf("err: %v", err)
	}
	return count > 0
}

func (u *UserLoginByPasswordDTO) GetUser() *User {
	sql := "SELECT `id`, `username` FROM `users` WHERE `email` = ?;"
	var user User
	err := global.Database.DbHandle.QueryRow(sql, u.Email).Scan(&user.Id, &user.Username)
	if err != nil {
		global.Logger.Errorf("err: %v", err)
	}
	return &user
}

func (u *UserLoginByCodeDTO) GetUser() *User {
	sql := "SELECT `id`, `username` FROM `users` WHERE `email` = ?;"
	var user User
	err := global.Database.DbHandle.QueryRow(sql, u.Email).Scan(&user.Id, &user.Username)
	if err != nil {
		global.Logger.Errorf("err: %v", err)
	}
	return &user
}
