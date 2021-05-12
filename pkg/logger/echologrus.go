package logger

import (
	"io"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/sirupsen/logrus"
)

type Logrus struct {
	*logrus.Entry
}

var EchoLogger *logrus.Entry

func GetEchoLogger() Logrus {
	return Logrus{EchoLogger}
}

func (l Logrus) Level() log.Lvl {
	switch l.Logger.Level {
	case logrus.DebugLevel:
		return log.DEBUG
	case logrus.WarnLevel:
		return log.WARN
	case logrus.ErrorLevel:
		return log.ERROR
	case logrus.InfoLevel:
		return log.INFO
	default:
		l.Panic("Invalid level")
	}

	return log.OFF
}

func (l Logrus) SetHeader(_ string) {}

func (l Logrus) SetPrefix(s string) {}

func (l Logrus) Prefix() string {
	return ""
}

func (l Logrus) SetLevel(lvl log.Lvl) {
	switch lvl {
	case log.DEBUG:
		EchoLogger.Logger.SetLevel(logrus.DebugLevel)
	case log.WARN:
		EchoLogger.Logger.SetLevel(logrus.WarnLevel)
	case log.ERROR:
		EchoLogger.Logger.SetLevel(logrus.ErrorLevel)
	case log.INFO:
		EchoLogger.Logger.SetLevel(logrus.InfoLevel)
	default:
		l.Panic("invalid level")
	}
}

func (l Logrus) Output() io.Writer {
	return l.Logger.Out
}

func (l Logrus) SetOutput(w io.Writer) {
	EchoLogger.Logger.SetOutput(w)
}

func (l Logrus) Printj(j log.JSON) {
	EchoLogger.WithFields(logrus.Fields(j)).Print()
}

func (l Logrus) Debugj(j log.JSON) {
	EchoLogger.WithFields(logrus.Fields(j)).Debug()
}

func (l Logrus) Infoj(j log.JSON) {
	EchoLogger.WithFields(logrus.Fields(j)).Info()
}

func (l Logrus) Warnj(j log.JSON) {
	EchoLogger.WithFields(logrus.Fields(j)).Warn()
}

func (l Logrus) Errorj(j log.JSON) {
	EchoLogger.WithFields(logrus.Fields(j)).Error()
}

func (l Logrus) Fatalj(j log.JSON) {
	EchoLogger.WithFields(logrus.Fields(j)).Fatal()
}

func (l Logrus) Panicj(j log.JSON) {
	EchoLogger.WithFields(logrus.Fields(j)).Panic()
}

func (l Logrus) Print(i ...interface{}) {
	EchoLogger.Print(i[0].(string))
}

func (l Logrus) Debug(i ...interface{}) {
	EchoLogger.Debug(i[0].(string))
}

func (l Logrus) Info(i ...interface{}) {
	EchoLogger.Info(i[0].(string))
}

func (l Logrus) Warn(i ...interface{}) {
	EchoLogger.Warn(i[0].(string))
}

func (l Logrus) Error(i ...interface{}) {
	EchoLogger.Error(i[0].(string))
}

func (l Logrus) Fatal(i ...interface{}) {
	EchoLogger.Fatal(i[0].(string))
}

func (l Logrus) Panic(i ...interface{}) {
	EchoLogger.Panic(i[0].(string))
}

func logrusMiddlewareHandler(c echo.Context, next echo.HandlerFunc) error {
	req := c.Request()
	res := c.Response()
	start := time.Now()
	err := next(c)
	if err != nil {
		c.Error(err)
	}
	stop := time.Now()

	p := req.URL.Path

	bytesIn := req.Header.Get(echo.HeaderContentLength)

	l := EchoLogger.WithFields(map[string]interface{}{
		"time_rfc3339":  time.Now().Format(time.RFC3339),
		"remote_ip":     c.RealIP(),
		"host":          req.Host,
		"uri":           req.RequestURI,
		"method":        req.Method,
		"path":          p,
		"referer":       req.Referer(),
		"user_agent":    req.UserAgent(),
		"status":        res.Status,
		"latency":       strconv.FormatInt(stop.Sub(start).Nanoseconds()/1000, 10),
		"latency_human": stop.Sub(start).String(),
		"bytes_in":      bytesIn,
		"bytes_out":     strconv.FormatInt(res.Size, 10),
	})

	if err != nil {
		l.WithError(err).Error("handled request")
	} else {
		l.Info("handled request")
	}

	return nil
}

func echoLogger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		return logrusMiddlewareHandler(c, next)
	}
}

// Hook is a function to process middleware.
func NewEchoLogrus() echo.MiddlewareFunc {
	return echoLogger
}

