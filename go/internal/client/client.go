/*
 * The package implements a basic go client for both geomserver and dataserver
 */

package client

/*
 * Client Package
 */
import (
	"context"
	"crypto/tls"
	"sync"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

/*
 * Public Functions
 */

// Service Client Struct contains all the variables passed to all the clients
type ServiceClient struct {
	addr         string          // server/destination address
	logger       *logrus.Logger  // logger
	ctx          context.Context // client context (used for propagating shutdown signals)
	wg           *sync.WaitGroup // shared wait group (with main routine)
	tlsConf      *tls.Config     // transport credentials/TLS Config
	insecureConn bool            // whether to use insecure connection (i.e. skip TLS)
}

func NewServiceClient(wg *sync.WaitGroup, addr string, tlsConf *tls.Config, insecureConn bool, logger *logrus.Logger, ctx context.Context) *ServiceClient {
	return &ServiceClient{
		addr:         addr,
		logger:       logger,
		ctx:          ctx,
		wg:           wg,
		tlsConf:      tlsConf,
		insecureConn: insecureConn,
	}
}

func (c *ServiceClient) TLSOptions() grpc.DialOption {

	// grpc dial options (with or without TLS)
	if c.insecureConn { // if tls is diabled
		return grpc.WithTransportCredentials(insecure.NewCredentials())
	}
	// dial option for grpc - tls
	return grpc.WithTransportCredentials(credentials.NewTLS(c.tlsConf))
}
