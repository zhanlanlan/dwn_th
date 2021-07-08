package main

import (
	"testing"

	"github.com/golang/glog"
)

func TestFuck(t *testing.T) {
	glog.Error("fuck")
	glog.Warning("fuck")
	glog.Info("fuck")
}
