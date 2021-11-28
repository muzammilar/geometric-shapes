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

// Create a Generator Service client and perform the functions required
func GeneratorClient(c *ServiceClient) {

	// notify the wait group when client finishes
	defer c.wg.Done()

	conn, err := grpc.Dial(c.addr, grpc.WithTransportCredentials(c.creds))
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	// example mux
	storeClient := shapestore.NewStoreClient(conn)
	generatorClient := shapestore.NewGeneratorClient(conn)

	c.logger.Infof("Store service client: %#v", storeClient)
	c.logger.Infof("Generator service client: %#v", generatorClient)

	fmt.Println(storeClient, c.ctx)
}
