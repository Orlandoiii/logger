package logger

import (
	"fmt"
	"io"

	"strings"
	"time"

	"os"
	"path/filepath"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

const dateFormat = "2006-01-02T15:04:05.999999"

type FilteredWriter struct {
	w     zerolog.LevelWriter
	level zerolog.Level
}

func (w *FilteredWriter) Write(p []byte) (n int, err error) {
	return w.w.Write(p)
}
func (w *FilteredWriter) WriteLevel(level zerolog.Level, p []byte) (n int, err error) {
	if level == w.level {
		return w.w.WriteLevel(level, p)
	}
	return len(p), nil
}

var log zerolog.Logger

var name string

func load(config *ConfigLogger, asConsole bool) error {

	serverName, err := os.Hostname()
	if err != nil {
		return err
	}

	path := config.Ruta

	var writersTrace io.Writer
	var writersDebug io.Writer
	var writersInfo io.Writer
	var writersWarn io.Writer
	var writersError io.Writer
	var writersFatal io.Writer

	var writers []io.Writer

	var traceFilter FilteredWriter
	var debugFilter FilteredWriter
	var infoFilter FilteredWriter
	var warnFilter FilteredWriter
	var errorFilter FilteredWriter
	var fatalFilter FilteredWriter

	zerolog.TimeFieldFormat = dateFormat

	if asConsole && config.Console {
		if config.BeutifyConsoleLog {
			output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
			output.FormatLevel = func(i interface{}) string {
				return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
			}
			writers = append(writers, output)
		} else {
			writers = append(writers, os.Stdout)
		}
	}
	switch config.MinLevel {
	case "trace":
		{
			tracePath := filepath.Join(path, "TRACE")
			writersTrace = &lumberjack.Logger{
				Filename:   tracePath + "/Trace.log", // File name
				MaxSize:    config.RotationMaxSizeMB, // Size in MB before file gets rotated
				MaxBackups: config.MaxBackups,        // Max number of files kept before being overwritten
				MaxAge:     config.MaxAgeDay,         // Max number of days to keep the files
				Compress:   config.Compress,
			}
			traceWriter := zerolog.MultiLevelWriter(writersTrace)
			traceFilter = FilteredWriter{traceWriter, zerolog.TraceLevel}
			writers = append(writers, &traceFilter)
		}
		fallthrough
	case "debug":
		{
			debugPath := filepath.Join(path, "DEBUG")
			writersDebug = &lumberjack.Logger{
				Filename:   debugPath + "/Debug.log", // File name
				MaxSize:    config.RotationMaxSizeMB, // Size in MB before file gets rotated
				MaxBackups: config.MaxBackups,        // Max number of files kept before being overwritten
				MaxAge:     config.MaxAgeDay,         // Max number of days to keep the files
				Compress:   config.Compress,
			}
			debugWriter := zerolog.MultiLevelWriter(writersDebug)
			debugFilter = FilteredWriter{debugWriter, zerolog.DebugLevel}
			writers = append(writers, &debugFilter)

		}
		fallthrough
	case "info":
		{
			infoPath := filepath.Join(path, "INFO")
			writersInfo = &lumberjack.Logger{
				Filename:   infoPath + "/Info.log",
				MaxSize:    config.RotationMaxSizeMB, // Size in MB before file gets rotated
				MaxBackups: config.MaxBackups,        // Max number of files kept before being overwritten
				MaxAge:     config.MaxAgeDay,         // Max number of days to keep the files
				Compress:   config.Compress,
			}
			infoWritter := zerolog.MultiLevelWriter(writersInfo)
			infoFilter = FilteredWriter{infoWritter, zerolog.InfoLevel}
			writers = append(writers, &infoFilter)

		}
		fallthrough
	case "warn":
		{
			warnPath := filepath.Join(path, "WARN")
			writersWarn = &lumberjack.Logger{
				Filename:   warnPath + "/Warn.log",   // File name
				MaxSize:    config.RotationMaxSizeMB, // Size in MB before file gets rotated
				MaxBackups: config.MaxBackups,        // Max number of files kept before being overwritten
				MaxAge:     config.MaxAgeDay,         // Max number of days to keep the files
				Compress:   config.Compress,
			}
			warnWritter := zerolog.MultiLevelWriter(writersWarn)
			warnFilter = FilteredWriter{warnWritter, zerolog.WarnLevel}
			writers = append(writers, &warnFilter)

		}
		fallthrough
	case "error":
		{
			errorPath := filepath.Join(path, "ERROR")

			writersError = &lumberjack.Logger{
				Filename:   errorPath + "/Error.log", // File name
				MaxSize:    config.RotationMaxSizeMB, // Size in MB before file gets rotated
				MaxBackups: config.MaxBackups,        // Max number of files kept before being overwritten
				MaxAge:     config.MaxAgeDay,         // Max number of days to keep the files
				Compress:   config.Compress,
			}
			errWritter := zerolog.MultiLevelWriter(writersError)
			errorFilter = FilteredWriter{errWritter, zerolog.ErrorLevel}
			writers = append(writers, &errorFilter)

		}
		fallthrough
	case "fatal":
		{
			fatalPath := filepath.Join(path, "FATAL")
			writersFatal = &lumberjack.Logger{
				Filename:   fatalPath + "/Fatal.log", // File name
				MaxSize:    config.RotationMaxSizeMB, // Size in MB before file gets rotated
				MaxBackups: config.MaxBackups,        // Max number of files kept before being overwritten
				MaxAge:     config.MaxAgeDay,         // Max number of days to keep the files
				Compress:   config.Compress,
			}
			fatalWriter := zerolog.MultiLevelWriter(writersFatal)
			fatalFilter = FilteredWriter{fatalWriter, zerolog.FatalLevel}
			writers = append(writers, &fatalFilter)

		}
	}
	w := zerolog.MultiLevelWriter(writers...)

	log = zerolog.New(w).With().Timestamp().Str(serverName, name).Logger()

	return nil
}

func Log() zerolog.Logger {
	return log
}

func Trace() *zerolog.Event {
	return log.Trace()

}

func Debug() *zerolog.Event {
	return log.Debug()

}

func Info() *zerolog.Event {
	return log.Info()

}

func Warn() *zerolog.Event {
	return log.Warn()

}

func Error() *zerolog.Event {
	return log.Error()

}
func Fatal() *zerolog.Event {
	return log.Fatal()
}

func GetLogggerWithIdentifiers(identifiers map[string]string) zerolog.Logger {

	child := log.With()
	for key, val := range identifiers {
		child = child.Str(key, val)
	}

	return child.Logger()
}

func Load(config *ConfigLogger, asConsole bool) error {
	err := load(config, asConsole)
	if err != nil {
		return err
	}
	return nil
}
