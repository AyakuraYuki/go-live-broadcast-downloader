package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/opentracing/opentracing-go/ext"
	"go-live-broadcast-downloader/plugins/log"
	"go-live-broadcast-downloader/plugins/trace"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"strings"
	"time"
)

const (
	dbLogName = "DBLog"
	dbSQLSlow = "DBSQLSlow"
)

type LogLevel int

const (
	LogLevelSilent LogLevel = iota + 1
	LogLevelError
	LogLevelWarn
	LogLevelInfo
)

type dbLog struct {
	logger.Config
	infoStr, warnStr, errStr, traceStr, traceErrStr, traceWarnStr string
}

func NewDBLog(config logger.Config) logger.Interface {
	var (
		infoStr      = "%s "
		warnStr      = "%s "
		errStr       = "%s "
		traceStr     = "%s [%.3fms] [rows:%v] %s"
		traceWarnStr = "%s %s [%.3fms] [rows:%v] %s"
		traceErrStr  = "%s %s [%.3fms] [rows:%v] %s"
	)

	return &dbLog{
		Config:       config,
		infoStr:      infoStr,
		warnStr:      warnStr,
		errStr:       errStr,
		traceStr:     traceStr,
		traceWarnStr: traceWarnStr,
		traceErrStr:  traceErrStr,
	}
}

// LogMode log mode
func (l *dbLog) LogMode(level logger.LogLevel) logger.Interface {
	newlogger := *l
	newlogger.LogLevel = level
	return &newlogger
}

// Info print info
func (l dbLog) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Info {
		log.Info(dbLogName).Msgf(l.infoStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Warn print warn messages
func (l dbLog) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Warn {
		log.Warn(dbLogName).Msgf(l.warnStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Error print error messages
func (l dbLog) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Error {
		log.Error(dbLogName).Msgf(l.errStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Trace print sql message
func (l dbLog) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	codeLine := utils.FileWithLineNum()
	if arr := strings.Split(codeLine, "/business/"); len(arr) == 2 {
		codeLine = "/business/" + arr[1]
	}
	span, _ := trace.StartSpanFromContextWithSt(ctx, codeLine, begin)
	defer span.Finish()
	span.SetTag(string(ext.DBStatement), sql)
	span.SetTag(string(ext.DBType), "sql")

	switch {
	case err != nil && l.LogLevel >= logger.Error && (!errors.Is(err, logger.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		if rows == -1 {
			log.Error(dbLogName).Msgf(l.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			log.Error(dbLogName).Msgf(l.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= logger.Warn:
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
		if rows == -1 {
			log.Warn(dbSQLSlow).Msgf(l.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			log.Warn(dbSQLSlow).Msgf(l.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case l.LogLevel == logger.Info:
		if rows == -1 {
			log.Info(dbLogName).Msgf(l.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			log.Info(dbLogName).Msgf(l.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	}
}
