package logger

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"template/pkg/converter"
	"template/pkg/dotenv"
	"time"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sirupsen/logrus"
)

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func(l *Logger) Logging(ctx context.Context) {
	var (
		logr = logrus.New()
	)

	logr.SetLevel(logrus.DebugLevel)
	logr.SetFormatter(new(formatter))
	logr.SetReportCaller(true)
	logr.SetOutput(os.Stdout)

	var errString string

	if l.Error != nil { errString = l.Error.Error() }

	logFormat := map[string]interface{}{
		"LOG LEVEL":         strings.ToUpper(l.LogLevel),
		"SERVICE_NAME":      strings.ToUpper(l.ServiceName),
		"MODULE_NAME":       strings.ToUpper(l.ModuleName),
		"COLLECTION_NAME":   strings.ToUpper(l.CollectionName),
		"URL":               l.UrlSource,
		"METHOD":            strings.ToUpper(l.Method),
		"ERROR":             errString,
		"DETAIL":            l.Detail,
	}

	if Log.newRelic != nil {
		jsonData, _ := converter.ConvertInterfaceToJSON(l)
		logEvent := newrelic.LogData{
			Timestamp: time.Now().Unix(),
			Severity:  l.CollectionName,
			Message:   string(jsonData),
		}
		Log.newRelic.RecordLog(logEvent)
	}

	if os.Getenv("IS_USE_FILE") == "true" {
		pathName := dotenv.GetString("LOG_PATH", "logs/")
		if ok, err := pathExists(pathName); !ok {
			if err != nil {
				fmt.Println(err)
			}

			err := os.Mkdir(pathName, os.ModePerm)
			if err != nil {
				fmt.Printf("Failed make directory logs: %v\n", err)
			}
		}

		var (
			file     *os.File
			filename = fmt.Sprintf("%s-%s.log", strings.ToUpper(l.ModuleName), DateNow)
			path     = fmt.Sprintf("%s%s", pathName, filename)
		)

		if _, err := os.Stat(path); err != nil {
			file, err = os.Create(path)
			if err != nil {
				fmt.Printf("Failed create file logs: %v\n", err)
			}
		} else {
			file, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
			if err != nil {
				fmt.Printf("Failed open file logs: %v\n", err)
			}
		}
		defer file.Close()

		logr.SetOutput(io.MultiWriter(file, os.Stdout))
		logr.SetReportCaller(true)
	}

	switch l.LogLevel {
	case "Error":
		logr.WithFields(logFormat).Error(l.CollectionName)
	case "Warning":
		logr.WithFields(logFormat).Warn(l.CollectionName)
	case "Info":
		logr.WithFields(logFormat).Info(l.CollectionName)
	default:
		logr.WithFields(logFormat).Trace(l.CollectionName)
	}
}