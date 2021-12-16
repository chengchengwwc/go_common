package shutdown

import "sync"

type ShutdownCallback interface {
	OnShutdown(string) error
}

type ShutdownFunc func(string) error

func (f ShutdownFunc) OnShutdown(shutdownManager string) error {
	return f(shutdownManager)
}

type ShutdownManager interface {
	GetName() string
	Start(gs GSInterface) error
	ShutdownStart() error
	ShutdownFinish() error
}

type GSInterface interface {
	StartShutdown(sm ShutdownManager)
	ReportError(err error)
	AddShutdownCallback(shutdownCallback ShutdownCallback)
}

type ErrorHandler interface {
	OnError(err error)
}

type GracefulShutdown struct {
	callbacks    []ShutdownCallback
	managers     []ShutdownManager
	errorHandler ErrorHandler
}

func New() *GracefulShutdown {
	return &GracefulShutdown{
		callbacks: make([]ShutdownCallback, 0, 10),
		managers:  make([]ShutdownManager, 0, 10),
	}
}

func (gs *GracefulShutdown) Start() error {
	for _, manger := range gs.managers {
		if err := manger.Start(gs); err != nil {
			return err
		}
	}
	return nil
}

func (gs *GracefulShutdown) AddShutdownManager(manager ShutdownManager) {
	gs.managers = append(gs.managers, manager)
}

func (gs *GracefulShutdown) AddShutdownCallback(shutdownCallback ShutdownCallback) {
	gs.callbacks = append(gs.callbacks, shutdownCallback)
}

func (gs *GracefulShutdown) SetErrorHandler(errorHandler ErrorHandler) {
	gs.errorHandler = errorHandler
}

func (gs *GracefulShutdown) StartShutdown(sm ShutdownManager) {
	gs.ReportError(sm.ShutdownStart())
	var wg sync.WaitGroup
	for _, shutdownCallback := range gs.callbacks {
		wg.Add(1)
		go func(shutdownCallback ShutdownCallback) {
			defer wg.Done()
			gs.ReportError(shutdownCallback.OnShutdown(sm.GetName()))
		}(shutdownCallback)
	}
	wg.Wait()
	gs.ReportError(sm.ShutdownFinish())

}

func (gs *GracefulShutdown) ReportError(err error) {
	if err != nil && gs.errorHandler != nil {
		gs.errorHandler.OnError(err)
	}
}
