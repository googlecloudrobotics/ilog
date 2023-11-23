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
	"fmt"
	"io"
	"log/slog"
	"os"
)

type multiCloser struct {
	closer []io.Closer
}

// Close implements the io.Closer interface
func (mc *multiCloser) Close() error {
	for _, c := range mc.closer {
		if err := c.Close(); err != nil {
			return err
		}
	}
	return nil
}

// Setup changes the log output to the provided files.
// It keeps logging to /dev/stderr as well.
func Setup(fs ...string) (io.Closer, error) {
	writer := make([]io.Writer, 0, len(fs))
	mc := multiCloser{
		closer: make([]io.Closer, 0, len(fs)),
	}
	// /dev/stderr is special. If passed as argument it would fail during bootup.
	// It also should not be closed.
	// See https://pkg.go.dev/os@go1.20.6#pkg-variables
	writer = append(writer, os.Stderr)
	for _, f := range fs {
		w, err := os.OpenFile(f, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			mc.Close()
			return nil, fmt.Errorf("failed to open %s: %v", f, err)
		}
		mc.closer = append(mc.closer, w)
		writer = append(writer, w)
	}
	nh := NewLogHandler(slog.LevelInfo, io.MultiWriter(writer...))
	slog.SetDefault(slog.New(nh))
	return &mc, nil
}
