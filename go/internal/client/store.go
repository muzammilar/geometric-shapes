/*
 * The package implements a basic go client for both geomserver and dataserver
 */

package client

/*
 * Client Package
 */
import (
	"sync"

	"github.com/muzammilar/geometric-shapes/protos/shapestore"
	"google.golang.org/grpc"
)

/*
 * Public Functions
 */

// Create a Store Service client and perform the functions required
func StoreClient(c *ServiceClient) {

	// notify the wait group when client finishes
	defer c.wg.Done()

	// create a connection
	//conn, err := grpc.Dial(c.addr, grpc.WithTransportCredentials(c.creds))
	conn, err := grpc.Dial(c.addr, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// setup a shared wait group
	var wg sync.WaitGroup

	// setup a service client
	storeClient := shapestore.NewStoreClient(conn)
	c.logger.Infof("Store service client: %#v", storeClient)

	wg.Add(1)
	go func() {

		wg.Done()
	}()

	// wait for internal wait groups
	wg.Wait()
}
