package gortcp


import (
	"os"

	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("gortcp")

func InitLogger(debug bool) {
	format := logging.MustStringFormatter(
		`%{color}%{time:06-01-02 15:04:05.000} %{level:.4s} @%{shortfile}%{color:reset} %{message}`,
	)
	logging.SetFormatter(format)
	logging.SetBackend(logging.NewLogBackend(os.Stdout, "", 0))

	if debug {
		logging.SetLevel(logging.DEBUG, "gortcp")
	} else {
		logging.SetLevel(logging.INFO, "gortcp")
	}
}

func GetLogger() *logging.Logger {
	return logger
}
