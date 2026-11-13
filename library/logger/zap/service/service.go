package service

import (
	"io"

	"github.com/air-go/rpc/library/logger"
	"github.com/air-go/rpc/library/logger/zap"
)

type Config struct {
	InfoFile  string
	ErrorFile string
	Level     string
}

type ServiceLogger struct {
	*zap.ZapLogger
	config *Config
}

type options struct {
	infoWriter  io.Writer
	errorWriter io.Writer
}

type Option func(l *options)

func defaultOptions() *options {
	return &options{}
}

// WithWriter sets both log destinations, overriding the rotating files created
// from Config.InfoFile and Config.ErrorFile; a single writer is rejected.
func WithWriter(infoWriter, errorWriter io.Writer) Option {
	return func(o *options) {
		o.infoWriter, o.errorWriter = infoWriter, errorWriter
	}
}

func NewServiceLogger(serviceName string, config *Config, options ...Option) (*ServiceLogger, error) {
	opts := defaultOptions()
	for _, o := range options {
		o(opts)
	}

	infoWriter, errorWriter, err := logger.ResolveWriters(opts.infoWriter, opts.errorWriter, config.InfoFile, config.ErrorFile)
	if err != nil {
		return nil, err
	}

	zl, err := zap.NewLogger(
		zap.WithModule(logger.ModuleHTTP),
		zap.WithServiceName(serviceName),
		zap.WithInfoWriter(infoWriter),
		zap.WithErrorWriter(errorWriter),
		zap.WithLevel(config.Level),
	)
	if err != nil {
		return nil, err
	}

	return &ServiceLogger{ZapLogger: zl, config: config}, nil
}
