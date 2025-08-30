package rrkurrentdb

import (
	"fmt"

	"github.com/kurrent-io/KurrentDB-Client-Go/kurrentdb"
)

type From struct {
	Kind  string
	Index *int
}

type StreamState struct {
	Kind     string
	Revision *uint64
}

func (ss *StreamState) ToKurrentDBStreamState() kurrentdb.StreamState {
	switch ss.Kind {
	case "StreamExists":
		return kurrentdb.StreamExists{}
	case "NoStream":
		return kurrentdb.NoStream{}
	case "Any":
		return kurrentdb.Any{}
	case "Revision":
		if ss.Revision == nil {
			return kurrentdb.Any{}
		}
		return kurrentdb.Revision(*ss.Revision)
	default:
		panic(fmt.Sprintf("unknown streamState kind %q", ss.Kind))
	}
}
