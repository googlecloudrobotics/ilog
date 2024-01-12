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
	"io"
	"log/slog"
)

// NewLogHandler creates a new LogHandler. Only log messages with log level l or
// higher will get written to the writer.
func NewLogHandler(l slog.Level, w io.Writer) *slog.JSONHandler {
	return slog.NewJSONHandler(w, &slog.HandlerOptions{
		AddSource: true,
		Level:     l,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			switch a.Key {
			case slog.MessageKey:
				a.Key = "message"
			case slog.TimeKey:
				a.Key = "timestamp"
			case slog.LevelKey:
				a.Key = "severity"
			}
			return a
		},
	})
}
