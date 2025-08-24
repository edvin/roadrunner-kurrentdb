package rrkurrentdb

import (
	"context"
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
	cfg *Config
	log *zap.Logger
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
	const op = errors.Op("kurrentdb plugin serve")
	s.log.Info("kurrentdb plugin serve")
	return nil
}

func (s *Plugin) Stop(ctx context.Context) error {
	return nil
}

func (s *Plugin) Name() string {
	return PluginName
}
