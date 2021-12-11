package main

import "github.com/godofcc/go-common/lib/log"

func main() {
	opts := &log.Options{
		Level:            "debug",
		Format:           "json",
		EnableColor:      false,
		OutputPaths:      []string{"test.log", "stdout"},
		ErrorOutputPaths: []string{"error.log"},
	}
	log.Init(opts)
	log.Info("hello")
	log.Infow("Message printed with Errorw", "X-Request-ID", "fbf54504-64da-4088-9b86-67824a7fb508")

}
