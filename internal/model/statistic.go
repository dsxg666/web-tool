package model

import (
	"github.com/dsxg666/web-tool/global"
	"github.com/dsxg666/web-tool/pkg/util"
)

type Dau struct {
	Id        string `json:"id"`
	UserId    string `json:"userId"`
	UserIp    string `json:"userIp"`
	CreatedAt string `json:"createdAt"`
}

func (d *Dau) GetPastWeekData(arr []string) []int {
	datas := make([]int, 7)
	sql := "SELECT DATE(created_at) AS date, COUNT(DISTINCT `user_id`) AS active_users FROM `dau` WHERE `created_at` >= NOW() - INTERVAL 7 DAY GROUP BY DATE(`created_at`) ORDER BY DATE(`created_at`) DESC;"
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

func (d *Dau) Add() {
	sql := "INSERT INTO `dau` (`user_id`,`user_ip`) VALUES (?, ?);"
	_, err := global.Database.DbHandle.Exec(sql, d.UserId, d.UserIp)
	if err != nil {
		global.Logger.Errorf(" add error: %v", err)
	}
}

func (d *Dau) GetTodayDauNum() int {
	sql := "SELECT COUNT(DISTINCT `user_id`) AS dau_num FROM `dau` WHERE DATE(`created_at`) = CURDATE();"
	var count int
	err := global.Database.DbHandle.QueryRow(sql).Scan(&count)
	if err != nil {
		global.Logger.Errorf(" GetTodayDauNum error: %v", err)
	}
	return count
}

func (d *Dau) GetYesterdayDauNum() int {
	sql := "SELECT COUNT(DISTINCT `user_id`) AS dau_num FROM `dau` WHERE DATE(`created_at`) = DATE_SUB(CURDATE(), INTERVAL 1 DAY);"
	var count int
	err := global.Database.DbHandle.QueryRow(sql).Scan(&count)
	if err != nil {
		global.Logger.Errorf("GetYesterdayDauNum error: %v", err)
	}
	return count
}
