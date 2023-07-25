package clickhousedb

import (
	"go-live-broadcast-downloader/plugins/db"
	"go-live-broadcast-downloader/plugins/log"
	"gorm.io/driver/clickhouse"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

const (
	DBInitErr = "DBInitErr"
)

type ClickHouseConfig struct {
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

func DB(c ClickHouseConfig) (*gorm.DB, error) {
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

	mdb, err := gorm.Open(clickhouse.New(clickhouse.Config{
		DSN: c.DSN,
	}), &gorm.Config{
		Logger: db.NewDBLog(
			logger.Config{
				SlowThreshold:             time.Duration(c.SlowThreshold) * time.Second, // 慢sql定义
				LogLevel:                  logger.LogLevel(c.LogLevel),                  // 日志级别
				IgnoreRecordNotFoundError: true,                                         // 忽略找不到记录的err
			},
		),
	})

	sqldb, err := mdb.DB()
	if err != nil {
		log.Error(DBInitErr).Msg(err)
		return nil, err
	}
	sqldb.SetMaxOpenConns(c.MaxOpenConns)
	sqldb.SetMaxIdleConns(c.MaxIdleConns)
	sqldb.SetConnMaxLifetime(time.Duration(c.MaxLifetime) * time.Second)
	sqldb.SetConnMaxIdleTime(time.Duration(c.MaxIdleTime) * time.Second)
	return mdb, nil
}
