/*
 * The package implements a basic examle of a couple of data grpc server
 */

package dataserver

/*
 * Data Server Package
 */
import (
	"fmt"
	"io"
	"sync"

	"github.com/muzammilar/geometric-shapes/protos/shape"
	"github.com/muzammilar/geometric-shapes/protos/shapestore"
	"github.com/sirupsen/logrus"
)

/*
 * Private Structs (or Public)
 */

// StoreServer implements the Store service.
type StoreServer struct {
	// add a context that (when closed on a graceful shutdown,) the Cuboid or AsyncReply functions can use to call a `SendAndClose` to the clients
	// embed the Store Server
	shapestore.UnimplementedStoreServer
	// Other internal use variables
	logger   *logrus.Logger  // a shared logger (can be a bottleneck)
	wg       *sync.WaitGroup // a wait group to track all the request
	chCuboid chan<- *shape.Cuboid
}

/*
 * Functions
 */

func (s *StoreServer) Cuboid(stream shapestore.Store_CuboidServer) error {
	// sample non-blocking server
	nonBlockingCh := make(chan *shape.Cuboid)
	defer close(nonBlockingCh)

	// Keep reading the stream while there is no error
	// note: the stream.stream.SendAndClose() function is not called so we can use a loop to read only
	var cuboid *shape.Cuboid
	var err error
	for cuboid, err = stream.Recv(); err == nil; cuboid, err = stream.Recv() {
		// blocking write
		s.chCuboid <- cuboid

		// non blocking write
		select {
		case nonBlockingCh <- cuboid:
			// do nothing
		default:
			// do nothing
		}

	}
	// check if the stream was closed by the client
	if err == io.EOF {
		return nil
	}

	// make sure error is handled
	if err != nil {
		return err
	}
	return nil
}

func (s *StoreServer) Replay(stream shapestore.Store_ReplayServer) error {

	// using a channel and avoiding a wait group
	errCh := make(chan error, 2)
	dataCh := make(chan *shape.Identifier)
	// reader
	go func() {
	readerLoop:
		for {
			in, err := stream.Recv()

			// read done.
			if err == io.EOF {
				errCh <- nil
				close(dataCh)
				return
			}
			// error handling
			if err != nil {
				errCh <- fmt.Errorf("Failed to receive shape info: %v", err)
			}
			// got some data!
			select {
			case dataCh <- in:
			default:
				// looks like data channel is closed by the writer
				errCh <- nil
				break readerLoop
			}

		}
	}()

	// writer
	go func() {
		var err error
		for data := range dataCh {
			if err = stream.Send(data); err != nil {
				// close the data channel (probably not a good idea to do it here)
				errCh <- fmt.Errorf("Failed to send shape info: %v", err)
				close(dataCh)
				break
			}
		}
		// clean exit on no error
	}()

	// There are two times that this function can be called (probably not a good idea to call like this)
	err := <-errCh
	if err != nil {
		return err
	}
	err = <-errCh
	if err != nil {
		return err
	}

	// close error channel
	close(errCh)

	// return
	return nil
}
