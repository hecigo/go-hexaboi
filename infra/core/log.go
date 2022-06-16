package core

import (
	log "github.com/sirupsen/logrus"
)

func InitLogger() {
	// Log as JSON instead of the default ASCII formatter.
	if GetBoolEnv("LOG_AS_JSON", false) {
		log.SetFormatter(&log.JSONFormatter{})
	} else {
		log.SetFormatter(&log.TextFormatter{
			ForceColors: true,
		})
	}
	switch Getenv("LOG_LEVEL", "Warn") {
	case "Trace":
		log.SetLevel(log.TraceLevel)
	case "Debug":
		log.SetLevel(log.DebugLevel)
	case "Info":
		log.SetLevel(log.InfoLevel)
	case "Warn":
		log.SetLevel(log.WarnLevel)
	case "Error":
		log.SetLevel(log.ErrorLevel)
	case "Fatal":
		log.SetLevel(log.FatalLevel)
	case "Panic":
		log.SetLevel(log.PanicLevel)
	}

	log.SetReportCaller(GetBoolEnv("LOG_METHOD_NAME", true))
}
