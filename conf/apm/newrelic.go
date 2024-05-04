package apm

import (
	"net/http"
	"os"
	"time"

	"github.com/newrelic/go-agent/v3/integrations/nrlogrus"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sirupsen/logrus"
)

type NewrelicConfig struct {
	AppName    string
	LicenseKey string
}

func InstanceNewrelic(nr *NewrelicConfig) *newrelic.Application {
	if nr == nil { panic("Failed to connect newrelic. Invalid credentials") }

	if nr.AppName == "" || nr.LicenseKey == "" { panic("Failed to connect newrelic. Invalid credentials") }

	httpTransport := &http.Transport{
		MaxIdleConns:          150,
		MaxIdleConnsPerHost:   50,
		MaxConnsPerHost:       150,
		IdleConnTimeout:       60,
		ResponseHeaderTimeout: 60 * time.Second,
	}


	newRelic, _ := newrelic.NewApplication(newrelic.ConfigAppName(
		nr.AppName),
		newrelic.ConfigLicense(nr.LicenseKey),
		newrelic.ConfigDistributedTracerEnabled(true),
		newrelic.ConfigAppLogForwardingEnabled(true),
		newrelic.ConfigCodeLevelMetricsEnabled(true),
		newrelic.ConfigAppLogDecoratingEnabled(false),
		newrelic.ConfigDebugLogger(os.Stdout),
		newrelic.ConfigAppLogEnabled(true),
		newrelic.ConfigModuleDependencyMetricsEnabled(true),
		newrelic.ConfigEnabled(true),
		func(configNR *newrelic.Config) {
			configNR.Enabled = true
			configNR.Transport = httpTransport
			configNR.HostDisplayName = nr.AppName
			configNR.TransactionEvents.Enabled = true
			log := logrus.New()
			log.SetLevel(logrus.DebugLevel)
			log.SetReportCaller(true)
			configNR.Logger = nrlogrus.StandardLogger()
		},
	)

	return newRelic
}