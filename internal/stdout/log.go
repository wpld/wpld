package stdout

import (
	"bytes"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

var log *logrus.Logger

func init() {
	log = logrus.New()
	log.SetLevel(logrus.WarnLevel)
}

func SetLogLevel(verbosity int) {
	if verbosity == 1 {
		log.SetLevel(logrus.InfoLevel)
	} else if verbosity == 2 {
		log.SetLevel(logrus.DebugLevel)
	} else if verbosity > 2 {
		log.SetLevel(logrus.TraceLevel)
	}
}

func Traceln(args ...interface{}) {
	PauseSpinner()
	log.Traceln(args...)
	StartSpinner()
}

func Debugln(args ...interface{}) {
	PauseSpinner()
	log.Debugln(args...)
	StartSpinner()
}

func Debugf(format string, args ...interface{}) {
	PauseSpinner()
	log.Debugf(format, args...)
	ResumeSpinner()
}

func Debugy(msg string, data interface{}) {
	buffer := bytes.NewBufferString("")

	encoder := yaml.NewEncoder(buffer)
	encoder.SetIndent(4)

	if err := encoder.Encode(map[string]interface{}{msg: data}); err == nil {
		PauseSpinner()
		log.Debugln(buffer.String())
		ResumeSpinner()
	}
}

func Infoln(args ...interface{}) {
	PauseSpinner()
	log.Infoln(args...)
	StartSpinner()
}

func Infof(format string, args ...interface{}) {
	PauseSpinner()
	log.Infof(format, args...)
	ResumeSpinner()
}

func Warnln(args ...interface{}) {
	PauseSpinner()
	log.Warnln(args...)
	StartSpinner()
}
