// 文档地址： https://gorm.io/docs/

package mysqldb

import (
	"go-live-broadcast-downloader/plugins/db"
	"go-live-broadcast-downloader/plugins/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

const (
	DBInitErr = "DBInitErr"
)

type MySQLConfig struct {
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	//
	// 配置 eg: dev:passwd@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&interpolateParams=true&parseTime=true&timeout=200ms&writeTimeout=2s&readTimeout=2s
	DSN string `json:"DSN"`
	// 最大活跃链接数
	MaxOpenConns int `json:"MaxOpenConns"`
	// 最大空闲链接数
	MaxIdleConns int `json:"MaxIdleConns"`
	// 最大生命时间(秒)
	MaxLifetime int `json:"MaxLifetime"`
	// 最大空闲时间(秒)
	MaxIdleTime int `json:"MaxIdleTime"`
	// 日志级别
	LogLevel db.LogLevel `json:"LogLevel"`
	// 慢sql 定义（单位秒）
	SlowThreshold int `json:"SlowThreshold"`
}

func DB(c MySQLConfig) (*gorm.DB, error) {
	if c.MaxOpenConns <= 0 {
		c.MaxOpenConns = 64
	}
	if c.MaxIdleConns <= 0 {
		c.MaxIdleConns = 8
	}
	if c.MaxLifetime < 0 || c.MaxLifetime > 5*60 {
		c.MaxLifetime = 5 * 60
	}
	if c.MaxIdleTime <= 0 || c.MaxIdleTime > 5*60 {
		c.MaxIdleTime = 5 * 60
	}
	if c.LogLevel == 0 {
		c.LogLevel = db.LogLevelInfo
	}
	if c.SlowThreshold <= 0 {
		c.SlowThreshold = 1
	}

	mdb, err := gorm.Open(mysql.Open(c.DSN), &gorm.Config{
		Logger: db.NewDBLog(
			logger.Config{
				SlowThreshold:             time.Second,                 // 慢sql定义
				LogLevel:                  logger.LogLevel(c.LogLevel), // 日志级别
				IgnoreRecordNotFoundError: true,
				// 忽略找不到记录的err
			},
		),
	})
	if err != nil {
		log.Error(DBInitErr).Msg(err)
		return nil, err
	}
	sqldb, err := mdb.DB()
	if err != nil {
		log.Error(DBInitErr).Msg(err)
		return nil, err
	}

	sqldb.SetMaxOpenConns(c.MaxOpenConns)
	sqldb.SetMaxIdleConns(c.MaxIdleConns)
	sqldb.SetConnMaxLifetime(time.Duration(c.MaxLifetime) * time.Second)
	sqldb.SetConnMaxIdleTime(time.Duration(c.MaxIdleTime) * time.Second)
	go func() {
		timer := time.NewTicker(time.Minute)
		select {
		case <-timer.C:
			log.Debug("sqldb.Stats").Msgf("%+v", sqldb.Stats())
		}
	}()

	return mdb, nil
}
