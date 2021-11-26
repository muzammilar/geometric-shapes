/*
 * The package implements a basic go client for both geomserver and dataserver
 */

package client

/*
 * HTTP Server Package
 */
import (
	"context"
	"fmt"

	"github.com/muzammilar/geometric-shapes/protos/shapestore"
	"google.golang.org/grpc"
)

/*
 * Public Functions
 */

// Create an HTTP Server and register all the required endpoints
func StartClient(addr string) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		// TODO error
	}

	defer conn.Close()

	storeClient := shapestore.NewStoreClient(conn)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fmt.Println(storeClient, ctx)
}
