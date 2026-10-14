package log

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"localhost.com/my_web/entity/public"
)

// log的外部接口
type LogInter interface {
	NewLogger(out io.Writer, prefix string, flag int)
	NewFileLogger(logPath string) error
	WriteLog(level int, v ...interface{})
	WriteError(v ...interface{})
	WriteInfo(v ...interface{})
	WriteInfoNotPre(v ...interface{})
}

// Logger 是一个自定义的日志记录器结构体
type logReal struct {
	logger *log.Logger
}

func GetLogInter() *logReal {
	return &logReal{}
}

// NewLogger 创建一个返回到标准输出的Logger实例
func (lr *logReal) NewLogger(out io.Writer, prefix string, flag int) {
	lr.logger = log.New(out, prefix, flag)
}

// NewFileLogger 创建一个将日志输出到指定文件的Logger实例
func (lr *logReal) NewFileLogger(logPath string) error {
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	lr.logger = log.New(file, "", 0)
	return nil
}

func (lr *logReal) WriteLog(level int, v ...interface{}) {
	if level == public.LOGGER_LEVEL_ERROR {
		lr.WriteError(v...)
	} else if level == public.LOGGER_LEVEL_INFO {
		lr.WriteInfo(v...)
	} else {
		lr.WriteInfoNotPre(v...)
	}
}

// 错误信息
func (lr *logReal) WriteError(v ...interface{}) {
	if lr != nil && lr.logger != nil {
		now := time.Now().Format("2006-01-02 15:04:05")
		lr.logger.Println("[ERROR]" + now + " " + fmt.Sprint(v...))
	}
}

// 实时信息
func (lr *logReal) WriteInfo(v ...interface{}) {
	if lr != nil && lr.logger != nil {
		now := time.Now().Format("2006-01-02 15:04:05")
		lr.logger.Println("[INFO]" + now + " " + fmt.Sprint(v...))
	}

}

// 实时信息
func (lr *logReal) WriteInfoNotPre(v ...interface{}) {
	if lr != nil && lr.logger != nil {
		lr.logger.Println(fmt.Sprint(v...))
	}
}
