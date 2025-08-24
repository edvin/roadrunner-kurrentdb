package rrkurrentdb

import "go.uber.org/zap"

type rpc struct {
	plugin *Plugin
	log    *zap.Logger
}

func (s *rpc) Hello(input string, output *string) error {
	*output = input
	// s.plugin.Foo() <-- you may also use methods from the Plugin itself
	s.log.Debug("foo")
	return nil
}

func (s *Plugin) RPC() any {
	return &rpc{plugin: s, log: s.log}
}
