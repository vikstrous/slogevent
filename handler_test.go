package slogevent_test

import (
	"context"
	"log/slog"
	"testing"

	"github.com/vikstrous/slogevent"
)

func TestSlogevent(t *testing.T) {
	var eventRecord slogevent.Event
	l := slog.New(slogevent.NewHandler(func(ctx context.Context, e slogevent.Event) {
		eventRecord = e
	}, nil))

	l.With("key", "value").WithGroup("g").With("ingroup", true).Info("hello world", "inline", 1)

	if eventRecord.Level != slog.LevelInfo {
		t.Fatalf("wrong level: %s", eventRecord.Level)
	}
	if eventRecord.Message != "hello world" {
		t.Fatalf("unexpected message: %s", eventRecord.Message)
	}
	attrs := eventRecord.Attrs
	if len(attrs) != 3 {
		t.Fatalf("unexpected number of attributes: %d", len(attrs))
	}
	if attrs["g.ingroup"].Bool() != true {
		t.Fatalf("wrong value for g.ingroup: %#v", attrs["g.ingroup"])
	}
	if attrs["g.inline"].Int64() != 1 {
		t.Fatalf("wrong value for g.inline: %#v", attrs["g.inline"])
	}
	if attrs["key"].String() != "value" {
		t.Fatalf("wrong value for key: %#v", attrs["key"])
	}
}
