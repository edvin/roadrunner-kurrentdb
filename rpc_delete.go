package rrkurrentdb

import (
	"context"

	"github.com/kurrent-io/KurrentDB-Client-Go/kurrentdb"
)

type DeleteStreamRequest struct {
	Stream         string      `msg:"Stream"`
	StreamState    StreamState `msgpack:"StreamState"`
	Deadline       *int64      `msgpack:"Deadline,omitempty"`
	RequiresLeader bool        `msgpack:"RequiresLeader"`
}

func (rpc *RPC) DeleteStream(in *DeleteStreamRequest, out *kurrentdb.DeleteResult) error {
	result, err := rpc.plugin.Client.DeleteStream(context.Background(), in.Stream, kurrentdb.DeleteStreamOptions{
		StreamState:    in.StreamState.ToKurrentDBStreamState(),
		Deadline:       deadlineMsToDuration(in.Deadline),
		RequiresLeader: in.RequiresLeader,
	})
	if err != nil {
		return err
	}

	*out = *result
	return nil
}
