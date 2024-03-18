package filmoteka_test

import (
	"context"
	"io"
	"log"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/famusovsky/VkTestTask/internal/filmoteka"
	"github.com/stretchr/testify/assert"
)

// testServer - структура для тестирования GraceRun.
type testServer struct {
	ListenAndServeCalled bool
	ShutdownCalled       bool
}

func (ts *testServer) ListenAndServe() error {
	ts.ListenAndServeCalled = true
	for ts.ShutdownCalled == false {
		time.Sleep(time.Millisecond * 100)
	}
	return nil
}
func (ts *testServer) Shutdown(context.Context) error {
	ts.ShutdownCalled = true
	return nil
}

func TestGraceRun(t *testing.T) {
	sendTerminationSignal := func() error {
		pid := os.Getpid()
		process, err := os.FindProcess(pid)
		if err != nil {
			return err
		}

		return process.Signal(syscall.SIGTERM)
	}
	t.Run("testing graceful shutdown of the server", func(t *testing.T) {
		assert.Equal(t, 1, 1)
		logger := log.New(io.Discard, "test", log.LstdFlags)
		app := filmoteka.CreateApp(":8080", logger, logger, nil, false)
		srvr := &testServer{}

		go func() {
			filmoteka.GraceRun(srvr, app)
		}()

		time.Sleep(time.Second)
		err := sendTerminationSignal()
		time.Sleep(time.Second)

		assert.NoError(t, err)
		assert.True(t, srvr.ShutdownCalled)
		assert.True(t, srvr.ListenAndServeCalled)
	})
}
