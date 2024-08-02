package logging

import (
	"log"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

func SettingUpLoggerWithconfigurations(logger *logrus.Logger, loggingConfig interface{}) {
	cfg, ok := loggingConfig.(map[string]interface{})
	if !ok {
		log.Fatal("Some error occurred while loading loggers configuration")
	}
	if cfg["level"] != nil {
		logger.SetLevel(transletLevel(cfg["level"].(string)))
	}

	if cfg["logFileConfigs"] != nil {
		logFileConfigs, ok := cfg["logFileConfigs"].(map[string]interface{})
		if ok {

			name := logFileConfigs["filename"].(string)
			path := logFileConfigs["path"].(string)
			maxAge := logFileConfigs["maxAge"].(int)
			maxBackup := logFileConfigs["maxBackup"].(int)
			isCompressed := logFileConfigs["compression"].(bool)
			maxSize := logFileConfigs["maxSize"].(int)
			logPath := path + "/" + name
			logger.SetOutput(&lumberjack.Logger{
				Filename:   logPath,
				MaxSize:    maxSize,      // Max size in MB
				MaxBackups: maxBackup,    // Max number of old log files to keep
				MaxAge:     maxAge,       // Max age in days to keep a log file
				Compress:   isCompressed, // Compress old log files
			})
		}
	}
}

func transletLevel(level string) (internalLevel logrus.Level) {
	switch level {
	case "PANIC":
		internalLevel = logrus.PanicLevel
	case "FATEL":
		internalLevel = logrus.FatalLevel
	case "ERROR":
		internalLevel = logrus.ErrorLevel
	case "WARN":
		internalLevel = logrus.WarnLevel
	case "INFO":
		internalLevel = logrus.InfoLevel
	case "DEBUG":
		internalLevel = logrus.DebugLevel
	case "TRACE":
		internalLevel = logrus.TraceLevel
	default:
		internalLevel = logrus.InfoLevel
	}
	return

}
