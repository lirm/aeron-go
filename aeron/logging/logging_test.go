package logging

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	logname = "logtest"
)

func TestLogLevels(t *testing.T) {
	assert := assert.New(t)
	logger := MustGetLogger(logname)

	// We initialzie to INFO
	assert.EqualValues(INFO, GetLevel("logtest"))
	assert.EqualValues(INFO, logger.GetLevel())

	// Set to Warning (via name)
	SetLevel(WARNING, logname)
	assert.EqualValues(WARNING, GetLevel("logtest"))
	assert.EqualValues(WARNING, logger.GetLevel())

	// Set to Info (via logger)
	logger.SetLevel(INFO)
	assert.EqualValues(INFO, GetLevel("logtest"))
	assert.EqualValues(INFO, logger.GetLevel())

	// logger.Fatal("Log at Fatal")
	logger.Error("Log at Error")
	logger.Warning("Log at Warning")
	logger.Info("Log at Info")
	logger.Debug("Log at Debug") // Silent

}
