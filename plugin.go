package rrkurrentdb

import (
	"context"
	"fmt"

	"github.com/kurrent-io/KurrentDB-Client-Go/kurrentdb"
	"github.com/roadrunner-server/errors"
	"go.uber.org/zap"
)

const PluginName = "kurrentdb"

type Configurer interface {
	UnmarshalKey(name string, out any) error
	Has(name string) bool
}

type Logger interface {
	NamedLogger(name string) *zap.Logger
}
type Plugin struct {
	cfg    *Config
	log    *zap.Logger
	Client *kurrentdb.Client
}

func (s *Plugin) Init(cfg Configurer, log Logger) error {
	op := errors.Op("kurrentdb plugin init")

	if !cfg.Has(PluginName) {
		return errors.E(op, errors.Disabled)
	}

	err := cfg.UnmarshalKey(PluginName, &s.cfg)
	if err != nil {
		return errors.E(op, err)
	}

	s.log = log.NamedLogger(PluginName)

	s.cfg.InitDefaults()
	return nil
}

func (s *Plugin) Serve() chan error {
	errCh := make(chan error, 1)

	go func() {
		defer close(errCh)

		const op = errors.Op("kurrentdb plugin serve")

		settings, err := kurrentdb.ParseConnectionString(s.cfg.Address)
		if err != nil {
			errCh <- errors.E(op, err)
			return
		}

		settings.Logger = func(level kurrentdb.LogLevel, format string, args ...interface{}) {
			msg := fmt.Sprintf(format, args...)

			switch level {
			case kurrentdb.LogDebug:
				s.log.Debug(msg)
			case kurrentdb.LogInfo:
				s.log.Info(msg)
			case kurrentdb.LogWarn:
				s.log.Warn(msg)
			case kurrentdb.LogError:
				s.log.Error(msg)
			default:
				s.log.Info(msg)
			}
		}

		db, err := kurrentdb.NewClient(settings)
		if err != nil {
			errCh <- errors.E(op, err)
			return
		}

		s.Client = db
		s.log.Info("kurrentdb client configured")
	}()

	return errCh
}

func (s *Plugin) Stop(ctx context.Context) error {
	if s.Client != nil {
		s.log.Info("closing kurrentdb connection")
		return s.Client.Close()
	}
	return nil
}

func (s *Plugin) Name() string {
	return PluginName
}

func (s *Plugin) RPC() any {

	return &RPC{plugin: s}
}
