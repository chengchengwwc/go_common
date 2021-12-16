package main

import (
	"fmt"
	"github.com/godofcc/go-common/lib/shutdown"
	"time"
)

func main() {
	gs := shutdown.New()
	gs.AddShutdownManager(shutdown.NewPosixSignalManager())
	gs.AddShutdownCallback(shutdown.ShutdownFunc(func(string) error {
		fmt.Println("SSSS")
		return nil
	}))

	if err := gs.Start(); err != nil {
		fmt.Println(err)
		return
	}
	time.Sleep(time.Hour)

	//opts := &log.Options{
	//	Level:            "debug",
	//	Format:           "json",
	//	EnableColor:      false,
	//	OutputPaths:      []string{"test.log", "stdout"},
	//	ErrorOutputPaths: []string{"error.log"},
	//}
	//log.Init(opts)
	//log.Info("hello")
	//log.Infow("Message printed with Errorw", "X-Request-ID", "fbf54504-64da-4088-9b86-67824a7fb508")

}
