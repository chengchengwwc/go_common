package shutdown

import (
	"os"
	"os/signal"
	"syscall"
)

const PosixSignalManagerName = "PosixSignalManager"

type PosixSignalManager struct {
	signals []os.Signal
}

func NewPosixSignalManager(sig ...os.Signal) *PosixSignalManager {
	if len(sig) == 0 {
		sig = make([]os.Signal, 2)
		sig[0] = os.Interrupt
		sig[1] = syscall.SIGTERM
	}
	return &PosixSignalManager{
		signals: sig,
	}
}

func (p *PosixSignalManager) GetName() string {
	return PosixSignalManagerName
}

func (p *PosixSignalManager) Start(gs GSInterface) error {
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, p.signals...)
		<-c
		gs.StartShutdown(p)
	}()
	return nil
}

func (p *PosixSignalManager) ShutdownStart() error {
	return nil
}

// ShutdownFinish exits the app with os.Exit(0).
func (p *PosixSignalManager) ShutdownFinish() error {
	os.Exit(0)
	return nil
}
