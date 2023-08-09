package slogevent

import (
	"context"
	"log/slog"
	"strings"
	"time"
)

// Event is similar to slog.Record but contains all attributes in a key value map with the groups collected as dot separated keys
type Event struct {
	Time    time.Time
	Message string
	Level   slog.Level
	PC      uintptr
	Attrs   map[string]slog.Value
}

type EventFunc func(ctx context.Context, e Event)

type Handler struct {
	eventFunc EventFunc
	with      *groupOrAttrs
	next      slog.Handler
}

func NewHandler(e EventFunc, next slog.Handler) *Handler {
	return &Handler{
		eventFunc: e,
		next:      next,
	}
}

// Enabled implements slog.Handler.
func (*Handler) Enabled(context.Context, slog.Level) bool {
	return true
}

// Handle implements slog.Handler.
func (h *Handler) Handle(ctx context.Context, r slog.Record) error {
	c := attrsCollector{
		attrs: map[string]slog.Value{},
	}
	groups := h.with.Apply(c.formatAttr)
	r.Attrs(func(a slog.Attr) bool {
		c.formatAttr(groups, a)
		return true
	})
	h.eventFunc(ctx, Event{
		Time:    r.Time,
		Message: r.Message,
		PC:      r.PC,
		Level:   r.Level,
		Attrs:   c.attrs,
	})
	if h.next != nil {
		return h.next.Handle(ctx, r)
	}
	return nil
}

// WithAttrs implements slog.Handler.
func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &Handler{h.eventFunc, h.with.WithAttrs(attrs), h.next}
}

// WithGroup implements slog.Handler.
func (h *Handler) WithGroup(name string) slog.Handler {
	return &Handler{h.eventFunc, h.with.WithGroup(name), h.next}
}

var _ = slog.Handler(&Handler{})

type attrsCollector struct {
	attrs map[string]slog.Value
}

func (at *attrsCollector) formatAttr(groups []string, a slog.Attr) {
	if a.Value.Kind() == slog.KindGroup {
		gs := a.Value.Group()
		if len(gs) == 0 {
			return
		}
		if a.Key != "" {
			groups = append(groups, a.Key)
		}
		for _, g := range gs {
			at.formatAttr(groups, g)
		}
	} else if key := a.Key; key != "" {
		if len(groups) > 0 {
			key = strings.Join(groups, ".") + "." + key
		}
		at.attrs[key] = a.Value
	}
}
