// The geomserver is TODO application
package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"time"
)

var (
	version string
	commit  string
)

func main() {
	// initialize global pseudo random generator
	rand.Seed(time.Now().Unix())

	// flags
	var grpcPort int
	httpPortPtr := flag.Int("httpport", 8123, "The port for the HTTP web server .") // skip validation
	flag.IntVar(&grpcPort, "port", 8120, "The port for the gRPC data application.") // skip validation

	//parse flags
	flag.Parse()

	// http addresses
	httpAddr := net.JoinHostPort("", strconv.Itoa(*httpPortPtr))

	// context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	geomserver.Serve(nil, ctx, grpcPort, certFile string, keyFile string, logger *logrus.Logger)

}
