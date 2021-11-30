/*
 * The package implements a basic examle of a geometry grpc server
 */

package geomserver

/*
 * Geometry Server Package
 */
import (
	"context"
	"sync"
	"time"

	"github.com/muzammilar/geomrpc/pkg/grpcserver"
	"github.com/muzammilar/geomrpc/protos/shapecalc"
	"github.com/sirupsen/logrus"
)

/*
 * Public Functions
 */

// Serve creates a gRPC Geometry Server, and registers all the required endpoints. This function is blocking (and should run in a go routine).
func Serve(wg *sync.WaitGroup, ctx context.Context, port int, certFile string, keyFile string, version string, logger *logrus.Logger) {

	// if there's a wait group implemented, then notify about the thread finishing
	if wg != nil {
		defer wg.Done()
	}

	// create a grpc server
	serverRegistrar := grpcserver.CreateServerWithStatsAndTLS(certFile, keyFile, logger)

	// geometry handler Wait Group
	var geomWg sync.WaitGroup

	// create a handler for geometry server
	geometryHandler := &GeometryServer{
		logger:  logger,
		wg:      &geomWg,
		name:    "geometry",
		version: version,
	}

	// register the geometry handler with the server for both services
	shapecalc.RegisterGeometryServer(serverRegistrar, geometryHandler)
	shapecalc.RegisterInfoServer(serverRegistrar, geometryHandler)

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

	// wait for store to finish reading
	geomWg.Wait()

	// server is shutdown
	logger.Infof("gRPC Server has shutdown: %#v", serverRegistrar)

}

/*
 * Private Functions
 */
