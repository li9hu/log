package log

import (
	"fmt"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

type LogConf struct {
	ConsoleEnable bool
	Color         bool
	FileEnable    bool
	FormatJson    bool
	Level         string
	Filename      string
	MaxFileSize   int
	MaxFileBackup int
}

func SetUp(l *LogConf) {
	if l.ConsoleEnable != true && l.FileEnable != true {
		return
	}
	var level zapcore.Level
	switch l.Level {
	case "info":
		level = zapcore.InfoLevel
	case "debug":
		level = zapcore.DebugLevel
	case "warn":
		level = zapcore.WarnLevel
	case "err":
		level = zapcore.ErrorLevel
	case "panic":
		level = zapcore.DPanicLevel
	default:
		level = zapcore.InfoLevel
	}
	l.setLogs(level)
}

func (l *LogConf) setLogs(level zapcore.Level) {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = timeEncoder // 设置自定义 TimeEncoder
	if l.FileEnable == false && l.ConsoleEnable == true && l.FormatJson == false && l.Color == true {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	var encoder zapcore.Encoder
	if l.FormatJson == true {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	var zapWriteSyncer zapcore.WriteSyncer
	if l.FileEnable == true {
		//日志切割
		hook := lumberjack.Logger{
			Filename:   l.Filename,
			MaxSize:    l.MaxFileSize,
			MaxBackups: l.MaxFileBackup,
		}

		if l.ConsoleEnable == true {
			// 终端和文件输出均开启
			zapWriteSyncer = zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stderr), zapcore.AddSync(&hook))
		} else {
			// 只开启文件输出
			zapWriteSyncer = zapcore.AddSync(&hook)
		}
	} else {
		// 只开启终端输出
		zapWriteSyncer = zapcore.AddSync(os.Stderr)
	}

	core := zapcore.NewCore(
		encoder,
		zapWriteSyncer,
		level,
	)
	logger := zap.New(core, zap.AddCaller())
	zap.ReplaceGlobals(logger)
}

func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(fmt.Sprintf("%04d-%02d-%02d-%02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second()))
}
