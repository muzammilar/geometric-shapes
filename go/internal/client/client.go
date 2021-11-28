/*
 * The package implements a basic go client for both geomserver and dataserver
 */

package client

/*
 * Client Package
 */
import (
	"context"
	"sync"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/credentials"
)

/*
 * Public Functions
 */

// Service Client Struct contains all the variables passed to all the clients
type ServiceClient struct {
	addr   string                           // server/destination address
	logger *logrus.Logger                   // logger
	ctx    context.Context                  // client context (used for propagating shutdown signals)
	wg     *sync.WaitGroup                  // shared wait group (with main routine)
	creds  credentials.TransportCredentials // transport creds
}

func NewServiceClient(wg *sync.WaitGroup, addr string, creds credentials.TransportCredentials, logger *logrus.Logger, ctx context.Context) *ServiceClient {
	return &ServiceClient{
		addr:   addr,
		logger: logger,
		ctx:    ctx,
		wg:     wg,
		creds:  creds,
	}
}
