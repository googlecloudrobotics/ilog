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
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestMultiLogger(t *testing.T) {
	d := t.TempDir()
	f1 := filepath.Join(d, "t1")
	f2 := filepath.Join(d, "t2")
	f3 := filepath.Join(d, "t3")
	c, err := Setup(f1, f2, f3)
	if err != nil {
		t.Fatalf("Could not setup logger: %v", err)
	}
	log.Printf("TestTestTest")
	c.Close()
	b1, err := os.ReadFile(f1)
	if err != nil {
		t.Errorf("Failed to read f1: %v", err)
	}
	b2, err := os.ReadFile(f2)
	if err != nil {
		t.Errorf("Failed to read f2: %v", err)
	}
	b3, err := os.ReadFile(f3)
	if err != nil {
		t.Errorf("Failed to read f3: %v", err)
	}
	if len(b1) == 0 {
		t.Errorf("b1 is empty")
	}
	if !bytes.Equal(b1, b2) {
		t.Errorf("b1 != b2")
	}
	if !bytes.Equal(b2, b3) {
		t.Errorf("b2 != b3")
	}
}

// This 'test' is more of a development tool than an actual test.
// The output can only be observed and verified manually.
// It should print the log message twice: once for stderr and once for stdout.
func TestMultiLoggerStdErr(t *testing.T) {
	c, err := Setup("/dev/stdout")
	if err != nil {
		t.Fatalf("Could not setup logger: %v", err)
	}
	log.Printf("TESTtestTEST")
	c.Close()
}
