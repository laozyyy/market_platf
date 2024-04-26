package common

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
