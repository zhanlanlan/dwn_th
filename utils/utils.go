package utils

import (
	"os"
	"strconv"

	"github.com/golang/glog"
)

func MustGetENV(key string) (value string) {
	key = os.Getenv(key)
	if key == "" {
		glog.Fatalf("key: %s is empty string", key)
	}

	return
}

func OptGetENV(key string) (value string) {
	return os.Getenv(key)
}

func OptGetEnvBool(key string) (b bool) {
	value := os.Getenv(key)

	b, _ = strconv.ParseBool(value)
	return
}
