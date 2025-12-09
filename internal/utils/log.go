package utils

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

// LogLevel 日志级别
type LogLevel int

const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelWarn
	LevelError
)

var (
	logLevelNames = []string{"DEBUG", "INFO", "WARN", "ERROR"}
)

// Logger 日志记录器
type Logger struct {
	Level    LogLevel
	debugLog *log.Logger
	infoLog  *log.Logger
	warnLog  *log.Logger
	errorLog *log.Logger
	file     *os.File
}

// LogConfig 日志配置
type LogConfig struct {
	Level      LogLevel  // 日志级别
	Output     io.Writer // 输出位置
	FilePath   string    // 文件路径（如果输出到文件）
	MaxSize    int64     // 最大文件大小（字节）
	MaxBackups int       // 最大备份文件数
	MaxAge     int       // 最大保存天数
}

// DefaultLogConfig 默认日志配置
func DefaultLogConfig() LogConfig {
	return LogConfig{
		Level:    LevelInfo,
		Output:   os.Stdout,
		FilePath: "", // 默认输出到控制台
	}
}

// NewLogger 创建新的日志记录器
func NewLogger(cfg LogConfig) (*Logger, error) {
	var output io.Writer = cfg.Output
	var file *os.File

	// 如果指定了文件路径，则输出到文件
	if cfg.FilePath != "" {
		// 确保目录存在
		dir := filepath.Dir(cfg.FilePath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create log directory: %w", err)
		}

		f, err := os.OpenFile(cfg.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %w", err)
		}
		file = f
		output = f

		// 如果是控制台输出，同时输出到控制台和文件
		if cfg.Output == os.Stdout || cfg.Output == os.Stderr {
			output = io.MultiWriter(cfg.Output, f)
		}
	}

	return &Logger{
		Level:    cfg.Level,
		debugLog: log.New(output, "[DEBUG] ", log.Ldate|log.Ltime|log.Lshortfile),
		infoLog:  log.New(output, "[INFO]  ", log.Ldate|log.Ltime|log.Lshortfile),
		warnLog:  log.New(output, "[WARN]  ", log.Ldate|log.Ltime|log.Lshortfile),
		errorLog: log.New(output, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile),
		file:     file,
	}, nil
}

// ContextLogger 带上下文的日志记录器
type ContextLogger struct {
	*Logger
	ctx context.Context
}

// NewContextLogger 创建带上下文的日志记录器
func NewContextLogger(ctx context.Context, logger *Logger) *ContextLogger {
	return &ContextLogger{
		Logger: logger,
		ctx:    ctx,
	}
}

// WithContext 为日志记录器添加上下文
func (l *Logger) WithContext(ctx context.Context) *ContextLogger {
	return NewContextLogger(ctx, l)
}

// getCallerInfo 获取调用者信息
func getCallerInfo() string {
	_, file, line, ok := runtime.Caller(3) // 跳过日志函数的调用栈
	if !ok {
		return ""
	}
	return fmt.Sprintf("%s:%d", filepath.Base(file), line)
}

// Debug 记录调试日志
func (l *Logger) Debug(format string, v ...interface{}) {
	if l.Level <= LevelDebug {
		caller := getCallerInfo()
		if caller != "" {
			format = fmt.Sprintf("%s %s", caller, format)
		}
		l.debugLog.Printf(format, v...)
	}
}

func (cl *ContextLogger) Debug(format string, v ...interface{}) {
	if cl.Level <= LevelDebug {
		// 从上下文中获取请求ID等信息
		reqID := ""
		if id, ok := cl.ctx.Value("request_id").(string); ok {
			reqID = fmt.Sprintf("[req:%s] ", id)
		}
		caller := getCallerInfo()
		if caller != "" {
			format = fmt.Sprintf("%s%s %s", reqID, caller, format)
		} else {
			format = fmt.Sprintf("%s%s", reqID, format)
		}
		cl.debugLog.Printf(format, v...)
	}
}

// Info 记录信息日志
func (l *Logger) Info(format string, v ...interface{}) {
	if l.Level <= LevelInfo {
		caller := getCallerInfo()
		if caller != "" {
			format = fmt.Sprintf("%s %s", caller, format)
		}
		l.infoLog.Printf(format, v...)
	}
}

func (cl *ContextLogger) Info(format string, v ...interface{}) {
	if cl.Level <= LevelInfo {
		reqID := ""
		if id, ok := cl.ctx.Value("request_id").(string); ok {
			reqID = fmt.Sprintf("[req:%s] ", id)
		}
		caller := getCallerInfo()
		if caller != "" {
			format = fmt.Sprintf("%s%s %s", reqID, caller, format)
		} else {
			format = fmt.Sprintf("%s%s", reqID, format)
		}
		cl.infoLog.Printf(format, v...)
	}
}

// Warn 记录警告日志
func (l *Logger) Warn(format string, v ...interface{}) {
	if l.Level <= LevelWarn {
		caller := getCallerInfo()
		if caller != "" {
			format = fmt.Sprintf("%s %s", caller, format)
		}
		l.warnLog.Printf(format, v...)
	}
}

func (cl *ContextLogger) Warn(format string, v ...interface{}) {
	if cl.Level <= LevelWarn {
		reqID := ""
		if id, ok := cl.ctx.Value("request_id").(string); ok {
			reqID = fmt.Sprintf("[req:%s] ", id)
		}
		caller := getCallerInfo()
		if caller != "" {
			format = fmt.Sprintf("%s%s %s", reqID, caller, format)
		} else {
			format = fmt.Sprintf("%s%s", reqID, format)
		}
		cl.warnLog.Printf(format, v...)
	}
}

// Error 记录错误日志
func (l *Logger) Error(format string, v ...interface{}) {
	if l.Level <= LevelError {
		caller := getCallerInfo()
		if caller != "" {
			format = fmt.Sprintf("%s %s", caller, format)
		}
		l.errorLog.Printf(format, v...)
	}
}

func (cl *ContextLogger) Error(format string, v ...interface{}) {
	if cl.Level <= LevelError {
		reqID := ""
		if id, ok := cl.ctx.Value("request_id").(string); ok {
			reqID = fmt.Sprintf("[req:%s] ", id)
		}
		caller := getCallerInfo()
		if caller != "" {
			format = fmt.Sprintf("%s%s %s", reqID, caller, format)
		} else {
			format = fmt.Sprintf("%s%s", reqID, format)
		}
		cl.errorLog.Printf(format, v...)
	}
}

// Fatal 记录致命错误并退出
func (l *Logger) Fatal(format string, v ...interface{}) {
	caller := getCallerInfo()
	if caller != "" {
		format = fmt.Sprintf("%s %s", caller, format)
	}
	l.errorLog.Fatalf(format, v...)
}

// Close 关闭日志文件
func (l *Logger) Close() error {
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

// GlobalLogger 全局日志记录器
var (
	globalLogger *Logger
)

// InitGlobalLogger 初始化全局日志记录器
func InitGlobalLogger(cfg LogConfig) error {
	logger, err := NewLogger(cfg)
	if err != nil {
		return err
	}
	globalLogger = logger
	return nil
}

// GetLogger 获取全局日志记录器
func GetLogger() *Logger {
	if globalLogger == nil {
		// 默认使用控制台输出
		logger, _ := NewLogger(DefaultLogConfig())
		globalLogger = logger
	}
	return globalLogger
}

// LoggerFromContext 从上下文中获取日志记录器
func LoggerFromContext(ctx context.Context) *ContextLogger {
	if logger, ok := ctx.Value("logger").(*ContextLogger); ok {
		return logger
	}
	return GetLogger().WithContext(ctx)
}

// LogLevelFromString 将字符串转换为日志级别
func LogLevelFromString(level string) LogLevel {
	switch level {
	case "debug":
		return LevelDebug
	case "info":
		return LevelInfo
	case "warn":
		return LevelWarn
	case "error":
		return LevelError
	default:
		return LevelInfo
	}
}

// StringFromLogLevel 将日志级别转换为字符串
func StringFromLogLevel(level LogLevel) string {
	if level >= LevelDebug && level <= LevelError {
		return logLevelNames[level]
	}
	return "INFO"
}
