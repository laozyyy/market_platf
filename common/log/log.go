package log

import (
	"github.com/sirupsen/logrus"
)

var (
	Log *logrus.Logger
)

func init() {
	Log = logrus.New()
	Log.SetReportCaller(true)
	Log.SetFormatter(&logrus.TextFormatter{
		//以下设置只是为了使输出更美观
		DisableColors:   true,
		TimestampFormat: "2006-01-02 15:03:04",
	})
}

func Infof(format string, args ...interface{}) {
	if args == nil {
		Log.Info(format)
	} else {
		Log.Infof(format, args...)
	}

}

func Info(format string) {
	Log.Info(format)
}

func Errorf(format string, args ...interface{}) {
	if args == nil {
		Log.Error(format)
	} else {
		Log.Errorf(format, args...)
	}
}

func Error(format string) {
	Log.Error(format)
}
