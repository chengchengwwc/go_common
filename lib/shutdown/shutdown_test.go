package shutdown

import (
	"fmt"
	"testing"
	"time"
)

func TestShutdown(t *testing.T) {
	gs := New()
	gs.AddShutdownManager(NewPosixSignalManager())
	gs.AddShutdownCallback(ShutdownFunc(func(string) error {
		fmt.Println("SSSS")
		return nil
	}))

	if err := gs.Start(); err != nil {
		t.Log(err)
		return
	}
	time.Sleep(time.Hour)
}
