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
	"github.com/muzammilar/geometric-shapes/protos/shapestore"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/peer"
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
	wg     *sync.WaitGroup // a wait group to track all the request
}

/*
 * Functions
 */

func (g *GeneratorServer) Cuboid(e *emptypb.Empty, stream shapestore.Generator_CuboidServer) error {

	if remotePeer, ok := peer.FromContext(stream.Context()); ok {
		g.logger.Debugf("Recieved a connection from '%s' for '%T'", remotePeer.Addr.String(), stream)
	}

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
			return err
		}

		// wait a while before sending more data
		time.Sleep(time.Duration(streamSendDelayMs) * time.Millisecond)
	}

	// no error
	return nil
}
