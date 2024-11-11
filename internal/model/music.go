package model

import (
	"fmt"
	"github.com/dsxg666/web-tool/global"
	"strings"
)

type Songs struct {
	Id     string `json:"id"`
	Title  string `json:"title"`
	Artist string `json:"artist"`
	Time   string `json:"time"`
}

func (s *Songs) ListByFavorite(ids []string) []*Songs {
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}
	sql := fmt.Sprintf("SELECT * FROM `songs` WHERE id IN (%s)", strings.Join(placeholders, ", "))
	rows, err := global.Database.DbHandle.Query(sql, args...)
	if err != nil {
		global.Logger.Errorf("ListByFavorite list %v", err)
		return nil
	}
	defer rows.Close()

	var ss []*Songs
	for rows.Next() {
		var temp Songs
		err = rows.Scan(&temp.Id, &temp.Title, &temp.Artist, &temp.Time)
		if err != nil {
			global.Logger.Errorf("Song scan %v", err)
		}
		ss = append(ss, &temp)
	}

	return ss
}

func (s *Songs) List() []*Songs {
	sql := "SELECT * FROM `songs`"
	rows, err := global.Database.DbHandle.Query(sql)
	if err != nil {
		global.Logger.Errorf("Song list %v", err)
		return nil
	}
	defer rows.Close()

	var ss []*Songs
	for rows.Next() {
		var temp Songs
		err = rows.Scan(&temp.Id, &temp.Title, &temp.Artist, &temp.Time)
		if err != nil {
			global.Logger.Errorf("Song scan %v", err)
		}
		ss = append(ss, &temp)
	}

	return ss
}

type SongFavorites struct {
	Id     string `json:"id"`
	UserId string `json:"user_id"`
	SongId string `json:"song_id"`
}

func (sf *SongFavorites) List() []string {
	sql := "SELECT `song_id` FROM song_favorites WHERE `user_id` = ?;"
	rows, err := global.Database.DbHandle.Query(sql, sf.UserId)
	if err != nil {
		global.Logger.Errorf("Song list %v", err)
		return nil
	}
	defer rows.Close()

	var ts []string
	for rows.Next() {
		var temp string
		err = rows.Scan(&temp)
		if err != nil {
			global.Logger.Errorf("Song scan %v", err)
		}
		ts = append(ts, temp)
	}

	return ts
}

func (sf *SongFavorites) Add() {
	sql := "INSERT INTO `song_favorites` (`song_id`, `user_id`) VALUES (?, ?)"
	_, err := global.Database.DbHandle.Exec(sql, sf.SongId, sf.UserId)
	if err != nil {
		global.Logger.Errorf("Add error: %v", err)
	}
}

func (sf *SongFavorites) Delete() {
	sql := "DELETE FROM `song_favorites` WHERE `song_id` = ? AND `user_id` = ?;"
	_, err := global.Database.DbHandle.Exec(sql, sf.SongId, sf.UserId)
	if err != nil {
		global.Logger.Errorf("Delete error: %v", err)
	}
}
