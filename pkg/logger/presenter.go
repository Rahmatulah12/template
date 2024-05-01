package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sirupsen/logrus"
)

type logger struct {
	*logrus.Logger
	Env      string
	Name     string
	Date     string
	newRelic *newrelic.Application
}

var Log logger

type DetailLogger struct {
	IPClient	string	`json:"ipClient"`
	Header		string	`json:"header"`
	QueryParams	string	`json:"queryParams"`
	Request       string `json:"request"`
	Response      string `json:"response"`
}

type Logger struct {
	LogLevel        string  `json:"logLevel"`
	ServiceName     string  `json:"serviceName"`
	ModuleName      string  `json:"moduleName"`
	CollectionName  string  `json:"collectionName"`
	UrlSource       string  `json:"urlSource"`
	Method          string  `json:"method"`
	Error           error   `json:"error"`
	Detail			*DetailLogger `json:"detailLog"`
}

const (
	RED    = 31
	YELLOW = 33
	BLUE   = 36
	GRAY   = 37
)

var DateNow = time.Now().Format("2006-01-02")

type formatter struct{}

// Format custom formatter
func (f *formatter) Format(entry *logrus.Entry) ([]byte, error) {
	const formatColor = "\x1b[%dm"

	// output buffer
	b := &bytes.Buffer{}
	defer b.Reset()
	levelColor := getColorByLevel(entry.Level)

	_, _ = fmt.Fprintf(b, formatColor, levelColor)
	// Log Date Time
	now := time.Now().Format(time.RFC3339)
	b.WriteString("[")
	b.WriteString(now)
	b.WriteString("]")

	// Log level
	b.WriteString("[")
	level := strings.ToUpper(entry.Level.String())
	b.WriteString(level)
	b.WriteString("]")

	// Log direction
	// if entry.HasCaller() && f.env == "development" {
	// 	b.WriteString("[")
	// 	if f.isLocal {
	// 		_, _ = fmt.Fprintf(b, formatColor, levelColor)
	// 	}
	// 	fmt.Fprintf(
	// 		b,
	// 		"%s:%d",
	// 		entry.Caller.Function,
	// 		entry.Caller.Line,
	// 	)
	// 	if f.isLocal {
	// 		_, _ = fmt.Fprintf(b, formatColor, colorGray)
	// 	}
	// 	b.WriteString("]")
	// }

	// Log message
	if entry.Message != "" {
		b.WriteString("[")
		b.WriteString(entry.Message)
		b.WriteString("]")
	}

	keys := make([]string, 0, len(entry.Data))

	for key := range entry.Data {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	for _, key := range keys {
		json, _ := json.Marshal(entry.Data[key])
		b.WriteString(key)
		_, _ = fmt.Fprintf(b, formatColor, GRAY)
		b.WriteString(":")
		b.WriteString(string(json))
		_, _ = fmt.Fprintf(b, formatColor, levelColor)
	}

	b.WriteByte('\n')
	b.WriteByte('\n')
	return b.Bytes(), nil
}

func getColorByLevel(level logrus.Level) int {
	switch level {
	case logrus.DebugLevel:
		return GRAY
	case logrus.WarnLevel:
		return YELLOW
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		return RED
	default:
		return BLUE
	}
}
