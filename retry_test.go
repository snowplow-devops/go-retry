//
// Copyright (c) 2021-2026 Snowplow Analytics Ltd. All rights reserved.
//
// This program is licensed to you under the Apache License Version 2.0,
// and you may not use this file except in compliance with the Apache License Version 2.0.
// You may obtain a copy of the Apache License Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0.
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the Apache License Version 2.0 is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the Apache License Version 2.0 for the specific language governing permissions and limitations there under.
//

package retry

import (
	"errors"
	"io"
	"os"
	"testing"
	"time"

	pkgerrors "github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	testSleep  = time.Millisecond
	testPrefix = "test prefix"
)

var errTest = errors.New("test error")

func TestMain(m *testing.M) {
	logrus.SetOutput(io.Discard)
	os.Exit(m.Run())
}

// failingFunc returns a function that fails for the first failures calls and
// then succeeds, along with a pointer to the number of times it was called.
func failingFunc(failures int) (func() error, *int) {
	calls := 0
	return func() error {
		calls++
		if calls <= failures {
			return errTest
		}
		return nil
	}, &calls
}

func TestExponential(t *testing.T) {
	tests := []struct {
		name      string
		attempts  int
		failures  int
		wantCalls int
		wantErr   bool
	}{
		{name: "succeeds first time", attempts: 3, failures: 0, wantCalls: 1},
		{name: "succeeds after retries", attempts: 3, failures: 2, wantCalls: 3},
		{name: "fails on every attempt", attempts: 3, failures: 5, wantCalls: 3, wantErr: true},
		{name: "single attempt fails without retry", attempts: 1, failures: 1, wantCalls: 1, wantErr: true},
		{name: "zero attempts still calls once", attempts: 0, failures: 1, wantCalls: 1, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, calls := failingFunc(tt.failures)

			err := Exponential(tt.attempts, testSleep, testPrefix, f)

			if *calls != tt.wantCalls {
				t.Errorf("calls = %d, want %d", *calls, tt.wantCalls)
			}
			if !tt.wantErr {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}
			assertWrapped(t, err)
		})
	}
}

func TestExponentialWithInterface(t *testing.T) {
	tests := []struct {
		name      string
		attempts  int
		failures  int
		wantCalls int
		wantErr   bool
	}{
		{name: "succeeds first time", attempts: 3, failures: 0, wantCalls: 1},
		{name: "succeeds after retries", attempts: 3, failures: 2, wantCalls: 3},
		{name: "fails on every attempt", attempts: 3, failures: 5, wantCalls: 3, wantErr: true},
		{name: "single attempt fails without retry", attempts: 1, failures: 1, wantCalls: 1, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, calls := failingFunc(tt.failures)

			res, err := ExponentialWithInterface(tt.attempts, testSleep, testPrefix, func() (interface{}, error) {
				if err := f(); err != nil {
					return "partial", err
				}
				return 42, nil
			})

			if *calls != tt.wantCalls {
				t.Errorf("calls = %d, want %d", *calls, tt.wantCalls)
			}
			if !tt.wantErr {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if got, ok := res.(int); !ok || got != 42 {
					t.Errorf("result = %v, want 42", res)
				}
				return
			}
			if res != nil {
				t.Errorf("result = %v, want nil on failure", res)
			}
			assertWrapped(t, err)
		})
	}
}

func TestExponentialBacksOff(t *testing.T) {
	f, _ := failingFunc(3)

	start := time.Now()
	err := Exponential(4, testSleep, testPrefix, f)
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Sleeps are at least 1x, 2x and 4x the initial sleep, before jitter.
	if minWait := 7 * testSleep; elapsed < minWait {
		t.Errorf("elapsed = %s, want at least %s", elapsed, minWait)
	}
}

func TestNoSleep(t *testing.T) {
	for _, sleep := range []time.Duration{0, -time.Second} {
		t.Run(sleep.String(), func(t *testing.T) {
			f, calls := failingFunc(5)
			assertWrapped(t, Exponential(3, sleep, testPrefix, f))
			if *calls != 3 {
				t.Errorf("Exponential calls = %d, want 3", *calls)
			}

			f, calls = failingFunc(5)
			_, err := ExponentialWithInterface(3, sleep, testPrefix, func() (interface{}, error) {
				return nil, f()
			})
			assertWrapped(t, err)
			if *calls != 3 {
				t.Errorf("ExponentialWithInterface calls = %d, want 3", *calls)
			}
		})
	}
}

func assertWrapped(t *testing.T, err error) {
	t.Helper()

	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if want := testPrefix + ": " + errTest.Error(); err.Error() != want {
		t.Errorf("error = %q, want %q", err.Error(), want)
	}
	if !errors.Is(err, errTest) {
		t.Errorf("error does not wrap the original error")
	}
	if pkgerrors.Cause(err) != errTest {
		t.Errorf("errors.Cause = %v, want %v", pkgerrors.Cause(err), errTest)
	}
}
