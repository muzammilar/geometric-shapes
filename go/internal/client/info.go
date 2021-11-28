/*
 * The package implements a basic go client for both geomserver and dataserver
 */

package client

/*
 * Client Package
 */
import (
	"context"
	"fmt"
	"sync"

	"github.com/muzammilar/geometric-shapes/protos/shapestore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

/*
 * Public Functions
 */

// Create an Info Service client and perform the functions required
func InfoClient(wg *sync.WaitGroup, addr string, creds credentials.TransportCredentials, ctx context.Context) {

	// notify the wait group when client finishes
	defer wg.Done()

	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	storeClient := shapestore.NewStoreClient(conn)

	fmt.Println(storeClient, ctx)
}
