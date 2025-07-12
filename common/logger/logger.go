package logger

import (
	"context"
	"fmt"
	"io"
	"os"
	"runtime"
	"time"

	"log"

	"github.com/cloudwego/kitex/pkg/klog"
	amqp "github.com/rabbitmq/amqp091-go"
)

type MyLogger struct {
	channel *amqp.Channel
	queue   *amqp.Queue
	prefix  string
	ctx     context.Context
	level   klog.Level
}

func NewLogger(channel *amqp.Channel, queue *amqp.Queue, prefix string, ctx context.Context) *MyLogger {
	return &MyLogger{
		channel: channel,
		queue:   queue,
		prefix:  prefix,
		ctx:     ctx,
		level:   klog.LevelInfo,
	}
}

func (logger *MyLogger) SetLevel(level klog.Level) {
	logger.level = level
}

// 因为日志已经确定要送往消息队列，因此这个接口弃用
func (logger *MyLogger) SetOutput(io.Writer) {
	return
}

func formatHeader(now time.Time, prefix string, file string, line int, level string) string {
	return fmt.Sprintf("%s %s: %s:%d: %s ", now.Format("2006/01/02 15:04:05.000000"), prefix, file, line, level)
}

func (logger MyLogger) logf(level klog.Level, format *string, v ...interface{}) {
	if logger.level > level {
		return
	}

	// 这一段是从标准库的log中抄过来的
	now := time.Now()
	const calldepth = 4
	_, file, line, ok := runtime.Caller(calldepth)
	if !ok {
		file = "???"
		line = 0
	}

	msg := formatHeader(now, logger.prefix, file, line, toString(level))

	if format != nil {
		msg += fmt.Sprintf(*format, v...)
	} else {
		msg += fmt.Sprint(v...)
	}
	msg += "\n"

	var err error
	if logger.ctx == nil {
		err = logger.channel.Publish(
			"",                // exchange
			logger.queue.Name, // routing key
			false,             // mandatory
			false,             // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(msg),
			},
		)
	} else {
		err = logger.channel.PublishWithContext(logger.ctx,
			"",                // exchange
			logger.queue.Name, // routing key
			false,             // mandatory
			false,             // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(msg),
			},
		)
	}
	if err != nil {
		log.Fatal(err)
	}
	if level == klog.LevelFatal {
		os.Exit(1)
	}
}

func (logger *MyLogger) Close() {
	logger.channel.Close()
}

func (logger *MyLogger) Trace(v ...interface{}) {
	logger.logf(klog.LevelTrace, nil, v...)
}

func (logger *MyLogger) Debug(v ...interface{}) {
	logger.logf(klog.LevelDebug, nil, v...)
}

func (logger *MyLogger) Info(v ...interface{}) {
	logger.logf(klog.LevelInfo, nil, v...)
}

func (logger *MyLogger) Notice(v ...interface{}) {
	logger.logf(klog.LevelNotice, nil, v...)
}

func (logger *MyLogger) Warn(v ...interface{}) {
	logger.logf(klog.LevelWarn, nil, v...)
}

func (logger *MyLogger) Error(v ...interface{}) {
	logger.logf(klog.LevelError, nil, v...)
}

func (logger *MyLogger) Fatal(v ...interface{}) {
	logger.logf(klog.LevelFatal, nil, v...)
}

func (logger *MyLogger) Tracef(format string, v ...interface{}) {
	logger.logf(klog.LevelTrace, &format, v...)
}

func (logger *MyLogger) Debugf(format string, v ...interface{}) {
	logger.logf(klog.LevelDebug, &format, v...)
}

func (logger *MyLogger) Infof(format string, v ...interface{}) {
	logger.logf(klog.LevelInfo, &format, v...)
}

func (logger *MyLogger) Noticef(format string, v ...interface{}) {
	logger.logf(klog.LevelNotice, &format, v...)
}

func (logger *MyLogger) Warnf(format string, v ...interface{}) {
	logger.logf(klog.LevelWarn, &format, v...)
}

func (logger *MyLogger) Errorf(format string, v ...interface{}) {
	logger.logf(klog.LevelError, &format, v...)
}

func (logger *MyLogger) Fatalf(format string, v ...interface{}) {
	logger.logf(klog.LevelFatal, &format, v...)
}

// 这边仿照klog.defaultLogger的做法
// 因为没有放ctx的地方，就直接当没有
func (logger *MyLogger) CtxTracef(ctx context.Context, format string, v ...interface{}) {
	logger.Tracef(format, v...)
}

func (logger *MyLogger) CtxDebugf(ctx context.Context, format string, v ...interface{}) {
	logger.Debugf(format, v...)
}

func (logger *MyLogger) CtxInfof(ctx context.Context, format string, v ...interface{}) {
	logger.Infof(format, v...)
}

func (logger *MyLogger) CtxNoticef(ctx context.Context, format string, v ...interface{}) {
	logger.Noticef(format, v...)
}

func (logger *MyLogger) CtxWarnf(ctx context.Context, format string, v ...interface{}) {
	logger.Warnf(format, v...)
}

func (logger *MyLogger) CtxErrorf(ctx context.Context, format string, v ...interface{}) {
	logger.Errorf(format, v...)
}

func (logger *MyLogger) CtxFatalf(ctx context.Context, format string, v ...interface{}) {
	logger.Fatalf(format, v...)
}

var strs = []string{
	"[Trace] ",
	"[Debug] ",
	"[Info] ",
	"[Notice] ",
	"[Warn] ",
	"[Error] ",
	"[Fatal] ",
}

func toString(level klog.Level) string {
	if level >= klog.LevelTrace && level <= klog.LevelFatal {
		return strs[level]
	}
	return fmt.Sprintf("[?%d] ", level)
}

// type FormatLogger interface {
// 	Tracef(format string, v ...interface{})
// 	Debugf(format string, v ...interface{})
// 	Infof(format string, v ...interface{})
// 	Noticef(format string, v ...interface{})
// 	Warnf(format string, v ...interface{})
// 	Errorf(format string, v ...interface{})
// 	Fatalf(format string, v ...interface{})
// }

// type Logger interface {
// 	Trace(v ...interface{})
// 	Debug(v ...interface{})
// 	Info(v ...interface{})
// 	Notice(v ...interface{})
// 	Warn(v ...interface{})
// 	Error(v ...interface{})
// 	Fatal(v ...interface{})
// }

// type CtxLogger interface {
// 	CtxTracef(ctx context.Context, format string, v ...interface{})
// 	CtxDebugf(ctx context.Context, format string, v ...interface{})
// 	CtxInfof(ctx context.Context, format string, v ...interface{})
// 	CtxNoticef(ctx context.Context, format string, v ...interface{})
// 	CtxWarnf(ctx context.Context, format string, v ...interface{})
// 	CtxErrorf(ctx context.Context, format string, v ...interface{})
// 	CtxFatalf(ctx context.Context, format string, v ...interface{})
// }

// type Control interface {
// 	SetLevel(Level)
// 	SetOutput(io.Writer)
// }
