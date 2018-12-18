package log

import (
	"bytes"
	"io/ioutil"
	"path"
	"path/filepath"
	"time"

	"github.com/copernet/whccommon/model"
	"github.com/lestrrat/go-file-rotatelogs"
	"github.com/rifflock/lfshook"

	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

var DefaultTraceLabel = "traceID"

type ResultWriter struct {
	resp *bytes.Buffer
	gin.ResponseWriter
}

func (rw *ResultWriter) Write(p []byte) (int, error) {
	size, err := rw.resp.Write(p)
	if err != nil {
		return size, err
	}
	return rw.ResponseWriter.Write(p)
}

func LogContext(filter []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := c.Query(DefaultTraceLabel)
		if traceID == "" {
			uniqueID, _ := uuid.NewV4()
			traceID = uniqueID.String()
		}

		// unique track request identification, use c.MustGet("reqID") to
		// get this unique environment variable
		c.Set(DefaultTraceLabel, traceID)

		rw := &ResultWriter{
			resp:           bytes.NewBuffer(nil),
			ResponseWriter: c.Writer,
		}
		c.Writer = rw

		c.Next()

		if !inArray(c.Request.URL.Path, filter) {
			logrus.WithFields(logrus.Fields{
				"traceID": traceID,
				"IP":      c.ClientIP(),
				"URI":     c.Request.RequestURI,
				"Params":  c.Request.Form,
			}).Info("request information:")

			logrus.WithFields(logrus.Fields{
				"traceID":  traceID,
				"IP":       c.ClientIP(),
				"URI":      c.Request.RequestURI,
				"Params":   c.Request.Form,
				"Response": rw.resp.String(),
			}).Info("response information:")
		}
	}
}

func inArray(field string, array []string) bool {
	if array == nil {
		return false
	}

	for _, item := range array {
		if field == item {
			return true
		}
	}

	return false
}

func WithCtx(c context.Context) *logrus.Entry {
	if ginCtx, ok := c.(*gin.Context); ok {
		if ginCtx == nil {
			return logrus.WithField("traceID", "Empty TraceID")
		}

		uid, existed := ginCtx.Get(DefaultTraceLabel)
		if !existed {
			uniqueID, _ := uuid.NewV4()
			uid = uniqueID.String()

			ginCtx.Set(DefaultTraceLabel, uid)
		}
		return logrus.WithFields(logrus.Fields{
			"traceID": uid,
			"IP":      ginCtx.ClientIP(),
			"URI":     ginCtx.Request.RequestURI,
			"Params":  ginCtx.Request.Form,
		})
	}

	return logrus.WithField("traceID", c.Value(DefaultTraceLabel))
}

const (
	DefaultLogLevel = logrus.DebugLevel
	DefaultLogDir   = "/logs/"
)

func SetLogLevel(level string) error {
	l, err := logrus.ParseLevel(level)
	if err == nil {
		logrus.SetLevel(l)
	}

	return err
}

func InitLog(conf *model.LogOption) {
	// format log output
	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true, TimestampFormat: "2006-01-02 15:04:05", DisableTimestamp: false})

	// set log level
	err := SetLogLevel(conf.Level)
	if err != nil {
		logrus.SetLevel(DefaultLogLevel)
	}

	currentPath, err := filepath.Abs("./")
	if err != nil {
		panic("can not get current file abs path")
	}

	baseLogPath := path.Join(currentPath+DefaultLogDir, conf.Filename)
	writer, err := rotatelogs.New(
		baseLogPath+".%Y-%m-%d",
		rotatelogs.WithLinkName(baseLogPath),
		rotatelogs.WithMaxAge(conf.MaxAge*time.Hour), //default 7 days
		rotatelogs.WithRotationTime(24*time.Hour),
	)

	if err != nil {
		logrus.Errorf("config local file system logger error. %+v ", err)
	}

	lfHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: writer,
		logrus.InfoLevel:  writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: writer,
		logrus.FatalLevel: writer,
		logrus.PanicLevel: writer,
	}, &logrus.TextFormatter{FullTimestamp: true, TimestampFormat: "2006-01-02 15:04:05"})

	logrus.AddHook(lfHook)

	logrus.SetOutput(ioutil.Discard)
}
