/*
 * The package implements a basic examle of a couple of data grpc server
 */

package dataserver

/*
 * Data Server Package
 */

import (
	"context"
	"sync"
	"time"

	"github.com/muzammilar/geometric-shapes/pkg/grpcserver"
	"github.com/muzammilar/geometric-shapes/protos/shapestore"
	"github.com/sirupsen/logrus"
)

/*
 * Public Functions
 */

// Serve creates a gRPC Data Server, and registers all the required endpoints. This function is blocking (and should run in a go routine).
func Serve(wg *sync.WaitGroup, ctx context.Context, port int, certFile string, keyFile string, version string, logger *logrus.Logger) {

	// if there's a wait group implemented, then notify about the thread finishing
	if wg != nil {
		defer wg.Done()
	}

	// create a grpc server
	serverRegistrar := grpcserver.CreateServerWithStatsAndTLS(certFile, keyFile, logger)

	// store server Wait Group
	var storeWg sync.WaitGroup

	// create a handler for store service
	storeHandler := &StoreServer{
		logger: logger,
		wg:     &storeWg,
	}

	// create a handler for generator service
	generatorHandler := &GeneratorServer{
		logger: logger,
		wg:     &storeWg,
	}

	// register the handlers with the server
	shapestore.RegisterGeneratorServer(serverRegistrar, generatorHandler)
	shapestore.RegisterStoreServer(serverRegistrar, storeHandler)

	// create a tcp listener
	listener := grpcserver.TCPListener("", port)

	// register the shutdown handler
	go func() {
		<-ctx.Done() // blocking wait
		logger.Infof("Initiating gRPC Server of shutdown; %#v", serverRegistrar)
		grpcserver.ShutDownServerWithTimeout(serverRegistrar, 20*time.Second)
	}()

	logger.Infof("Starting gRPC Server: %#v", serverRegistrar)

	// start the registrar server/registrar (blocking)
	if err := serverRegistrar.Serve(listener); err != nil {
		logger.Errorf("gRPC Server '%T' failed to serve on the listener with err: %s", serverRegistrar, err)
	}
	// server is shutdown
	logger.Errorf("gRPC Server has shutdown: %#v", serverRegistrar)

}

/*
 * Private Functions
 */
