package model

import "github.com/dsxg666/web-tool/global"

type Code struct {
	Id        string
	Email     string
	Code      string
	Type      string
	CreatedAt string
}

type CodeDTO struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

func (c *CodeDTO) AddRegisterCode() {
	sql := "INSERT INTO codes (`email`, `code`, `type`) VALUES (?, ?, ?);"
	_, err := global.Database.DbHandle.Exec(sql, c.Email, c.Code, "1")
	if err != nil {
		global.Logger.Errorf("err: %v", err)
	}
}

func (c *CodeDTO) AddVerificationCode() {
	sql := "INSERT INTO codes (`email`, `code`) VALUES (?, ?);"
	_, err := global.Database.DbHandle.Exec(sql, c.Email, c.Code)
	if err != nil {
		global.Logger.Errorf("err: %v", err)
	}
}

func (c *CodeDTO) IsValidRegisterCode() bool {
	sql := "SELECT COUNT(*) AS count FROM codes WHERE `email` = ? AND `code` = ? AND `type` = ? AND `created_at` >= NOW() - INTERVAL 5 MINUTE;"
	var count int
	err := global.Database.DbHandle.QueryRow(sql, c.Email, c.Code, "1").Scan(&count)
	if err != nil {
		global.Logger.Errorf("IsValidRegisterCode err: %v", err)
	}
	return count > 0
}

func (c *CodeDTO) IsValidVerificationCode() bool {
	sql := "SELECT COUNT(*) AS count FROM codes WHERE email = ? AND code = ? AND `type` = ?  AND created_at >= NOW() - INTERVAL 5 MINUTE;"
	var count int
	err := global.Database.DbHandle.QueryRow(sql, c.Email, c.Code, "0").Scan(&count)
	if err != nil {
		global.Logger.Errorf("err: %v", err)
	}
	return count > 0
}

func (c *CodeDTO) IsTooQuickRegister() bool {
	sql := "SELECT COUNT(*) AS count FROM codes WHERE email = ? AND `type` = ? AND created_at >= NOW() - INTERVAL 1 MINUTE;"
	var count int
	err := global.Database.DbHandle.QueryRow(sql, c.Email, "1").Scan(&count)
	if err != nil {
		global.Logger.Errorf("err: %v", err)
	}
	return count > 0
}

func (c *CodeDTO) IsTooQuickVerification() bool {
	sql := "SELECT COUNT(*) AS count FROM codes WHERE email = ? AND `type` = ? AND created_at >= NOW() - INTERVAL 1 MINUTE;"
	var count int
	err := global.Database.DbHandle.QueryRow(sql, c.Email, "0").Scan(&count)
	if err != nil {
		global.Logger.Errorf("err: %v", err)
	}
	return count > 0
}
