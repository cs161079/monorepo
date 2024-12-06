package logger

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	logger "github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
	gormlogger "gorm.io/gorm/logger"
)

type GormLogger struct {
	LogLevel  gormlogger.LogLevel
	Formatter *easy.Formatter
}

const (
	rootLogsDirpath = "/logs"
)

var Logger *logger.Logger

func InitLogger(applicationName string) {
	Logger = &logger.Logger{
		Out:   os.Stderr,
		Level: logger.InfoLevel,
		Formatter: &easy.Formatter{
			TimestampFormat: "2006-01-02 15:04:05",
			LogFormat:       "[%lvl%]: %time% - %msg%",
		},
	}
	directoryPath := filepath.Join(rootLogsDirpath, applicationName)
	err := os.Mkdir(directoryPath, 0777)
	if err != nil {
		fmt.Printf("error create directory file: %v\n", err)
	}
	fmt.Println("Folder create succefully for logs....")
	var runMode = os.Getenv("application.mode")
	if runMode == "PROD" {
		fileName := filepath.Join(rootLogsDirpath, applicationName, "oasaLogs.log")
		//open a file
		f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
		if err != nil {
			fmt.Printf("error opening file: %v\n", err)
		}
		Logger.SetOutput(f)
	}

}

func INFO(str string) {
	Logger.Println(fmt.Sprintf("%s\n", str))
}

func WARN(str string) {
	Logger.SetLevel(logger.DebugLevel)
	Logger.Warn(fmt.Sprintf("%s\n", str))
}

func ERROR(str string) {
	Logger.SetLevel(logger.ErrorLevel)
	Logger.Error(fmt.Sprintf("%s\n", str))
}

// LogMode set log mode
func (l *GormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	newlogger := *l
	newlogger.LogLevel = level
	return &newlogger
}

// Info prints info
func (l *GormLogger) Info(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel >= gormlogger.Info {
		// Logger.Infof(str, args...)
		logger.WithFields(logger.Fields{
			"at": time.Now().Format("2006-01-02 15:04:05"),
		}).Infof(str+"\n", args)
	}
}

// Warn prints warn messages
func (l *GormLogger) Warn(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel >= gormlogger.Warn {
		// Logger.Warnf(str, args...)
		logger.WithFields(logger.Fields{
			"at": time.Now().Format("2006-01-02 15:04:05"),
		}).Warnf(str+"\n", args)
	}

}

// Error prints error messages
func (l *GormLogger) Error(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel >= gormlogger.Error {
		// Logger.Errorf(str, args...)
		//logger.WithFields(logger.Fields{
		//	"at": time.Now().Format("2006-01-02 15:04:05"),
		//}).Errorf(str+"\n", args)
		logger.Error(str, args)
	}
}

// Trace prints trace messages
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= gormlogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	switch {
	case err != nil && l.LogLevel >= gormlogger.Error: //&& (!errors.Is(err, ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		sql, rows := fc()
		if rows == -1 {
			//l.Printf(l.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
			//Logger.Error("[", elapsed.Milliseconds(), " ms, ", "sql -> ", sql, "\n")
			l.Error(nil, fmt.Sprintf("[%d ms, sql -> %s", elapsed.Milliseconds(), sql))
		} else {
			//l.Printf(l.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
			//Logger.Error("[", elapsed.Milliseconds(), " ms, ", rows, " rows] ", "sql -> ", sql, "\n")
			l.Error(nil, fmt.Sprintf("[%d ms, %d rows, sql -> %s", elapsed.Milliseconds(), rows, sql))
		}
	case l.LogLevel >= gormlogger.Warn:
		sql, rows := fc()
		// slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
		if rows == -1 {
			//l.Printf(l.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
			Logger.Warn("[", elapsed.Milliseconds(), " ms, ", "sql -> ", sql)
		} else {
			//l.Printf(l.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
			Logger.Info("[", elapsed.Milliseconds(), " ms, ", rows, " rows] ", "sql -> ", sql)
		}
	case l.LogLevel == gormlogger.Info:
		sql, rows := fc()
		if rows == -1 {

			Logger.Info("[", elapsed.Milliseconds(), " ms, ", "sql -> ", sql)
			//l.Info(l.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			//l.Printf(l.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
			Logger.Info("[", elapsed.Milliseconds(), " ms, ", rows, " rows] ", "sql -> ", sql)
		}
	}
}

func GetGormLogger() *GormLogger {
	return &GormLogger{
		LogLevel: gormlogger.Error,
		Formatter: &easy.Formatter{
			TimestampFormat: "2006-01-02 15:04:05",
			LogFormat:       "[%lvl%]: %time% - %msg%",
		},
	}
}
