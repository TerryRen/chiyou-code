package lib

import (
	"path"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"
)

// Config rotate logger
func ConfigRotateLogger(logLevel string, logPath string, logFileName string, maxAge time.Duration, rotationTime time.Duration) {
	baseLogPath := path.Join(logPath, logFileName)
	writer, err := rotatelogs.New(
		baseLogPath+"-%Y%m%d%H%M.log",
		// 生成软链，指向最新日志文件
		rotatelogs.WithLinkName(baseLogPath),
		rotatelogs.WithMaxAge(maxAge),             // 文件最大保存时间
		rotatelogs.WithRotationTime(rotationTime), // 日志切割时间间隔
	)
	if err != nil {
		log.Errorf("config local file system logger error. %+v", errors.WithStack(err))
	}
	lfHook := lfshook.NewHook(lfshook.WriterMap{
		// 为不同级别设置不同的输出目的
		log.TraceLevel: writer,
		log.DebugLevel: writer,
		log.InfoLevel:  writer,
		log.WarnLevel:  writer,
		log.ErrorLevel: writer,
		log.FatalLevel: writer,
		log.PanicLevel: writer,
	}, &log.TextFormatter{
		DisableColors:   true,
		TimestampFormat: "2006-01-02 15:04:05.000000",
		FullTimestamp:   true,
	})
	// 设置日志等级 (默认Info)
	level, err := log.ParseLevel(logLevel)
	if err == nil {
		log.SetLevel(level)
	} else {
		log.SetLevel(log.InfoLevel)
	}
	// 将函数名和行数放在日志里面
	log.SetReportCaller(true)
	log.AddHook(lfHook)
}
