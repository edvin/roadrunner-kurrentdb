package rrkurrentdb

import (
	"context"
	"errors"
	"io"
	"math"

	"github.com/kurrent-io/KurrentDB-Client-Go/kurrentdb"
)

type ReadStreamRequest struct {
	// The stream to read from
	Stream string `msgpack:"Stream"`
	// The maximum number of events to read.
	MaxEvents *uint64 `msgpack:"MaxEvents"`
	// Direction to read in the stream
	Direction int `msgpack:"Direction"`
	// Starting position of the read request.
	From *From `msgpack:"From"`
	// Whether the read request should resolve linkTo events to their linked events.
	ResolveLinkTos bool `msgpack:"ResolveLinkTos"`
	// A length of time to use for gRPC deadlines.
	Deadline *int64 `msgpack:"Deadline"`
	// Requires the request to be performed by the leader of the cluster.
	RequiresLeader bool `msgpack:"RequiresLeader"`
}

func (r *ReadStreamRequest) toReadStreamOptions() kurrentdb.ReadStreamOptions {
	return kurrentdb.ReadStreamOptions{
		Direction:      kurrentdb.Direction(r.Direction),
		From:           r.GetFrom(),
		ResolveLinkTos: r.ResolveLinkTos,
		Deadline:       deadlineMsToDuration(r.Deadline),
		RequiresLeader: r.RequiresLeader,
	}
}

func (r *ReadStreamRequest) GetFrom() kurrentdb.StreamPosition {
	if r.From == nil {
		return kurrentdb.Start{}
	}

	switch r.From.Kind {
	case "Start":
		return kurrentdb.Start{}
	case "End":
		return kurrentdb.End{}
	case "Index":
		return kurrentdb.Revision(uint64(*r.From.Index))
	default:
		panic("invalid from")
	}
}

func (r *ReadStreamRequest) GetCount() uint64 {
	if r.MaxEvents == nil {
		return math.MaxUint64
	} else {
		return *r.MaxEvents
	}
}

func (r *ReadStreamRequest) validateReadStreamRequest() error {
	if int(kurrentdb.Forwards) == r.Direction && r.From.Kind == "end" {
		return errors.New("cannot seek forwards from the end of the stream")
	}
	if int(kurrentdb.Backwards) == r.Direction && r.From.Kind == "start" {
		return errors.New("cannot seek backwards from the beginning of the stream")
	}
	return nil
}

func (rpc *RPC) ReadStream(in ReadStreamRequest, out *[]kurrentdb.ResolvedEvent) error {
	err := in.validateReadStreamRequest()
	if err != nil {
		return err
	}

	stream, err := rpc.plugin.Client.ReadStream(
		context.Background(),
		in.Stream,
		in.toReadStreamOptions(),
		in.GetCount(),
	)
	if err != nil {
		return err
	}

	defer stream.Close()

	var events []kurrentdb.ResolvedEvent

	for {
		event, err := stream.Recv()

		if err, ok := kurrentdb.FromError(err); !ok {
			if errors.Is(err, io.EOF) {
				break
			} else {
				return err
			}
		}

		events = append(events, *event)
	}
	*out = events
	return nil
}
