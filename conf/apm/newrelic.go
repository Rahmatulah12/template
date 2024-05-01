package apm

import (
	"net/http"
	"os"

	"github.com/newrelic/go-agent/v3/integrations/nrlogrus"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sirupsen/logrus"
)

type newrelicConfig struct {
	AppName    string
	LicenseKey string
}

func NewApmNewrelicConfig() *newrelicConfig {
	if os.Getenv("APP_NAME_NEWRELIC") == "" && os.Getenv("APP_KEY_NEWRELIC") == "" { panic("Failed to connect newrelic, invalid credentials.") }
	return &newrelicConfig{
		AppName:    os.Getenv("APP_NAME_NEWRELIC"),
		LicenseKey: os.Getenv("APP_KEY_NEWRELIC"),
	}
}

func InstanceNewrelic(nr *newrelicConfig) *newrelic.Application {
	if nr == nil { panic("Failed to connect newrelice. Onvalid credentials") }

	httpTransport := &http.Transport{
		MaxIdleConns:          50,
		MaxIdleConnsPerHost:   50,
		MaxConnsPerHost:       250,
		IdleConnTimeout:       10,
		ResponseHeaderTimeout: 10,
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