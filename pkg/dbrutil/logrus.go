package dbrutil

import (
	"github.com/sirupsen/logrus"
)

// LogrusEventReceiver is a sentinel EventReceiver; use it if the caller doesn't supply one
type LogrusEventReceiver struct {
	logger *logrus.Entry
}

func NewLogrusLogger(logger *logrus.Entry) *LogrusEventReceiver {
	logger = logger.WithField("system", "sql")
	return &LogrusEventReceiver{
		logger: logger,
	}
}

// Event receives a simple notification when various events occur
func (n *LogrusEventReceiver) Event(eventName string) {
	//n.logger.Info(eventName)
}

// EventKv receives a notification when various events occur along with
// optional key/value data
func (n *LogrusEventReceiver) EventKv(eventName string, kvs map[string]string) {
	//logger := n.logger
	//for k, v := range kvs {
	//	logger = logger.WithField(k, v)
	//}
	//logger.Info(eventName)
}

// EventErr receives a notification of an error if one occurs
func (n *LogrusEventReceiver) EventErr(eventName string, err error) error {
	n.logger.WithError(err).Error(eventName)
	return err
}

// EventErrKv receives a notification of an error if one occurs along with
// optional key/value data
func (n *LogrusEventReceiver) EventErrKv(eventName string, err error, kvs map[string]string) error {
	logger := n.logger
	for k, v := range kvs {
		logger = logger.WithField(k, v)
	}
	logger.WithError(err).Error(eventName)

	return err
}

// Timing receives the time an event took to happen
func (n *LogrusEventReceiver) Timing(eventName string, nanoseconds int64) {
	n.logger.WithField("time_ns", nanoseconds/1000000).Debug(eventName)
}

// TimingKv receives the time an event took to happen along with optional key/value data
func (n *LogrusEventReceiver) TimingKv(eventName string, nanoseconds int64, kvs map[string]string) {
	logger := n.logger
	for k, v := range kvs {
		logger = logger.WithField(k, v)
	}
	logger.WithField("time_ns", nanoseconds/1000000).Debug(eventName)
}
