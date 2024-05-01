package logger

import (
	"os"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sirupsen/logrus"
)

func InitLog(nr *newrelic.Application) {
	logr := logrus.New()
	logr.SetFormatter(new(formatter))
	logr.SetReportCaller(true)
	logr.SetOutput(os.Stdout)
	Log.Logger = logr

	if nr != nil {
		Log.newRelic = nr
	}
}