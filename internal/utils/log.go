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

// 声明当前常量。
const (
	// 更新当前变量或字段值。
	LevelDebug LogLevel = iota
	// 处理当前语句逻辑。
	LevelInfo
	// 处理当前语句逻辑。
	LevelWarn
	// 处理当前语句逻辑。
	LevelError
)

// Logger 日志记录器
type Logger struct {
	// 处理当前语句逻辑。
	Level LogLevel
	// 处理当前语句逻辑。
	debugLog *log.Logger
	// 处理当前语句逻辑。
	infoLog *log.Logger
	// 处理当前语句逻辑。
	warnLog *log.Logger
	// 处理当前语句逻辑。
	errorLog *log.Logger
	// 处理当前语句逻辑。
	file *os.File
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
	// 返回当前处理结果。
	return LogConfig{
		// 处理当前语句逻辑。
		Level: LevelInfo,
		// 处理当前语句逻辑。
		Output:   os.Stdout,
		FilePath: "", // 默认输出到控制台
	}
}

// NewLogger 创建新的日志记录器
func NewLogger(cfg LogConfig) (*Logger, error) {
	// 声明当前变量。
	var output io.Writer = cfg.Output
	// 声明当前变量。
	var file *os.File

	// 如果指定了文件路径，则输出到文件
	if cfg.FilePath != "" {
		// 确保目录存在
		dir := filepath.Dir(cfg.FilePath)
		// 判断条件并进入对应分支逻辑。
		if err := os.MkdirAll(dir, 0755); err != nil {
			// 返回当前处理结果。
			return nil, fmt.Errorf("failed to create log directory: %w", err)
		}

		// 定义并初始化当前变量。
		f, err := os.OpenFile(cfg.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		// 判断条件并进入对应分支逻辑。
		if err != nil {
			// 返回当前处理结果。
			return nil, fmt.Errorf("failed to open log file: %w", err)
		}
		// 更新当前变量或字段值。
		file = f
		// 更新当前变量或字段值。
		output = f

		// 如果是控制台输出，同时输出到控制台和文件
		if cfg.Output == os.Stdout || cfg.Output == os.Stderr {
			// 更新当前变量或字段值。
			output = io.MultiWriter(cfg.Output, f)
		}
	}

	// 返回当前处理结果。
	return &Logger{
		// 处理当前语句逻辑。
		Level: cfg.Level,
		// 调用log.New完成当前处理。
		debugLog: log.New(output, "[DEBUG] ", log.Ldate|log.Ltime|log.Lshortfile),
		// 调用log.New完成当前处理。
		infoLog: log.New(output, "[INFO]  ", log.Ldate|log.Ltime|log.Lshortfile),
		// 调用log.New完成当前处理。
		warnLog: log.New(output, "[WARN]  ", log.Ldate|log.Ltime|log.Lshortfile),
		// 调用log.New完成当前处理。
		errorLog: log.New(output, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile),
		// 处理当前语句逻辑。
		file: file,
		// 处理当前语句逻辑。
	}, nil
}

// ContextLogger 带上下文的日志记录器
type ContextLogger struct {
	*Logger
	// 处理当前语句逻辑。
	ctx context.Context
}

// NewContextLogger 创建带上下文的日志记录器
func NewContextLogger(ctx context.Context, logger *Logger) *ContextLogger {
	// 返回当前处理结果。
	return &ContextLogger{
		// 处理当前语句逻辑。
		Logger: logger,
		// 处理当前语句逻辑。
		ctx: ctx,
	}
}

// WithContext 为日志记录器添加上下文
func (l *Logger) WithContext(ctx context.Context) *ContextLogger {
	// 返回当前处理结果。
	return NewContextLogger(ctx, l)
}

// getCallerInfo 获取调用者信息
func getCallerInfo() string {
	_, file, line, ok := runtime.Caller(3) // 跳过日志函数的调用栈
	// 判断条件并进入对应分支逻辑。
	if !ok {
		// 返回当前处理结果。
		return ""
	}
	// 返回当前处理结果。
	return fmt.Sprintf("%s:%d", filepath.Base(file), line)
}

// Info 记录信息日志
func (l *Logger) Info(format string, v ...interface{}) {
	// 判断条件并进入对应分支逻辑。
	if l.Level <= LevelInfo {
		// 定义并初始化当前变量。
		caller := getCallerInfo()
		// 判断条件并进入对应分支逻辑。
		if caller != "" {
			// 更新当前变量或字段值。
			format = fmt.Sprintf("%s %s", caller, format)
		}
		// 调用l.infoLog.Printf完成当前处理。
		l.infoLog.Printf(format, v...)
	}
}

// Info 处理Info相关逻辑。
func (cl *ContextLogger) Info(format string, v ...interface{}) {
	// 判断条件并进入对应分支逻辑。
	if cl.Level <= LevelInfo {
		// 定义并初始化当前变量。
		reqID := ""
		// 判断条件并进入对应分支逻辑。
		if id, ok := cl.ctx.Value("request_id").(string); ok {
			// 更新当前变量或字段值。
			reqID = fmt.Sprintf("[req:%s] ", id)
		}
		// 定义并初始化当前变量。
		caller := getCallerInfo()
		// 判断条件并进入对应分支逻辑。
		if caller != "" {
			// 更新当前变量或字段值。
			format = fmt.Sprintf("%s%s %s", reqID, caller, format)
			// 进入新的代码块进行处理。
		} else {
			// 更新当前变量或字段值。
			format = fmt.Sprintf("%s%s", reqID, format)
		}
		// 调用cl.infoLog.Printf完成当前处理。
		cl.infoLog.Printf(format, v...)
	}
}

// Warn 记录警告日志
func (l *Logger) Warn(format string, v ...interface{}) {
	// 判断条件并进入对应分支逻辑。
	if l.Level <= LevelWarn {
		// 定义并初始化当前变量。
		caller := getCallerInfo()
		// 判断条件并进入对应分支逻辑。
		if caller != "" {
			// 更新当前变量或字段值。
			format = fmt.Sprintf("%s %s", caller, format)
		}
		// 调用l.warnLog.Printf完成当前处理。
		l.warnLog.Printf(format, v...)
	}
}

// Error 记录错误日志
func (l *Logger) Error(format string, v ...interface{}) {
	// 判断条件并进入对应分支逻辑。
	if l.Level <= LevelError {
		// 定义并初始化当前变量。
		caller := getCallerInfo()
		// 判断条件并进入对应分支逻辑。
		if caller != "" {
			// 更新当前变量或字段值。
			format = fmt.Sprintf("%s %s", caller, format)
		}
		// 调用l.errorLog.Printf完成当前处理。
		l.errorLog.Printf(format, v...)
	}
}

// Error 处理Error相关逻辑。
func (cl *ContextLogger) Error(format string, v ...interface{}) {
	// 判断条件并进入对应分支逻辑。
	if cl.Level <= LevelError {
		// 定义并初始化当前变量。
		reqID := ""
		// 判断条件并进入对应分支逻辑。
		if id, ok := cl.ctx.Value("request_id").(string); ok {
			// 更新当前变量或字段值。
			reqID = fmt.Sprintf("[req:%s] ", id)
		}
		// 定义并初始化当前变量。
		caller := getCallerInfo()
		// 判断条件并进入对应分支逻辑。
		if caller != "" {
			// 更新当前变量或字段值。
			format = fmt.Sprintf("%s%s %s", reqID, caller, format)
			// 进入新的代码块进行处理。
		} else {
			// 更新当前变量或字段值。
			format = fmt.Sprintf("%s%s", reqID, format)
		}
		// 调用cl.errorLog.Printf完成当前处理。
		cl.errorLog.Printf(format, v...)
	}
}

// Fatal 记录致命错误并退出
func (l *Logger) Fatal(format string, v ...interface{}) {
	// 定义并初始化当前变量。
	caller := getCallerInfo()
	// 判断条件并进入对应分支逻辑。
	if caller != "" {
		// 更新当前变量或字段值。
		format = fmt.Sprintf("%s %s", caller, format)
	}
	// 调用l.errorLog.Fatalf完成当前处理。
	l.errorLog.Fatalf(format, v...)
}

// Close 关闭日志文件
func (l *Logger) Close() error {
	// 判断条件并进入对应分支逻辑。
	if l.file != nil {
		// 返回当前处理结果。
		return l.file.Close()
	}
	// 返回当前处理结果。
	return nil
}

// GlobalLogger 全局日志记录器
var (
	// 处理当前语句逻辑。
	globalLogger *Logger
)

// InitGlobalLogger 初始化全局日志记录器
func InitGlobalLogger(cfg LogConfig) error {
	// 定义并初始化当前变量。
	logger, err := NewLogger(cfg)
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 返回当前处理结果。
		return err
	}
	// 更新当前变量或字段值。
	globalLogger = logger
	// 返回当前处理结果。
	return nil
}

// GetLogger 获取全局日志记录器
func GetLogger() *Logger {
	// 判断条件并进入对应分支逻辑。
	if globalLogger == nil {
		// 默认使用控制台输出
		logger, _ := NewLogger(DefaultLogConfig())
		// 更新当前变量或字段值。
		globalLogger = logger
	}
	// 返回当前处理结果。
	return globalLogger
}

// LoggerFromContext 从上下文中获取日志记录器
func LoggerFromContext(ctx context.Context) *ContextLogger {
	// 判断条件并进入对应分支逻辑。
	if logger, ok := ctx.Value("logger").(*ContextLogger); ok {
		// 返回当前处理结果。
		return logger
	}
	// 返回当前处理结果。
	return GetLogger().WithContext(ctx)
}

// LogLevelFromString 将字符串转换为日志级别
func LogLevelFromString(level string) LogLevel {
	// 根据表达式进入多分支处理。
	switch level {
	case "debug":
		// 返回当前处理结果。
		return LevelDebug
	case "info":
		// 返回当前处理结果。
		return LevelInfo
	case "warn":
		// 返回当前处理结果。
		return LevelWarn
	case "error":
		// 返回当前处理结果。
		return LevelError
	default:
		// 返回当前处理结果。
		return LevelInfo
	}
}
