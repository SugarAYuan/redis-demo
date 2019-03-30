package tools

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"

	"github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

// public use
var Log = logrus.New()

// a hack to private use

func InitLog() *logrus.Logger {
	Log.AddHook(&DefaultFieldsHook{})

	filePath := Config.GetString("log.filepath")
	if len(filePath) > 0 {
		if ok, _ := PathExists(filePath); !ok {
			if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
				panic(fmt.Errorf("Create log path failed: %s\n", err.Error()))
			}
		}

		writer, err := rotatelogs.New(
			fmt.Sprintf("%s%s", filePath, "log%Y%m%d.log"),
			rotatelogs.WithMaxAge(Config.GetDuration("log.savetime")),
			rotatelogs.WithRotationTime(Config.GetDuration("log.rotationime")),
			// rotatelogs.WithLinkName(filePath),
		)
		if err == nil {
			Log.Out = writer
		}
	}

	switch Config.GetString("log.type") {
	case "json":
		Log.Formatter = &logrus.JSONFormatter{}
	case "text":
		Log.Formatter = &logrus.TextFormatter{
			DisableColors:   true,
			TimestampFormat: "2006-01-02 15:04:05",
			FullTimestamp:   true,
		}
	case "log":
		Log.Formatter = &DefaultFormatter{
			TimestampFormat: "2006/01/02 15:04:05",
			LogFormat:       "%time% [%lvl%] (%file%:%line%) %msg%",
		}
	default:
		Log.Formatter = &DefaultFormatter{
			TimestampFormat: "2006/01/02 15:04:05",
			LogFormat:       "%time% [%lvl%] (%file%:%line%) %msg%",
		}
	}

	switch Config.GetString("log.level") {
	case "debug":
		Log.SetLevel(logrus.DebugLevel)
	case "info":
		Log.SetLevel(logrus.InfoLevel)
	case "warn":
		Log.SetLevel(logrus.WarnLevel)
	case "error":
		Log.SetLevel(logrus.ErrorLevel)
	case "fatal":
		Log.SetLevel(logrus.FatalLevel)
	case "panic":
		Log.SetLevel(logrus.PanicLevel)
	default:
		Log.SetLevel(logrus.InfoLevel)
	}

	Log.Debug("log init success")
	return Log
}

type DefaultFieldsHook struct {
}

func (df *DefaultFieldsHook) Fire(entry *logrus.Entry) error {
	if _, file, line, ok := runtime.Caller(7); ok {
		entry.Data["file"] = path.Base(file)
		entry.Data["line"] = line
	}

	return nil
}

func (df *DefaultFieldsHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// Formatter implements logrus.Formatter interface
type DefaultFormatter struct {
	TimestampFormat string
	LogFormat       string
}

// Format building log message
func (f *DefaultFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	// time field
	b.WriteString(entry.Time.Format(f.TimestampFormat))

	// level field
	b.WriteString(" [" + strings.ToUpper(entry.Level.String()[:4]+"]"))

	// file no
	b.WriteString(" (" + entry.Data["file"].(string) + ":" + strconv.Itoa(entry.Data["line"].(int)) + ")")

	// msg field
	b.WriteString(" " + entry.Message)

	b.WriteByte('\n')
	return b.Bytes(), nil
}

// check if file or directory exists
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		// does exist
		return true, nil
	} else if os.IsNotExist(err) {
		// does not exist
		return false, nil
	} else {
		// stat error
		return false, err
	}
}
