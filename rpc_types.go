package rrkurrentdb

import (
	"encoding/json"
	"fmt"
)

type From struct {
	Kind  string
	Index *int
}

func (f *From) UnmarshalJSON(data []byte) error {
	// Try int first
	var i int
	if err := json.Unmarshal(data, &i); err == nil {
		f.Kind = "index"
		f.Index = &i
		return nil
	}

	// Try string
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		switch s {
		case "start", "end":
			f.Kind = s
			f.Index = nil
			return nil
		default:
			return fmt.Errorf("invalid string value for From: %q", s)
		}
	}

	return fmt.Errorf("from must be int or \"start\"/\"end\", got: %s", string(data))
}

type StreamState struct {
	Kind     string
	Revision *uint64
}

func (f *StreamState) UnmarshalJSON(data []byte) error {
	// Try int first
	var i uint64
	if err := json.Unmarshal(data, &i); err == nil {
		f.Kind = "Revision"
		f.Revision = &i
		return nil
	}
	
	// Try string
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		switch s {
		case "Any", "StreamExists", "NoStream":
			f.Kind = s
			f.Revision = nil
			return nil
		default:
			return fmt.Errorf("invalid string value for From: %q", s)
		}
	}

	return fmt.Errorf("StreamState must be int (revision) or \"Any\"/\"StreamExists\" or \"NoStream\", got: %s", string(data))
}
