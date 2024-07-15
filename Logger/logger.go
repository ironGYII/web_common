package logger

import (
	"context"
	"github.com/ironGYI/commont/traceid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"time"
)

var logger *zap.Logger

type LoggerIndex string

func (l LoggerIndex) String() string {
	return string(l)
}

func InitLogger(level zapcore.Level, syncers []io.Writer, options []zap.Option, hooks []HookFunc) *zap.Logger {

	// 新增日志同步的输出位置
	var writeSyncers []zapcore.WriteSyncer // io writer
	for _, syncer := range syncers {
		writeSyncers = append(writeSyncers, zapcore.Lock(zapcore.AddSync(syncer)))
	}

	// 设置日志 level
	atomic_level := zap.NewAtomicLevel()
	atomic_level.SetLevel(level)

	// 设置日志编码格式
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.TimeEncoderOfLayout(time.DateTime)

	core := zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg), zapcore.NewMultiWriteSyncer(writeSyncers...), atomic_level)

	options = append(options,
		zap.AddStacktrace(zap.ErrorLevel), // 错误日志新增异常
		zap.AddCaller(),                   //日志打印输出 文件名, 行号, 函数名
	)

	// 设置hook
	core = RegisterHooks(core, hooks...)

	logger = zap.New(core, options...)
	return logger
}

func Info(ctx context.Context, index LoggerIndex, msg string, fields ...zap.Field) {
	if logger == nil {
		panic("logger not init")
	}
	fields = append(fields, getCtxfields(ctx, index)...)
	logger.Info(msg, fields...)
}

func Error(ctx context.Context, index LoggerIndex, msg string, fields ...zap.Field) {
	if logger == nil {
		panic("logger not init")
	}
	fields = append(fields, getCtxfields(ctx, index)...)
	logger.Error(msg, fields...)
}

func Panic(ctx context.Context, index LoggerIndex, msg string, fields ...zap.Field) {
	if logger == nil {
		panic("logger not init")
	}
	fields = append(fields, getCtxfields(ctx, index)...)
	logger.Panic(msg, fields...)
}

func DPanic(ctx context.Context, index LoggerIndex, msg string, fields ...zap.Field) {
	if logger == nil {
		panic("logger not init")
	}
	fields = append(fields, getCtxfields(ctx, index)...)
	logger.DPanic(msg, fields...)
}

func Warn(ctx context.Context, index LoggerIndex, msg string, fields ...zap.Field) {
	if logger == nil {
		panic("logger not init")
	}
	fields = append(fields, getCtxfields(ctx, index)...)
	logger.Warn(msg, fields...)
}

func Debug(ctx context.Context, index LoggerIndex, msg string, fields ...zap.Field) {
	if logger == nil {
		panic("logger not init")
	}
	fields = append(fields, getCtxfields(ctx, index)...)
	logger.Debug(msg, fields...)
}

func Fatal(ctx context.Context, index LoggerIndex, msg string, fields ...zap.Field) {
	if logger == nil {
		panic("logger not init")
	}
	fields = append(fields, getCtxfields(ctx, index)...)
	logger.Fatal(msg, fields...)
}

func getCtxfields(ctx context.Context, index LoggerIndex) []zap.Field {
	ctxFields := make([]zap.Field, 3)
	ctxFields[0] = zap.String("traceid", traceid.SpanFromContext(ctx).SpanContext().TraceID().String())
	ctxFields[1] = zap.String("spanid", traceid.SpanFromContext(ctx).SpanContext().SpanID().String())
	ctxFields[2] = zap.String("index", index.String())
	return ctxFields
}
