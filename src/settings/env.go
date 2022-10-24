package settings

import (
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
	"strings"
)

var ListenBind = "0.0.0.0"
var HttpPort = 8000 // HTTP_PORT

var EnableCacheHeaders = false // ENABLE_CACHE_HEADERS

func GetEnvOpts() {
	{
		envHttpPort, ok := os.LookupEnv("HTTP_PORT")
		if ok {
			i, err := strconv.Atoi(envHttpPort)
			if err != nil {
				logrus.Fatalf("Failed converting HTTP_PORT value '%s' to number: %s", envHttpPort, err)
			}
			HttpPort = i
		}
	}

	{
		_, ok := os.LookupEnv("DEBUG")
		if ok {
			logrus.SetLevel(logrus.DebugLevel)
		}
	}

	{
		value, ok := os.LookupEnv("ENABLE_CACHE_HEADERS")
		if ok && len(value) != 0 && (strings.ToLower(value) == "true" || value == "1") {
			logrus.Info("Cache Headers enabled")
			EnableCacheHeaders = true
		} else {
			logrus.Info("Cache Headers disabled")
		}
	}
}
