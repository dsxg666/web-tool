package model

import (
	"github.com/dsxg666/web-tool/global"
	"github.com/dsxg666/web-tool/pkg/util"
)

type Posts struct {
	Id          string `json:"id"`
	UserId      string `json:"userId"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	Category    string `json:"category"`
	ViewCount   string `json:"viewCount"`
	IsPublic    string `json:"isPublic"`
	PublishedAt string `json:"publishedAt"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

type PostsDTO struct {
	Id       string `json:"id"`
	UserId   string `json:"userId"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	Category string `json:"category"`
}

func (p *PostsDTO) Update() {
	sql := "UPDATE `posts` SET `title` = ?, `category` =?, `content` = ? WHERE `id` = ?;"
	_, err := global.Database.DbHandle.Exec(sql, p.Title, p.Category, p.Content, p.Id)
	if err != nil {
		global.Logger.Errorf("Update error: %v", err)
	}
}

func (p *PostsDTO) Delete() {
	sql := "DELETE FROM `posts` WHERE `id`=?"
	_, err := global.Database.DbHandle.Exec(sql, p.Id)
	if err != nil {
		global.Logger.Errorf("Delete error: %v", err)
	}
}

func (p *PostsDTO) AddView() {
	sql := "UPDATE `posts` SET `view_count` = `view_count` + 1 WHERE `id` = ?;"
	_, err := global.Database.DbHandle.Exec(sql, p.Id)
	if err != nil {
		global.Logger.Errorf("AddView error: %v", err)
	}
}

func (p *PostsDTO) IsExistAndValid() bool {
	sql := "SELECT COUNT(*) AS count FROM `posts` WHERE `id` = ? AND `user_id` = ?;"
	var count int
	err := global.Database.DbHandle.QueryRow(sql, p.Id, p.UserId).Scan(&count)
	if err != nil {
		global.Logger.Errorf("err: %v", err)
	}
	return count > 0
}

func (p *PostsDTO) IsExistAndValid2() bool {
	sql := "SELECT COUNT(*) AS count FROM `posts` WHERE `id` = ? AND `is_public` = '1';"
	var count int
	err := global.Database.DbHandle.QueryRow(sql, p.Id).Scan(&count)
	if err != nil {
		global.Logger.Errorf("err: %v", err)
	}
	return count > 0
}

func (p *PostsDTO) GetById() *Posts {
	sql := "SELECT * FROM `posts` WHERE `id`=?"
	var temp Posts
	err := global.Database.DbHandle.QueryRow(sql, p.Id).Scan(&temp.Id, &temp.UserId, &temp.Title, &temp.Content, &temp.Category, &temp.ViewCount,
		&temp.IsPublic, &temp.PublishedAt, &temp.CreatedAt, &temp.UpdatedAt)
	if err != nil {
		global.Logger.Errorf("getById err: %v", err)
	}
	return &temp
}

func (p *PostsDTO) ToPublic() {
	sql := "UPDATE `posts` SET `is_public` = '1', `published_at` = ? WHERE `id` = ?;"
	_, err := global.Database.DbHandle.Exec(sql, util.GetNowFormatTime(), p.Id)
	if err != nil {
		global.Logger.Errorf("ToPublic error: %v", err)
	}
}

func (p *PostsDTO) ToPrivate() {
	sql := "UPDATE `posts` SET `is_public` = '0' WHERE `id` = ?;"
	_, err := global.Database.DbHandle.Exec(sql, p.Id)
	if err != nil {
		global.Logger.Errorf("ToPrivate error: %v", err)
	}
}

func (p *PostsDTO) Add() {
	sql := "INSERT INTO `posts` (`user_id`, `title`, `content`, `category`) VALUES (?, ?, ?, ?)"
	_, err := global.Database.DbHandle.Exec(sql, p.UserId, p.Title, p.Content, p.Category)
	if err != nil {
		global.Logger.Errorf("Postes add error: %v", err)
	}
}

func (p *PostsDTO) List(offset int) []*Posts {
	sql := "SELECT * FROM `posts` WHERE `is_public` = '1' ORDER BY `created_at` DESC Limit ?, 100;"
	rows, err := global.Database.DbHandle.Query(sql, offset)
	if err != nil {
		global.Logger.Errorf("Posts list %v", err)
		return nil
	}
	defer rows.Close()

	var ps []*Posts
	for rows.Next() {
		var temp Posts
		err = rows.Scan(&temp.Id, &temp.UserId, &temp.Title, &temp.Content, &temp.Category, &temp.ViewCount,
			&temp.IsPublic, &temp.PublishedAt, &temp.CreatedAt, &temp.UpdatedAt)
		if err != nil {
			global.Logger.Errorf("List scan %v", err)
		}
		ps = append(ps, &temp)
	}

	return ps
}

func (p *PostsDTO) MyList(offset int) []*Posts {
	sql := "SELECT * FROM `posts` WHERE `user_id` = ? ORDER BY `created_at` DESC Limit ?, 100;"
	rows, err := global.Database.DbHandle.Query(sql, p.UserId, offset)
	if err != nil {
		global.Logger.Errorf("Posts list %v", err)
		return nil
	}
	defer rows.Close()

	var ps []*Posts
	for rows.Next() {
		var temp Posts
		err = rows.Scan(&temp.Id, &temp.UserId, &temp.Title, &temp.Content, &temp.Category, &temp.ViewCount,
			&temp.IsPublic, &temp.PublishedAt, &temp.CreatedAt, &temp.UpdatedAt)
		if err != nil {
			global.Logger.Errorf("NyList scan %v", err)
		}
		ps = append(ps, &temp)
	}

	return ps
}

func (p *PostsDTO) GetListTotalCount() int {
	sql := "SELECT COUNT(*) AS count FROM `posts` WHERE `is_public` = '1';"
	var count int
	err := global.Database.DbHandle.QueryRow(sql).Scan(&count)
	if err != nil {
		global.Logger.Errorf("err: %v", err)
	}
	return count
}

func (p *PostsDTO) GetMyListTotalCount() int {
	sql := "SELECT COUNT(*) AS count FROM `posts` WHERE `user_id` = ?;"
	var count int
	err := global.Database.DbHandle.QueryRow(sql, p.UserId).Scan(&count)
	if err != nil {
		global.Logger.Errorf("err: %v", err)
	}
	return count
}

type Likes struct {
	Id        string `json:"id"`
	UserId    string `json:"userId"`
	PostId    string `json:"postId"`
	CreatedAt string `json:"createdAt"`
}

func (l *Likes) Add() {
	sql := "INSERT INTO `likes` (`user_id`, `post_id`) VALUES (?, ?)"
	_, err := global.Database.DbHandle.Exec(sql, l.UserId, l.PostId)
	if err != nil {
		global.Logger.Errorf("Add error: %v", err)
	}
}

func (l *Likes) Delete() {
	sql := "DELETE FROM `likes` WHERE `post_id` = ? AND `user_id` = ?;"
	_, err := global.Database.DbHandle.Exec(sql, l.PostId, l.UserId)
	if err != nil {
		global.Logger.Errorf("Delete error: %v", err)
	}
}

func (l *Likes) IsLike() bool {
	sql := "SELECT COUNT(*) AS count FROM `likes` WHERE `post_id` = ? AND `user_id` = ?;"
	var count int
	err := global.Database.DbHandle.QueryRow(sql, l.PostId, l.UserId).Scan(&count)
	if err != nil {
		global.Logger.Errorf("err: %v", err)
	}
	return count > 0
}

func (l *Likes) GetCount() int {
	sql := "SELECT COUNT(*) AS count FROM `likes` WHERE `post_id` = ?"
	var count int
	err := global.Database.DbHandle.QueryRow(sql, l.PostId).Scan(&count)
	if err != nil {
		global.Logger.Errorf("err: %v", err)
	}
	return count
}

type Favorites struct {
	Id        string `json:"id"`
	UserId    string `json:"userId"`
	PostId    string `json:"postId"`
	CreatedAt string `json:"createdAt"`
}

func (f *Favorites) List(offset int) []string {
	sql := "SELECT `post_id` FROM `favorites` WHERE `user_id` = ? ORDER BY `created_at` DESC Limit ?, 100;"
	rows, err := global.Database.DbHandle.Query(sql, f.UserId, offset)
	if err != nil {
		global.Logger.Errorf("List %v", err)
		return nil
	}
	defer rows.Close()

	var fs []string
	for rows.Next() {
		var temp string
		err = rows.Scan(&temp)
		if err != nil {
			global.Logger.Errorf("Scan %v", err)
		}
		fs = append(fs, temp)
	}

	return fs
}

func (f *Favorites) Add() {
	sql := "INSERT INTO `favorites` (`user_id`, `post_id`) VALUES (?, ?)"
	_, err := global.Database.DbHandle.Exec(sql, f.UserId, f.PostId)
	if err != nil {
		global.Logger.Errorf("Add error: %v", err)
	}
}

func (f *Favorites) Delete() {
	sql := "DELETE FROM `favorites` WHERE `post_id` = ? AND `user_id` = ?;"
	_, err := global.Database.DbHandle.Exec(sql, f.PostId, f.UserId)
	if err != nil {
		global.Logger.Errorf("Delete error: %v", err)
	}
}

func (f *Favorites) IsFavorite() bool {
	sql := "SELECT COUNT(*) AS count FROM `favorites` WHERE `post_id` = ? AND `user_id` = ?;"
	var count int
	err := global.Database.DbHandle.QueryRow(sql, f.PostId, f.UserId).Scan(&count)
	if err != nil {
		global.Logger.Errorf("err: %v", err)
	}
	return count > 0
}

func (f *Favorites) GetPostCount() int {
	sql := "SELECT COUNT(*) AS count FROM `favorites` WHERE `post_id` = ?"
	var count int
	err := global.Database.DbHandle.QueryRow(sql, f.PostId).Scan(&count)
	if err != nil {
		global.Logger.Errorf("err: %v", err)
	}
	return count
}

func (f *Favorites) GetUserCount() int {
	sql := "SELECT COUNT(f.post_id) AS count FROM favorites f JOIN posts p ON f.post_id = p.id WHERE f.user_id = ? AND p.is_public = '1';"
	var count int
	err := global.Database.DbHandle.QueryRow(sql, f.UserId).Scan(&count)
	if err != nil {
		global.Logger.Errorf("err: %v", err)
	}
	return count
}
