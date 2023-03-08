package posixsignal

import (
	"golang-standards-project-example/pkg/shutdown"
	"os"
	"os/signal"
	"syscall"
)

// Name defines shutdown manager name.
const Name = "PosixSignalManager"

// PosixSignalManager implements ShutdownManager interface that is added
// to GracefulShutdown. Initialize with NewPosixSignalManager.
type PosixSignalManager struct {
	signals []os.Signal
}

// NewPosixSignalManager initializes the PosixSignalManager.
// As arguments you can provide os.Signal-s to listen to, if none are given,
// it will default to SIGINT and SIGTERM.
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

// GetName returns name of this ShutdownManager.
func (p *PosixSignalManager) GetName() string {
	return Name
}

// Start starts listening for posix signals.
func (p *PosixSignalManager) Start(gs shutdown.GSInterface) error {
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, p.signals...)
		// Block until a signal is received.
		<-c
		gs.StartShutdown(p)
	}()
	return nil
}

// ShutdownStart does nothing.
func (p *PosixSignalManager) ShutdownStart() error {
	return nil
}

// ShutdownFinish exits the app with os.Exit(0).
func (p *PosixSignalManager) ShutdownFinish() error {
	os.Exit(0)
	return nil
}
