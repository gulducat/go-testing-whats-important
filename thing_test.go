package main

import (
	"log"
	"os"
	"path"
	"testing"

	"github.com/shoenig/test"
	"github.com/shoenig/test/must"
)

// The default constructor should get the standard writer.
// Here we assert that it's configured a specific way.
// Then, other tests can mock away those implementation details
// while this ensures that the usual happy path is as desired.
func TestNewThing(t *testing.T) {
	thing := NewThing("hi")
	test.Eq(t, "hi", thing.Msg)
	w := thing.Writer.(*os.File)
	test.Eq(t, os.Stdout, w)
}

// TestThing has numbered comments to indicate the broad general components
// of test structure. After setup that applies to the whole test,
// each subtest does:
// 1. its own setup (with cleanup)
// 2. do the thing
// 3. assert results
// These being clearly delineated helps clarify intent for human readers.
func TestThing(t *testing.T) {
	// -1. any guards?
	if os.Getenv("SKIP_THING") != "" {
		t.Skip("skipping thing because the environment")
	}

	// 0. top-level setup - cleans up after itself
	t.Setenv("IMAGINATION", "vivid")
	tmpDir := t.TempDir()
	// sometimes quite a bit of stuff goes here

	// sub-tests are great, even laid out serially like this
	t.Run("logfile", func(t *testing.T) {
		// 1. setup
		f := newTestFile(t, tmpDir, "logfile.log")
		thing := Thing{
			Msg: "hi logfile",
			// LogFile ties behavior directly to a filesystem.
			LogFile: f.Name(),
		}
		t.Logf("thing: %#v", thing)

		// 2. do the important thing
		err := thing.WriteLogFile()

		// 3. assert stuff about the results
		must.NoError(t, err) // must ~= require
		// since we wrote a real file, we need to check it
		test.FileContains(t, thing.LogFile, thing.Msg) // test ~= assert

		// no cleanup needed, because that's attached to setup
	})

	t.Run("logger", func(t *testing.T) {
		// 1. setup it
		f := newTestFile(t, tmpDir, "logger.log")
		logger := log.New(f, "logger", log.Flags())
		thing := Thing{
			Msg: "hi logger",
			// Logger is much more flexible, but still a specific type
			Logger: logger,
		}
		t.Logf("thing: %#v", thing)

		// 2. do it
		thing.WriteLogger()

		// 3. assert it
		test.FileContains(t, f.Name(), thing.Msg)
	})

	t.Run("writer", func(t *testing.T) {
		// 1. setup it
		writer := &MockWriteStorer{}
		thing := Thing{
			Msg: "hi writer",
			// using a mock implementing the Writer interface indicates that
			// we don't actually care about disks really, just Write()ing.
			Writer: writer,
		}
		t.Logf("thing: %#v", thing)

		// 2. do it
		err := thing.WriteWriter()

		// 3. assert it
		test.NoError(t, err)                       // arguably this just tests our mock.
		writer.assertWrote(t, []byte("hi writer")) // but this tests behavior.
	})
}

func newTestFile(t *testing.T, dir, name string) *os.File {
	t.Helper() // always mark your helpers as helpers.

	f, err := os.Create(path.Join(dir, name))
	// if the helper can't help, just fail the test.
	must.NoError(t, err) // "must" = t.Fatal() or t.FailNow()

	// clean up after ourselves.
	//defer f.Close() // this wouldn't work in a helper
	//t.Cleanup(f.Close) // using Cleanup() more strongly suggests error checking
	t.Cleanup(func() {
		if err := f.Close(); err != nil {
			t.Errorf("cleanup error closing file: %s", err)
		}
	})

	return f
}

// MockWriteStorer helps make assertions about Write()s,
// without getting a disk involved.
type MockWriteStorer struct {
	wrote [][]byte
	// probably should lock but meh.
}

func (w *MockWriteStorer) Write(b []byte) (int, error) {
	w.wrote = append(w.wrote, b)
	return len(b), nil
}

// assertWrote ties assertions about the Writer to the mock
func (w *MockWriteStorer) assertWrote(t *testing.T, b []byte) {
	t.Helper()
	str := string(b)
	for _, wrote := range w.wrote {
		if string(wrote) == str {
			return
		}
	}
	t.Errorf("'%s' was not written", str)
}
