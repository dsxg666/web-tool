package setting

import "time"

type ServerSetting struct {
	RunMode      string
	HttpPort     string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	FrontendHost string
}

type AppSetting struct {
	LogSavePath      string
	LogFileName      string
	LogFileExtension string
}

type MySQLSetting struct {
	DBType       string
	UserName     string
	Password     string
	Host         string
	DBName       string
	Charset      string
	ParseTime    bool
	MaxIdleConns int
	MaxOpenConns int
}

type EmailSetting struct {
	Host     string
	Port     int
	UserName string
	Password string
	IsSSL    bool
	From     string
}

type JwtTokenSetting struct {
	SecretKey      string
	ExpirationTime time.Duration
}

type RedisSetting struct {
	Host     string
	Password string
	DB       int
}
