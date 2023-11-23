// Copyright 2023 The Cloud Robotics Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ilog

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"sync"
	"time"
)

type groupOrAttrs struct {
	group string
	attrs []slog.Attr
}

// LogHandler implements slog.Handler.
// Use NewLogHandler to create instances with correct internal state.
type LogHandler struct {
	level  slog.Level
	goas   []groupOrAttrs
	mu     *sync.Mutex
	writer io.Writer
}

// Enabled is part of the slog.Handler interface
func (h *LogHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level
}

// Handle is part of the slog.Handler interface
func (h *LogHandler) Handle(_ context.Context, r slog.Record) error {
	var buf bytes.Buffer
	buf.Grow(1024) // should cover most log messages
	if !r.Time.IsZero() {
		buf.WriteString(r.Time.UTC().Format(time.RFC3339))
		buf.WriteRune(' ')
	}
	buf.WriteString(r.Level.String())
	buf.WriteRune(' ')
	buf.WriteString(r.Message)
	buf.WriteRune(' ')
	// TODO: The layout of groups and extra attrs needs some work. So far they
	// have not been used so it does not matter.
	for _, goa := range h.goas {
		if goa.group != "" {
			buf.WriteString(goa.group)
			buf.WriteString(": ")
		} else {
			for _, a := range goa.attrs {
				h.appendAttr(&buf, a)
			}
		}
	}

	r.Attrs(func(a slog.Attr) bool {
		h.appendAttr(&buf, a)
		return true
	})
	buf.WriteRune('\n')
	h.mu.Lock()
	defer h.mu.Unlock()
	_, err := h.writer.Write(buf.Bytes())
	return err
}

func (h *LogHandler) appendAttr(w io.Writer, a slog.Attr) (int, error) {
	if a.Equal(slog.Attr{}) {
		return 0, nil
	}
	switch a.Value.Kind() {
	case slog.KindString:
		return fmt.Fprintf(w, "%s=%s ", a.Key, a.Value.String())
	case slog.KindTime:
		return fmt.Fprintf(w, "%s=%s ", a.Key, a.Value.Time().UTC().Format(time.RFC3339))
	case slog.KindGroup:
		// TODO: how to handle KindGroup?
	default:
		// slog.KindAny:
		// slog.KindBool:
		// slog.KindDuration:
		// slog.KindFloat64:
		// slog.KindInt64:
		// slog.KindUint64:
		// slog.KindLogValuer:
		return fmt.Fprintf(w, "%s=%s ", a.Key, a.Value)
	}
	return 0, nil
}

func (h *LogHandler) withGroupOrAttrs(goa groupOrAttrs) *LogHandler {
	h2 := *h
	h2.goas = make([]groupOrAttrs, len(h.goas)+1)
	copy(h2.goas, h.goas)
	h2.goas[len(h2.goas)-1] = goa
	return &h2
}

// WithGroup is part of the slog.Handler interface
func (h *LogHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	return h.withGroupOrAttrs(groupOrAttrs{group: name})
}

// WithAttrs is part of the slog.Handler interface
func (h *LogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return h
	}
	return h.withGroupOrAttrs(groupOrAttrs{attrs: attrs})
}

// NewLogHandler creates a new LogHandler. Only log messages with log level l or
// higher will get written to the writer.
func NewLogHandler(l slog.Level, w io.Writer) *LogHandler {
	return &LogHandler{
		level:  l,
		mu:     &sync.Mutex{},
		writer: w,
	}
}
