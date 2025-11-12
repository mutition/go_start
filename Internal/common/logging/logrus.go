package logging

import (
	"strconv"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

func Init() {
	logrus.SetLevel(logrus.DebugLevel)
	SetFormatter()
}

func SetFormatter() {
	logrus.SetFormatter(&logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "time",
			logrus.FieldKeyLevel: "severity",
			logrus.FieldKeyMsg:   "message",
		},
	})
	if isLocal, _ := strconv.ParseBool(viper.GetString("local")); isLocal {
		logrus.SetFormatter(
			&prefixed.TextFormatter{
				ForceFormatting: true,
			},
		)
	}

}
