package zap

import (
	"context"
	"github.com/micro-in-cn/tutorials/microservice-in-micro/part5/basic"
	"github.com/micro-in-cn/tutorials/microservice-in-micro/part5/basic/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path/filepath"
	"sync"
)

var (
	l                              *Logger
	sp                             = string(filepath.Separator)
	errWS, warnWS, infoWS, debugWS zapcore.WriteSyncer       // IO输出
	debugConsoleWS                 = zapcore.Lock(os.Stdout) // 控制台标准输出
	errorConsoleWS                 = zapcore.Lock(os.Stderr)
)

func init() {
	l = &Logger{
		Opts: &Options{},
	}
	basic.Register(initLogger)
}

type Logger struct {
	zl *zap.Logger
	sync.RWMutex
	Opts      *Options `json:"opts"`
	zapConfig zap.Config
	inited    bool
}

// Debug logs a message at DebugLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func (log *Logger) Debug(ctx context.Context, msg string, fields ...zap.Field) {
	log.zl.Debug(msg, zap.Field{
		Key:    "trace-id",
		Type:   zapcore.StringType,
		String: ctx.Value("trace-id").(string)})
}

// Info logs a message at InfoLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func (log *Logger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	log.zl.Info(msg, zap.Field{
		Key:    "trace-id",
		Type:   zapcore.StringType,
		String: ctx.Value("trace-id").(string)})
}

// Warn() logs a message at Warn()Level. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func (log *Logger) Warn(ctx context.Context, msg string, fields ...zap.Field) {
	log.zl.Warn(msg, zap.Field{
		Key:    "trace-id",
		Type:   zapcore.StringType,
		String: ctx.Value("trace-id").(string)})
}

// Error() logs a message at Error()Level. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func (log *Logger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	log.zl.Error(msg, zap.Field{
		Key:    "trace-id",
		Type:   zapcore.StringType,
		String: ctx.Value("trace-id").(string)})
}

func initLogger() {

	l.Lock()
	defer l.Unlock()

	if l.inited {
		l.Info(context.Background(), "[initLogger] logger Inited")
		return
	}

	l.loadCfg()
	l.init()
	l.zl.Info("[initLogger] zap plugin initializing completed")
	l.inited = true
}

// GetLogger returns logger
func GetLogger() (ret *Logger) {
	return l
}

func (l *Logger) init() {

	l.setSyncers()
	var err error

	l.zl, err = l.zapConfig.Build(l.cores())
	if err != nil {
		panic(err)
	}

	defer l.zl.Sync()
}

func (l *Logger) loadCfg() {

	c := config.C()

	err := c.Path("zap", l.Opts, "zap")
	if err != nil {
		panic(err)
	}

	if l.Opts.Development {
		l.zapConfig = zap.NewDevelopmentConfig()
	} else {
		l.zapConfig = zap.NewProductionConfig()
	}

	// application log output path
	if l.Opts.OutputPaths == nil || len(l.Opts.OutputPaths) == 0 {
		l.zapConfig.OutputPaths = []string{"stdout"}
	}

	//  error of zap-self log
	if l.Opts.ErrorOutputPaths == nil || len(l.Opts.ErrorOutputPaths) == 0 {
		l.zapConfig.OutputPaths = []string{"stderr"}
	}

	// 默认输出到程序运行目录的logs子目录
	if l.Opts.LogFileDir == "" {
		l.Opts.LogFileDir, _ = filepath.Abs(filepath.Dir(filepath.Join(".")))
		l.Opts.LogFileDir += sp + "logs" + sp
	}

	if l.Opts.AppName == "" {
		l.Opts.AppName = "app"
	}

	if l.Opts.ErrorFileName == "" {
		l.Opts.ErrorFileName = "error.log"
	}

	if l.Opts.WarnFileName == "" {
		l.Opts.WarnFileName = "warn.log"
	}

	if l.Opts.InfoFileName == "" {
		l.Opts.InfoFileName = "info.log"
	}

	if l.Opts.DebugFileName == "" {
		l.Opts.DebugFileName = "debug.log"
	}

	if l.Opts.MaxSize == 0 {
		l.Opts.MaxSize = 50
	}
	if l.Opts.MaxBackups == 0 {
		l.Opts.MaxBackups = 3
	}
	if l.Opts.MaxAge == 0 {
		l.Opts.MaxAge = 30
	}
}

func (l *Logger) setSyncers() {

	f := func(fN string) zapcore.WriteSyncer {
		return zapcore.AddSync(&lumberjack.Logger{
			Filename:   l.Opts.LogFileDir + sp + l.Opts.AppName + "-" + fN,
			MaxSize:    l.Opts.MaxSize,
			MaxBackups: l.Opts.MaxBackups,
			MaxAge:     l.Opts.MaxAge,
			Compress:   true,
			LocalTime:  true,
		})
	}

	errWS = f(l.Opts.ErrorFileName)
	warnWS = f(l.Opts.WarnFileName)
	infoWS = f(l.Opts.InfoFileName)
	debugWS = f(l.Opts.DebugFileName)

	return
}

func (l *Logger) cores() zap.Option {

	fileEncoder := zapcore.NewJSONEncoder(l.zapConfig.EncoderConfig)
	consoleEncoder := zapcore.NewConsoleEncoder(l.zapConfig.EncoderConfig)

	errPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl > zapcore.WarnLevel && zapcore.WarnLevel-l.zapConfig.Level.Level() > -1
	})
	warnPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.WarnLevel && zapcore.WarnLevel-l.zapConfig.Level.Level() > -1
	})
	infoPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.InfoLevel && zapcore.InfoLevel-l.zapConfig.Level.Level() > -1
	})
	debugPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.DebugLevel && zapcore.DebugLevel-l.zapConfig.Level.Level() > -1
	})

	cores := []zapcore.Core{
		// region 日志文件

		// error 及以上
		zapcore.NewCore(fileEncoder, errWS, errPriority),

		// warn
		zapcore.NewCore(fileEncoder, warnWS, warnPriority),

		// info
		zapcore.NewCore(fileEncoder, infoWS, infoPriority),

		// debug
		zapcore.NewCore(fileEncoder, debugWS, debugPriority),

		// endregion

		// region 控制台

		// 错误及以上
		zapcore.NewCore(consoleEncoder, errorConsoleWS, errPriority),

		// 警告
		zapcore.NewCore(consoleEncoder, debugConsoleWS, warnPriority),

		// info
		zapcore.NewCore(consoleEncoder, debugConsoleWS, infoPriority),

		// debug
		zapcore.NewCore(consoleEncoder, debugConsoleWS, debugPriority),

		// endregion
	}

	return zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		return zapcore.NewTee(cores...)
	})
}
