package rrkurrentdb

import (
	"context"
	"errors"
	"io"
	"math"
	"time"

	"github.com/kurrent-io/KurrentDB-Client-Go/kurrentdb"
)

type ReadStreamRequest struct {
	// The stream to read from
	Stream string `json:"stream"`
	// The maximum number of events to read.
	MaxEvents *uint64 `json:"maxEvents"`
	// Direction to read in the stream
	Direction int `json:"direction"`
	// Starting position of the read request.
	From *From `json:"from"`
	// Whether the read request should resolve linkTo events to their linked events.
	ResolveLinkTos bool `json:"resolveLinkTos"`
	// A length of time to use for gRPC deadlines.
	Deadline *time.Duration `json:"deadline"`
	// Requires the request to be performed by the leader of the cluster.
	RequiresLeader bool `json:"requiresLeader"`
}

func (r *ReadStreamRequest) toReadStreamOptions() kurrentdb.ReadStreamOptions {
	return kurrentdb.ReadStreamOptions{
		Direction:      kurrentdb.Direction(r.Direction),
		From:           r.GetFrom(),
		ResolveLinkTos: r.ResolveLinkTos,
		Deadline:       r.Deadline,
		RequiresLeader: r.RequiresLeader,
	}
}

func (r *ReadStreamRequest) GetFrom() kurrentdb.StreamPosition {
	if r.From == nil {
		return kurrentdb.Start{}
	}

	switch r.From.Kind {
	case "start":
		return kurrentdb.Start{}
	case "end":
		return kurrentdb.End{}
	case "index":
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
