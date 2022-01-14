package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

func Init() *zap.SugaredLogger {
	// 建立一個info會觸發的條件
	infoLevel := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level == zapcore.InfoLevel
	})
	// 建立一個只有error和fatal會觸發的條件
	errorFatalLevel := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level == zapcore.ErrorLevel || level == zapcore.FatalLevel
	})
	// write syncers
	stdoutSyncer := zapcore.Lock(os.Stdout)
	stderrSyncer := zapcore.Lock(os.Stderr)
	// info訊息會寫到access log
	accessWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "./log/access.log", // 要寫入的檔案
		MaxSize:    10,                 // log檔最大的大小，單位：MB(megabytes)
		MaxBackups: 3,                  // 最多幾個備份
		//MaxAge:     28, // log檔要存放幾天後自動刪除
	})
	// error訊息會寫到error log
	errorWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "./log/error.log",
		MaxSize:    10,
		MaxBackups: 3,
	})
	// 建立兩個core，分別去寫入info和error
	core := zapcore.NewTee(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
			// 這邊用Multi是因為要同時寫入檔案也要同時show在console
			zapcore.NewMultiWriteSyncer(stdoutSyncer, accessWriter),
			infoLevel,
		),
		zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
			zapcore.NewMultiWriteSyncer(stderrSyncer, errorWriter),
			errorFatalLevel,
		),
	)

	// 要記得加個caller和stacktrace
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
	return logger.Sugar()
}
