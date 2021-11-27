/*
 * The package implements a signal handler for geometry services
 */

package sighandler

/*
 * HTTP Server Package
 */
import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
)

// SignalHandler function takes in a context cancel function (that it calls).
// This function must be run in the main program routine and not as a sub routine
// Since it is blocking, it should be called after starting the ain routines
func SignalHandler(cancel context.CancelFunc, logger *logrus.Logger) int {

	// list of signals to handle
	signals := []os.Signal{syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM}
	logger.Info("Starting Signal Handler for the signals: ", signals)

	// a channel that is used to notify that an exit signal is recieved
	exitChan := make(chan int)

	// a channel to read signals
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, signals...)
	go func() {
		for sig := range sigs { // `range` allows for blocking reads
			logger.Infof("Signal Received: %s", sig.String())
			switch sig {
			case syscall.SIGHUP:
				logger.Infof("Ignoring Received Signal: %s", sig.String())
			case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				logger.Infof("Respecting Shutdown Signal: %s", sig.String())
				// create an exit code
				exitChan <- 0
				// cancel the context
				cancel()
			default: // always need to be specified
			}
		}
	}()

	// wait here for exit code
	exitCode := <-exitChan
	close(exitChan) // close this channel in cas anyone is listening on the channel

	return exitCode
}
