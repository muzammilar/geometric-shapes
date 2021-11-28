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

	"github.com/muzammilar/geometric-shapes/pkg/geomgenerator"
	"github.com/muzammilar/geometric-shapes/pkg/grpcserver"
	"github.com/muzammilar/geometric-shapes/protos/shapestore"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/emptypb"
)

/*
 * Private Structs (or Public)
 */

// GeneratorServer implements the Generator service.
type GeneratorServer struct {
	// add a context that (when closed on a graceful shutdown,) the Cuboid function can use to close the clients
	ctx context.Context
	// embed the Store Server
	shapestore.UnimplementedGeneratorServer
	// Other internal use variables
	logger *logrus.Logger  // a shared logger (can be a bottleneck)
	wg     *sync.WaitGroup // a wait group - since the server's serve is blocking (and shutdown closes all workers), a wait group is not needed
}

/*
 * Functions
 */

func (g *GeneratorServer) Cuboid(e *emptypb.Empty, stream shapestore.Generator_CuboidServer) error {

	// get client information
	clientAddr := grpcserver.GetRemoteHostFromContext(stream.Context())
	errTemplate := "Error occured while streaming GeneratorServer/Cuboid to '%s': %s"
	g.logger.Debugf("Recieved a connection from '%s' for '%T'", clientAddr, stream)

cuboidLoop:
	for {
		select {
		case <-g.ctx.Done():
			break cuboidLoop // server is shutting down
			// we can also use `return nil` but it's a PoC
		default:
			// do nothing
		}
		// send data to the stream
		c := geomgenerator.Cuboid()
		if err := stream.Send(c); err != nil {
			g.logger.Infof(errTemplate, clientAddr, err.Error())
			return err
		}

		// wait a while before sending more data
		time.Sleep(time.Duration(streamSendDelayMs) * time.Millisecond)
	}

	// no error
	return nil
}
