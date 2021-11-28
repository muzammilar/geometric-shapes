/*
 * The package implements a basic go client for both geomserver and dataserver
 */

package client

/*
 * Client Package
 */
import (
	"github.com/muzammilar/geometric-shapes/protos/shapecalc"
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

	// example mux
	infoClient := shapecalc.NewInfoClient(conn)
	geometryClient := shapecalc.NewGeometryClient(conn)

	c.logger.Infof("Info service client: %#v", infoClient)
	c.logger.Infof("Geometry service client: %#v", geometryClient)

}
