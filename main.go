package main

import (
	"database/sql"
	"fmt"
	"github.com/dsxg666/web-tool/global"
	"github.com/dsxg666/web-tool/internal/routers"
	"github.com/dsxg666/web-tool/internal/ws"
	"github.com/dsxg666/web-tool/pkg/db"
	"github.com/dsxg666/web-tool/pkg/logger"
	"github.com/dsxg666/web-tool/pkg/setting"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"net/http"
	"time"
)

func init() {
	err := setupSetting()
	if err != nil {
		log.Fatalf("init.setupSetting err: %v", err)
	}
	err = setupLogger()
	if err != nil {
		log.Fatalf("init.setupLogger err: %v", err)
	}
	err = setupDatabase()
	if err != nil {
		log.Fatalf("init.setupDatabase err: %v", err)
	}
	global.OnlineUser = make([]string, 0)
}

// @title Tool API
// @version 1.0
// @description This is a sample server celler server.
func main() {
	serverSetting := global.ServerSetting
	gin.SetMode(serverSetting.RunMode)
	hub := ws.NewHub()
	go hub.Run()
	router := routers.NewRouter()
	router.GET("/ws", func(c *gin.Context) {
		ws.ServeWs(hub, c)
	})
	s := &http.Server{
		Addr:           ":" + serverSetting.HttpPort,
		Handler:        router,
		ReadTimeout:    serverSetting.ReadTimeout * time.Second,
		WriteTimeout:   serverSetting.WriteTimeout * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}

func setupSetting() error {
	s, err := setting.NewSetting()
	if err != nil {
		return err
	}
	err = s.ReadSection("Server", &global.ServerSetting)
	if err != nil {
		return err
	}
	err = s.ReadSection("App", &global.AppSetting)
	if err != nil {
		return err
	}
	err = s.ReadSection("Database", &global.DatabaseSetting)
	if err != nil {
		return err
	}
	err = s.ReadSection("Email", &global.EmailSetting)
	if err != nil {
		return err
	}
	err = s.ReadSection("JwtToken", &global.JwtTokenSetting)
	if err != nil {
		return err
	}
	return nil
}

func setupLogger() error {
	appSetting := global.AppSetting
	global.Logger = logger.NewLogger(&lumberjack.Logger{
		Filename:  appSetting.LogSavePath + "/" + appSetting.LogFileName + appSetting.LogFileExtension,
		MaxSize:   600,
		MaxAge:    10,
		LocalTime: true,
	}, "", log.LstdFlags).WithCaller(2)

	return nil
}

func setupDatabase() error {
	databaseSetting := global.DatabaseSetting
	dbHandle, err := sql.Open(databaseSetting.DBType, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=%t&loc=Local",
		databaseSetting.UserName,
		databaseSetting.Password,
		databaseSetting.Host,
		databaseSetting.DBName,
		databaseSetting.Charset,
		databaseSetting.ParseTime,
	))
	if err != nil {
		return err
	}

	dbHandle.SetMaxOpenConns(global.DatabaseSetting.MaxOpenConns)
	dbHandle.SetMaxIdleConns(global.DatabaseSetting.MaxIdleConns)

	global.Database = &db.Database{
		DbHandle: dbHandle,
	}

	return nil
}
