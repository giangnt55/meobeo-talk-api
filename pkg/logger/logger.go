package logger

import (
	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func Init(debug bool) {
	log = logrus.New()

	if debug {
		log.SetLevel(logrus.DebugLevel)
	} else {
		log.SetLevel(logrus.InfoLevel)
	}

	log.SetFormatter(&logrus.JSONFormatter{})
}

func Info(msg string, fields ...interface{}) {
	if len(fields) > 0 {
		log.WithFields(logrus.Fields{"data": fields}).Info(msg)
	} else {
		log.Info(msg)
	}
}

func Error(msg string, err error) {
	log.WithError(err).Error(msg)
}

func Fatal(msg string, err error) {
	log.WithError(err).Fatal(msg)
}

func Debug(msg string, fields ...interface{}) {
	if len(fields) > 0 {
		log.WithFields(logrus.Fields{"data": fields}).Debug(msg)
	} else {
		log.Debug(msg)
	}
}
