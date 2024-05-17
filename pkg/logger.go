package pkg

import (
	"context"
	"os"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server/binding"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"ququiz.org/lintang/quiz-query-service/config"

	hertzzap "github.com/hertz-contrib/logger/zap"
)

var lg *zap.Logger

// pake hertzlogger gak kayak pake uber/zap logger beneran
func InitZapLogger(cfg *config.Config) *hertzzap.Logger {
	productionCfg := zap.NewProductionEncoderConfig()
	productionCfg.TimeKey = "timestamp"
	productionCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	productionCfg.EncodeDuration = zapcore.SecondsDurationEncoder
	productionCfg.EncodeCaller = zapcore.ShortCallerEncoder

	developmentCfg := zap.NewDevelopmentEncoderConfig()
	developmentCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	// log encooder (json for prod, console for dev)
	consoleEncoder := zapcore.NewConsoleEncoder(developmentCfg)
	fileEncoder := zapcore.NewJSONEncoder(productionCfg)
	// loglevel
	logDevLevel := zap.NewAtomicLevelAt(zap.DebugLevel)
	logLevelProd := zap.NewAtomicLevelAt(zap.InfoLevel)

	//write sycer
	writeSyncerStdout, writeSyncerFile := GetLogWriter(cfg.MaxBackups, cfg.MaxAge)

	prodCfg := hertzzap.CoreConfig{
		Enc: fileEncoder,
		Ws:  writeSyncerFile,
		Lvl: logLevelProd,
	}

	devCfg := hertzzap.CoreConfig{
		Enc: consoleEncoder,
		Ws:  writeSyncerStdout,
		Lvl: logDevLevel,
	}
	logsCores := []hertzzap.CoreConfig{
		prodCfg,
		devCfg,
	}
	coreConsole := zapcore.NewCore(consoleEncoder, writeSyncerStdout, logDevLevel)
	coreFile := zapcore.NewCore(fileEncoder, writeSyncerFile, logLevelProd)
	core := zapcore.NewTee(
		coreConsole,
		coreFile,
	)
	lg = zap.New(core)
	zap.ReplaceGlobals(lg)

	prodAndDevLogger := hertzzap.NewLogger(hertzzap.WithZapOptions(zap.WithFatalHook(zapcore.WriteThenPanic)),
		hertzzap.WithCores(logsCores...))

	return prodAndDevLogger
}

func GetLogWriter(maxBackup, maxAge int) (writeSyncerStdout zapcore.WriteSyncer, writeSyncerFile zapcore.WriteSyncer) {
	file := zapcore.AddSync(&lumberjack.Logger{
		Filename: "./logs/app.log",

		MaxBackups: maxBackup,
		MaxAge:     maxAge,
	})
	stdout := zapcore.AddSync(os.Stdout)

	return stdout, file
}

type ValidateError struct {
	ErrType, FailField, Msg string
}

// Error implements error interface.
func (e *ValidateError) Error() string {
	if e.Msg != "" {
		return e.Msg
	}
	return e.ErrType + ": expr_path=" + e.FailField + ", cause=invalid"
}

type BindError struct {
	ErrType, FailField, Msg string
}

// Error implements error interface.
func (e *BindError) Error() string {
	if e.Msg != "" {
		return e.Msg
	}
	return e.ErrType + ": expr_path=" + e.FailField + ", cause=invalid"
}

func CreateCustomValidationError() *binding.ValidateConfig {
	validateConfig := &binding.ValidateConfig{}
	validateConfig.SetValidatorErrorFactory(func(failField, msg string) error {
		err := ValidateError{
			ErrType:   "validateErr",
			FailField: "[validateFailField]: " + failField,
			Msg:       msg,
		}

		return &err
	})
	return validateConfig
}

// accessLogger nbawaan zap bagus ini pas di load testing
func AccessLog() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		start := time.Now()
		path := string(ctx.Request.URI().Path()[:])
		query := string(ctx.Request.URI().QueryString()[:])
		ctx.Next(c)
		cost := time.Since(start)
		lg.Info(path,
			zap.Int("status", ctx.Response.StatusCode()),
			zap.String("method", string(ctx.Request.Header.Method())),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", ctx.ClientIP()),
			zap.String("user-agent", string(ctx.Host())),
			zap.String("errors", ctx.Errors.String()),
			zap.Duration("cost", cost),
		)
	}
}
