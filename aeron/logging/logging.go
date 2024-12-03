// Copyright (C) 2021 Talos, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Provides a transition layer from "github.com/op/go-logging" to
// "go.uber.org/zap" to simply resolve some reentrancy issues in go-logging.
//
// This provides a largely api compatible layer so we can quickly
// drop in a replacement.
package logging

import (
	"fmt"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Zaplogger is a container to wrap zap logging with the parts of the go-logging API we use
type ZapLogger struct {
	name   string
	sugar  *zap.SugaredLogger
	logger *zap.Logger
	config zap.Config
}

// Mapping of go-logging to zap log levels. It's imperfect but good enough
const (
	DEBUG    = zapcore.DebugLevel
	INFO     = zapcore.InfoLevel
	WARNING  = zapcore.WarnLevel
	NOTICE   = zapcore.InfoLevel // Map to info
	ERROR    = zapcore.ErrorLevel
	CRITICAL = zapcore.DPanicLevel
)

// The go-logging package creates named loggers and allows them to be
// accessed by name This is used in aeron to allow parent code to set
// the logging level of library components so we need to provide a
// mechanism for that
var namedLoggers sync.Map // [string]*ZapLogger

// MustGetLogger returns a new logger or panic()s
func MustGetLogger(name string) *ZapLogger {
	z := new(ZapLogger)
	z.name = name

	z.config = aeronLoggingConfig()
	logger, err := z.config.Build()
	if err != nil {
		panic("Failed to make logger")
	}
	z.logger = logger.Named(name).WithOptions(zap.AddCallerSkip(1))
	z.sugar = z.logger.Sugar()

	// Keep a reference
	namedLoggers.Store(name, z)

	return z
}

func GetLogger(name string) (*ZapLogger, error) {
	if z, ok := namedLoggers.Load(name); ok {
		return z.(*ZapLogger), nil
	} else {
		return nil, fmt.Errorf("unknown logger %s", name)
	}
}

// newAeronEncoderConfig returns a default opinionated EncoderConfig for aeron
func newAeronEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

// Default config
func aeronLoggingConfig() zap.Config {
	return zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "console",
		EncoderConfig:    newAeronEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}
}

// SetConfigAndRebuild so you can replace the default config
func (z *ZapLogger) SetConfigAndRebuild(c zap.Config) error {
	z.config = c
	logger, err := z.config.Build()
	if err != nil {
		return err
	}

	z.logger = logger.Named(z.name)
	z.sugar = z.logger.Sugar()
	return nil
}

// SetLevel on a named logger
func SetLevel(l zapcore.Level, name string) {
	z, ok := namedLoggers.Load(name)
	if ok {
		zlogger := z.(*ZapLogger)
		zlogger.SetLevel(l)
	}
}

// GetLevel on a named logger
func GetLevel(name string) zapcore.Level {
	z, ok := namedLoggers.Load(name)
	if ok {
		zlogger := z.(*ZapLogger)
		return zlogger.GetLevel()
	}

	// Bogus return but that's the API we are emulating and shouldn't happen
	return zapcore.InfoLevel
}

// Sugar returns the internal Sugared logger
func (z *ZapLogger) Sugar() *zap.SugaredLogger {
	return z.sugar
}

// SetSugar sets our internal Sugared logger
func (z *ZapLogger) SetSugar(s *zap.SugaredLogger) {
	z.sugar = s
}

// Logger returns the internal zap logger
func (z *ZapLogger) Logger() *zap.Logger {
	return z.logger
}

// SetLogger sets the internal zap logger
func (z *ZapLogger) SetLogger(l *zap.Logger) {
	z.logger = l
}

// SetLevel sets the log level at which we will log
func (z *ZapLogger) SetLevel(l zapcore.Level) {
	z.config.Level.SetLevel(l)
}

// GetLevel returns the log level at which we will log
func (z *ZapLogger) GetLevel() zapcore.Level {
	return z.config.Level.Level()
}

// IsEnabledFor returns true if logging is enabled for the specified level
func (z *ZapLogger) IsEnabledFor(level zapcore.Level) bool {
	return level >= z.GetLevel()
}

// Fatalf logs a formatted string at log level Fatal and will then always exit()
func (z *ZapLogger) Fatalf(template string, args ...interface{}) {
	z.sugar.Fatalf(template, args...)
}

// Errorf logs a formatted string at log level Error
func (z *ZapLogger) Errorf(template string, args ...interface{}) {
	z.sugar.Errorf(template, args...)
}

// Warningf logs a formatted string at log level Warning
func (z *ZapLogger) Warningf(template string, args ...interface{}) {
	z.sugar.Warnf(template, args...)
}

// Infof logs a formatted string at log level Info
func (z *ZapLogger) Infof(template string, args ...interface{}) {
	z.sugar.Infof(template, args...)
}

// Noticef logs a formatted string at log level *Info*
func (z *ZapLogger) Noticef(template string, args ...interface{}) {
	z.sugar.Infof(template, args...)
}

// Debugf logs a formatted string at log level Debug
func (z *ZapLogger) Debugf(template string, args ...interface{}) {
	z.sugar.Debugf(template, args...)
}

// Fatal logs it's arguments at log level Fatal
func (z *ZapLogger) Fatal(args ...interface{}) {
	z.sugar.Fatal(args...)
}

// Error logs it's arguments at log level Error
func (z *ZapLogger) Error(args ...interface{}) {
	z.sugar.Error(args...)
}

// Warning logs it's arguments at log level Warning
func (z *ZapLogger) Warning(args ...interface{}) {
	z.sugar.Warn(args...)
}

// Info logs it's arguments at log level Info
func (z *ZapLogger) Info(args ...interface{}) {
	z.sugar.Info(args...)
}

// Notice logs it's arguments at log level *Info*
func (z *ZapLogger) Notice(args ...interface{}) {
	z.sugar.Info(args...)
}

// Debug logs it's arguments at log level Debug
func (z *ZapLogger) Debug(args ...interface{}) {
	z.sugar.Debug(args...)
}
