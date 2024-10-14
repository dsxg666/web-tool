package model

import (
	"database/sql"
	"github.com/dsxg666/web-tool/global"
	"github.com/dsxg666/web-tool/pkg/util"
)

type TodoList struct {
	Id          string `json:"id"`
	UserId      string `json:"userId"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
	Status      string `json:"status"`
	DueDate     string `json:"dueDate"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

type TodoIdDTO struct {
	Id string `json:"id"`
}

type TodoListDTO struct {
	UserID    string `json:"userId"`
	CreatedAt string `json:"createdAt"`
}

type TodoAddDTO struct {
	UserID      string `json:"userId"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
	Status      string `json:"status"`
	DueDate     string `json:"dueDate"`
}

type TodoEditDTO struct {
	Id          string `json:"id"`
	UserID      string `json:"userId"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
	Status      string `json:"status"`
	DueDate     string `json:"dueDate"`
}

func (t *TodoList) GetPastWeekData(arr []string) []int {
	datas := make([]int, 7)
	sqlStat := "SELECT DATE(created_at) AS date, COUNT(`id`) AS num FROM `todolist` WHERE `created_at` >= NOW() - INTERVAL 7 DAY GROUP BY DATE(`created_at`) ORDER BY DATE(`created_at`) DESC;"
	rows, err := global.Database.DbHandle.Query(sqlStat)
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

func (t *TodoList) GetTodayTodolistNum() int {
	sqlStat := "SELECT COUNT(`id`) AS num FROM `todolist` WHERE DATE(`created_at`) = CURDATE();"
	var count int
	err := global.Database.DbHandle.QueryRow(sqlStat).Scan(&count)
	if err != nil {
		global.Logger.Errorf(" GetTodayTodolistNum error: %v", err)
	}
	return count
}

func (t *TodoList) GetYesterdayTodolistNum() int {
	sqlStat := "SELECT COUNT(`id`) AS num FROM `todolist` WHERE DATE(`created_at`) = DATE_SUB(CURDATE(), INTERVAL 1 DAY);"
	var count int
	err := global.Database.DbHandle.QueryRow(sqlStat).Scan(&count)
	if err != nil {
		global.Logger.Errorf(" GetYesterdayTodolistNum error: %v", err)
	}
	return count
}

func (t *TodoEditDTO) Edit() {
	sqlStat := "UPDATE `todolist` SET `title` = ?, `description` = ?, `priority` = ?, `status` = ?, `due_date` = ? WHERE `id` = ?;"
	_, err := global.Database.DbHandle.Exec(sqlStat, t.Title, t.Description, t.Priority, t.Status, t.DueDate, t.Id)
	if err != nil {
		global.Logger.Errorf("[todolist] edit error: %v", err)
	}
}

func (t *TodoEditDTO) IsExistAndValid() bool {
	sqlStat := "SELECT COUNT(*) AS count FROM `todolist` WHERE `id` = ? AND `user_id` = ?;"
	var count int
	err := global.Database.DbHandle.QueryRow(sqlStat, t.Id, t.UserID).Scan(&count)
	if err != nil {
		global.Logger.Errorf("err: %v", err)
	}
	return count > 0
}

func (t *TodoIdDTO) GetById() *TodoList {
	sqlStat := "SELECT * FROM `todolist` WHERE `id`=?"
	var temp TodoList
	err := global.Database.DbHandle.QueryRow(sqlStat, t.Id).Scan(&temp.Id, &temp.UserId, &temp.Title, &temp.Description,
		&temp.Priority, &temp.Status, &temp.DueDate, &temp.CreatedAt, &temp.UpdatedAt)
	if err != nil {
		global.Logger.Errorf("get todo list by id err: %v", err)
	}
	return &temp
}

func (t *TodoIdDTO) Delete() {
	sqlStat := "DELETE FROM todolist WHERE `id`=?"
	_, err := global.Database.DbHandle.Exec(sqlStat, t.Id)
	if err != nil {
		global.Logger.Errorf("[todolist] delete error: %v", err)
	}
}

func (t *TodoAddDTO) Add() {
	sqlStat := "INSERT INTO todolist (`user_id`, `title`, `description`, `priority`, `status`, `due_date`) VALUES (?, ?, ?, ?, ?, ?)"
	_, err := global.Database.DbHandle.Exec(sqlStat, t.UserID, t.Title, t.Description, t.Priority, t.Status, t.DueDate)
	if err != nil {
		global.Logger.Errorf("[todolist] add error: %v", err)
	}
}

func (t *TodoListDTO) List() []*TodoList {
	sqlStat := "SELECT * FROM todolist WHERE `user_id` = ? AND DATE(`created_at`) = ?;"
	var todos []*TodoList
	rows, err := global.Database.DbHandle.Query(sqlStat, t.UserID, t.CreatedAt)
	if err != nil {
		global.Logger.Errorf("[TodoListDTO.List] %v", err)
		return nil
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			global.Logger.Error(err.Error())
		}
	}(rows)

	for rows.Next() {
		var temp TodoList
		err = rows.Scan(&temp.Id, &temp.UserId, &temp.Title, &temp.Description,
			&temp.Priority, &temp.Status, &temp.DueDate, &temp.CreatedAt, &temp.UpdatedAt)
		if err != nil {
			global.Logger.Errorf("[TodoListDTO.Scan] %v", err)
		}
		todos = append(todos, &temp)
	}

	return todos
}
