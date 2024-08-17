package initializers

import (
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
)

func SetupLogger() {
	env := os.Getenv("ENV")
	switch env {
	case "dev":
		{
			logrus.SetLevel(logrus.DebugLevel)
			b, err := strconv.ParseBool(os.Getenv("REPORT_CALLER"))
			if err != nil {
				logrus.SetReportCaller(false)
			} else {
				logrus.SetReportCaller(b)
			}
		}
	case "prod":
		{
			logrus.SetLevel(logrus.InfoLevel)
		}
	default:
		{
			logrus.SetLevel(logrus.InfoLevel)
		}
	}
}
