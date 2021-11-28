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
	"github.com/muzammilar/geometric-shapes/protos/shape"
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

	// store service and worker service Wait Group
	var storeWg, workerWg sync.WaitGroup

	// create channels for storing the data
	// pointer channels are lightweight (as compared to struct)
	// however, the data is not necessarily thread safe generally (since the channel copies the pointer only and not the struct)
	// this is okay for unidirectional flow of pointers (i.e. one routine writes, others read)
	chCuboid := make(chan *shape.Cuboid, workerChanSize)

	// start the workers for handling storage
	workerWg.Add(numWorkersStore)
	for i := 0; i < numWorkersStore; i++ {
		go cuboidStoreWorker(&workerWg, i, chCuboid, logger)
	}

	// create a grpc server
	serverRegistrar := grpcserver.CreateServerWithStatsAndTLS(certFile, keyFile, logger)

	// create a handler for store service
	storeHandler := &StoreServer{
		logger:   logger,
		wg:       &storeWg,
		chCuboid: chCuboid,
	}

	// create a handler for generator service
	generatorHandler := &GeneratorServer{
		ctx:    ctx,
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

	// wait for store to finish reading
	//since the server's serve is blocking (and shutdown closes all workers), a wait group is not needed
	storeWg.Wait()

	// server is shutdown
	logger.Infof("gRPC Server has shutdown: %#v", serverRegistrar)

	// close the channels
	close(chCuboid)

	// wait for internal workers to finish
	logger.Infof("Waiting for internal storage workers to finish: %#v", workerWg)
	workerWg.Wait()
	logger.Infof("Waiting for internal storage workers to finish: %#v", workerWg)
}

/*
 * Private Functions
 */

func cuboidStoreWorker(wg *sync.WaitGroup, id int, chCuboid <-chan *shape.Cuboid, logger *logrus.Logger) {

	// notify the wait group
	defer wg.Done()

	for cuboid := range chCuboid { // blocking read (and will break on channel close)
		logger.Tracef("Database action taken on cuboid from worker-%d: %#v", id, cuboid)
	}

}
