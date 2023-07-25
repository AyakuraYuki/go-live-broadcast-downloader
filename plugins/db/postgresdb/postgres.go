package postgresdb

import (
	"go-live-broadcast-downloader/plugins/db"
	"go-live-broadcast-downloader/plugins/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

const (
	DBInitErr = "DBInitErr"
)

type PostgresConfig struct {
	// https://github.com/jackc/pgx
	// dsn := "host=localhost user=gorm password=gorm dbname=db_gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"

	DSN           string      `json:"DSN"`
	MaxOpenConns  int         `json:"MaxOpenConns"`  // maximum open connections
	MaxIdleConns  int         `json:"MaxIdleConns"`  // maximum idle connections
	MaxLifetime   int         `json:"MaxLifetime"`   // maximum lifetime in seconds
	MaxIdleTime   int         `json:"MaxIdleTime"`   // maximum idle time in seconds
	LogLevel      db.LogLevel `json:"LogLevel"`      // log level
	SlowThreshold int         `json:"SlowThreshold"` // defined execution time for slow sql, in seconds
}

func DB(c *PostgresConfig) (*gorm.DB, error) {
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

	mdb, err := gorm.Open(postgres.Open(c.DSN), &gorm.Config{

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
