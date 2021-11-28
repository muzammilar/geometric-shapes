/*
 * The package implements a basic go client for both geomserver and dataserver
 */

package client

/*
 * Client Package
 */
import (
	"fmt"

	"github.com/muzammilar/geometric-shapes/protos/shapestore"
	"google.golang.org/grpc"
)

/*
 * Public Functions
 */

// Create an Info Service client and perform the functions required
func InfoClient(c *ServiceClient) {

	// notify the wait group when client finishes
	defer c.wg.Done()

	conn, err := grpc.Dial(c.addr, grpc.WithTransportCredentials(c.creds))
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	storeClient := shapestore.NewStoreClient(conn)

	fmt.Println(storeClient, c.ctx)

}
