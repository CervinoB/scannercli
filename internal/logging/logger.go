package logging

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

func init() {
	Logger = logrus.New()

	// Set default configuration
	Logger.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
	})

	Logger.SetOutput(os.Stdout)
	Logger.SetLevel(logrus.InfoLevel) // Default level
}

func ConfigureLogger(verbose, debug bool) {
	switch {
	case debug:
		Logger.SetLevel(logrus.DebugLevel)
		Logger.SetReportCaller(true)
	case verbose:
		Logger.SetLevel(logrus.InfoLevel)
	default:
		Logger.SetLevel(logrus.WarnLevel)
	}
}
