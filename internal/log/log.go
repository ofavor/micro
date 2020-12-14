package log

import (
	"fmt"
	"log"
	"strings"
)

const (
	levelFatal = iota
	levelError
	levelWarn
	levelInfo
	levelDebug
)

var levels = map[string]int{
	"fatal": levelFatal,
	"error": levelError,
	"warn":  levelWarn,
	"info":  levelInfo,
	"debug": levelDebug,
}

var currentLevel = levelError

func getPrefix(lv int) string {
	p := ""
	switch lv {
	case levelFatal:
		p = "FATAL"
	case levelError:
		p = "ERROR"
	case levelWarn:
		p = "WARN"
	case levelInfo:
		p = "INFO"
	case levelDebug:
		p = "DEBUG"
	}
	return fmt.Sprintf("[MICRO][% 5s]", p)
}

func output(lv int, args ...interface{}) {
	if currentLevel > lv {
		arr := append([]interface{}{getPrefix(lv)}, args...)
		if lv == levelFatal {
			log.Fatal(arr...)
		} else {
			log.Println(arr...)
		}
	}
}

func outputf(lv int, format string, args ...interface{}) {
	if currentLevel > lv {
		str := fmt.Sprintf(format, args...)
		if lv == levelFatal {
			log.Fatalln(getPrefix(lv), str)
		} else {
			log.Println(getPrefix(lv), str)
		}
	}
}

// SetLevel set log level
func SetLevel(lv string) {
	if l, ok := levels[strings.ToLower(lv)]; !ok {
		Warn("invalid log level:", lv)
	} else {
		currentLevel = l
	}
}

// Debug logs a message at level Debug on the standard logger.
func Debug(args ...interface{}) {
	output(levelDebug, args...)
}

// Info logs a message at level Info on the standard logger.
func Info(args ...interface{}) {
	output(levelInfo, args...)
}

// Warn logs a message at level Warn on the standard logger.
func Warn(args ...interface{}) {
	output(levelWarn, args...)
}

// Error logs a message at level Error on the standard logger.
func Error(args ...interface{}) {
	output(levelError, args...)
}

// Fatal logs a message at level Fatal on the standard logger then the process will exit with status set to 1.
func Fatal(args ...interface{}) {
	output(levelFatal, args...)
}

// Debugf logs a message at level Debug on the standard logger.
func Debugf(format string, args ...interface{}) {
	outputf(levelDebug, format, args...)
}

// Infof logs a message at level Info on the standard logger.
func Infof(format string, args ...interface{}) {
	outputf(levelInfo, format, args...)
}

// Warnf logs a message at level Warn on the standard logger.
func Warnf(format string, args ...interface{}) {
	outputf(levelWarn, format, args...)
}

// Errorf logs a message at level Error on the standard logger.
func Errorf(format string, args ...interface{}) {
	outputf(levelError, format, args...)
}

// Fatalf logs a message at level Fatal on the standard logger then the process will exit with status set to 1.
func Fatalf(format string, args ...interface{}) {
	outputf(levelFatal, format, args...)
}
