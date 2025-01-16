package common

import (
	"testing"

	log "github.com/wifi32767/TikTokMall/common/logger"
)

func TestLogger(t *testing.T) {
	log.Init("debug")
	log.Debug("debug")
	log.Info("info")
	log.Warn("warn")
	log.Error("error")
}
