package events

import (
	"gitlab.petrichor.io/sandstone/v1/errors/errout"
	"gitlab.petrichor.io/sandstone/v1/errors/stackerr"
	"testing"
	"time"
)

func TestSyncEvent(t *testing.T) {
	var tests = []struct {
		fireTrigger     string
		fireValue       interface{}
		listenTrigger   string
		listenValue     interface{}
		expectedSuccess bool
	}{
		{"test", "abc", "test", "abc", true},
		{"test", "abc", "TEST", "abc", false},
		{"test", "abc", "test", "cba", false},
		{"test", 123, "test", 123, true},
	}

	for k, v := range tests {
		d := new(Dispatcher)

		var caught bool

		d.Call(func(data interface{}) error {
			caught = true

			if v.fireValue != data {
				t.Fatalf("Caught value does not match fired value (test %d)", k)
			}

			return nil
		}).On(v.listenTrigger)

		if err := d.Fire(v.fireTrigger, v.fireValue); err != nil {
			t.Fatalf("Could not fire event via dispatcher (test %d): %s", k, err.Error())
		}

		if !caught && v.expectedSuccess {
			t.Fatalf("Unable to catch event (test %d)", k)
		} else if caught && !v.expectedSuccess {
			t.Fatalf("Caught event when we shouldn't have (test %d)", k)
		}
	}
}

func TestAsyncEvent(t *testing.T) {
	var tests = []struct {
		fireTrigger     string
		fireValue       interface{}
		listenTrigger   string
		listenValue     interface{}
		expectedSuccess bool
	}{
		{"test", "abc", "test", "abc", true},
		{"test", "abc", "TEST", "abc", false},
		{"test", "abc", "test", "abc", true},
		{"test", 123, "test", 123, true},
	}

	for k, v := range tests {
		d := new(Dispatcher)

		var ok bool

		d.Call(func(data interface{}) error {
			ok = false

			if v.listenValue == data {
				ok = true
			}

			return nil
		}).On(v.listenTrigger)

		errCh := make(chan error)

		doneCh := d.FireAsync(v.fireTrigger, v.fireValue, errout.Chan(errCh))

		select {
		case <-doneCh:

		case err := <-errCh:
			if err != nil {
				t.Fatalf("Error caught during test %d", k)
			}
		}

		if !ok && v.expectedSuccess {
			t.Fatalf("Unable to catch event (test %d)", k)
		} else if ok && !v.expectedSuccess {
			t.Fatalf("Caught event when we shouldn't have (test %d)", k)
		}
	}
}

func TestAsyncChannels(t *testing.T) {
	d := NewDispatcher()

	var (
		caught   bool
		bindings int = 5
	)

	for n := 0; n < bindings; n++ {
		d.Call(func(data interface{}) error {
			caught = true

			return stackerr.New("Forced error")
		}).On("trigger")
	}

	errCh := make(chan error, bindings)

	doneCh := d.FireAsync("trigger", nil, errCh)

	select {
	case <-doneCh:
		errs := len(errCh)
		if errs != bindings {
			t.Fatalf("Only %d errors in channel, %d expected", errs, bindings)
		}
	case <-time.After(time.Second):
		t.Fatal("Timeout during event async fire")
	}

	if !caught {
		t.Fatal("Event not caught")
	}

}
