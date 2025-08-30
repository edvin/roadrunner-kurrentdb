package rrkurrentdb

import (
	"context"
	"fmt"
	"time"

	"github.com/kurrent-io/KurrentDB-Client-Go/kurrentdb"
)

type AppendToStreamRequest struct {
	// The stream to read from
	Stream string `json:"stream"`
	// A length of time to use for gRPC deadlines.
	Deadline *time.Duration `json:"deadline"`
	// Requires the request to be performed by the leader of the cluster.
	RequiresLeader bool         `json:"requiresLeader"`
	StreamState    *StreamState `json:"streamState"`
	// Events
	Events []kurrentdb.EventData `json:"events"`
}

func (request *AppendToStreamRequest) GetStreamState() kurrentdb.StreamState {
	switch request.StreamState.Kind {
	case "StreamExists":
		return kurrentdb.StreamExists{}
	case "NoStream":
		return kurrentdb.NoStream{}
	case "Any":
		return kurrentdb.Any{}
	case "Revision":
		if request.StreamState.Revision == nil {
			return kurrentdb.Any{}
		}
		return kurrentdb.Revision(*request.StreamState.Revision)
	default:
		panic(fmt.Sprintf("unknown streamState %q", request.StreamState))
	}
}

func (rpc *RPC) AppendToStream(in *AppendToStreamRequest, out *kurrentdb.WriteResult) error {
	rpc.plugin.log.Debug(fmt.Sprintf("About to writ: %+v", in))
	writeResult, err := rpc.plugin.Client.AppendToStream(
		context.Background(),
		in.Stream,
		kurrentdb.AppendToStreamOptions{
			StreamState:    in.GetStreamState(),
			Deadline:       in.Deadline,
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
