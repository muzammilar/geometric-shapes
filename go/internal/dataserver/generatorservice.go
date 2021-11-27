/*
 * The package implements a basic examle of a couple of data grpc server
 */

package dataserver

/*
 * Data Server Package
 */
import (
	"sync"

	"github.com/muzammilar/geometric-shapes/protos/shapestore"
	"github.com/sirupsen/logrus"
)

/*
 * Private Structs (or Public)
 */

// GeneratorServer implements the Generator service.
type GeneratorServer struct {
	// embed the Store Server
	shapestore.UnimplementedGeneratorServer
	// Other internal use variables
	logger *logrus.Logger  // a shared logger (can be a bottleneck)
	wg     *sync.WaitGroup // a wait group to track all the request
}

/*
 * Functions
 */
