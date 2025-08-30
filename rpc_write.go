package rrkurrentdb

import (
	"context"
	"fmt"

	"github.com/kurrent-io/KurrentDB-Client-Go/kurrentdb"
)

type AppendToStreamRequest struct {
	// The stream to read from
	Stream string `msgpack:"Stream"`
	// A length of time to use for gRPC deadlines.
	Deadline *int64 `msgpack:"Deadline"`
	// Requires the request to be performed by the leader of the cluster.
	RequiresLeader bool        `msgpack:"RequiresLeader"`
	StreamState    StreamState `msgpack:"StreamState"`
	// Events
	Events []kurrentdb.EventData `msgpack:"Events"`
}

func (rpc *RPC) AppendToStream(in *AppendToStreamRequest, out *kurrentdb.WriteResult) error {
	rpc.plugin.log.Debug(fmt.Sprintf("About to write: %+v", in))

	writeResult, err := rpc.plugin.Client.AppendToStream(
		context.Background(),
		in.Stream,
		kurrentdb.AppendToStreamOptions{
			StreamState:    in.StreamState.ToKurrentDBStreamState(),
			Deadline:       deadlineMsToDuration(in.Deadline),
			RequiresLeader: in.RequiresLeader,
		},
		in.Events...,
	)
	if err != nil {
		return err
	}

	*out = *writeResult
	return nil
}
