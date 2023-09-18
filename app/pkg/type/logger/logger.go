package logger

import (
	"fmt"
	"os"
	"small/pkg/constans"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type options struct {
	IsProduction bool
	Level        zapcore.Level
}

type Logger struct {
	logger      *zap.Logger
	sugarLogger *zap.SugaredLogger
}

var opt *options

func init() {

	opt = &options{IsProduction: false}
	if strings.ToLower(strings.TrimSpace(os.Getenv("IS_PRODUCTION"))) == "true" {
		opt.IsProduction = true
	}

	switch strings.ToUpper(strings.TrimSpace(os.Getenv("LOG_LEVEL"))) {
	case "ERR", "ERROR":
		opt.Level = zapcore.ErrorLevel
	case "WARN", "WARNING":
		opt.Level = zapcore.WarnLevel
	case "INFO":
		opt.Level = zapcore.InfoLevel
	case "DEBUG":
		opt.Level = zapcore.DebugLevel
	case "FATAL":
		opt.Level = zapcore.FatalLevel
	default:
		opt.Level = zapcore.InfoLevel
	}
}

func new() (*Logger, error) {
	var config zap.Config

	if opt.IsProduction {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
	}

	config.Level = zap.NewAtomicLevelAt(opt.Level)

	newLogger, err := config.Build(zap.AddCallerSkip(2))
	if err != nil {
		return nil, err
	}

	newLogger.Info("Set LOG_LEVEL", zap.Stringer("level", opt.Level))

	log := &Logger{logger: newLogger, sugarLogger: newLogger.Sugar()}

	return log, nil
}

func New() (*Logger, error) {
	return new()
}

func (l *Logger) DebugF(msg string, args ...interface{}) {
	l.sugarLogger.Debugf(msg, args...)
}

func (l *Logger) KafkaProcessMessage(topic string, partition int, message string, workerID int, offset int64, time time.Time) {
	l.logger.Debug(
		"Processing Kafka message",
		zap.String(constans.Topic, topic),
		zap.Int(constans.Partition, partition),
		zap.String(constans.Message, message),
		zap.Int(constans.WorkerID, workerID),
		zap.Int64(constans.Offset, offset),
		zap.Time(constans.Time, time),
	)
}

func (l *Logger) Debug(args ...interface{}) {
	l.sugarLogger.Debug(args...)
}

func (l *Logger) InfoF(template string, args ...interface{}) {
	l.sugarLogger.Infof(template, args)
}
func (l *Logger) Errorf(template string, args ...interface{}) {
	l.sugarLogger.Errorf(template, args...)
}
func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.logger.Info(msg, fields...)
}

func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.logger.Warn(msg, fields...)
}

func (l *Logger) WarnF(msg string, args ...interface{}) {
	l.sugarLogger.Warnf(msg, args...)
}
func (l *Logger) Err(msg string, err error) {
	l.logger.Error(msg, zap.Error(err))
}
func (l *Logger) Error(msg interface{}, fields ...zap.Field) {
	if msg == nil {
		return
	}

	switch v := msg.(type) {
	case string:
		l.logger.Error(v, fields...)
	case error:
		l.logger.Error(v.Error(), fields...)
	case fmt.Stringer:
		l.logger.Error(v.String(), fields...)
	default:
		l.logger.Error(fmt.Sprintf("%v", v), fields...)
	}
}

func (l *Logger) Fatal(msg interface{}, fields ...zap.Field) {
	if msg == nil {
		return
	}

	switch msg.(type) {
	case string:
		if v, ok := msg.(string); ok {
			l.logger.Fatal(v, fields...)
		}
	case error:
		if v, ok := msg.(error); ok {
			l.logger.Fatal(v.Error(), fields...)
		}
	case fmt.Stringer:
		if v, ok := msg.(fmt.Stringer); ok {
			l.logger.Fatal(v.String(), fields...)
		}
	default:
		l.logger.Fatal(fmt.Sprintf("%v", msg), fields...)
	}
}

func (l *Logger) WarnMsg(msg string, err error) {
	l.logger.Warn(msg, zap.String("error", err.Error()))
}

func (l *Logger) KafkaLogCommittedMessage(topic string, partition int, offset int64) {
	l.logger.Info(
		"Committed Kafka message",
		zap.String(constans.Topic, topic),
		zap.Int(constans.Partition, partition),
		zap.Int64(constans.Offset, offset),
	)
}
