// The geomserver is a server application for geometry
package main

import (
	"context"
	"flag"
	"math/rand"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/muzammilar/geomrpc/internal/geomserver"
	"github.com/muzammilar/geomrpc/internal/httpserver"
	"github.com/muzammilar/geomrpc/internal/sighandler"
	"github.com/muzammilar/geomrpc/pkg/logs"
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
	httpPortPtr := flag.Int("httpport", 8123, "The port for the HTTP web server .")                         // skip validation
	certFilePtr := flag.String("certfile", "/geometry/certs/server.grpc.crt", "The path of the cert file.") // skip path validation
	keyFilePtr := flag.String("keyfile", "/geometry/certs/server.grpc.key", "The path of the key file.")    // skip path validation
	logPathPtr := flag.String("logpath", "/var/log/geomserver.log", "The path of the logs file.")           // skip path validation
	logLevelPtr := flag.String("loglevel", "info", "The logging level for logrus.Logger.")                  // skip path validation
	logStdOutPtr := flag.Bool("logstdout", false, "The logging level for logrus.Logger.")                   // skip path validation
	flag.IntVar(&grpcPort, "port", 8120, "The port for the gRPC data application.")                         // skip validation

	//parse flags
	flag.Parse()

	// post parsing
	certFile := *certFilePtr
	keyFile := *keyFilePtr
	logPath := *logPathPtr
	logLevel := *logLevelPtr
	logStdOut := *logStdOutPtr

	// http addresses
	httpAddr := net.JoinHostPort("", strconv.Itoa(*httpPortPtr))

	// context - the cancel function is called by the sighandler
	ctx, cancel := context.WithCancel(context.Background())

	// setup logger
	c := logs.NewConfiguration("", logLevel, logPath, logStdOut)
	logger, err := logs.InitLoggerWithFileOutput(c)
	if err != nil {
		panic(err)
	}

	// setup a wait group
	var wg *sync.WaitGroup = new(sync.WaitGroup)

	// start the geometry server
	wg.Add(1)
	go geomserver.Serve(wg, ctx, grpcPort, certFile, keyFile, version, logger)

	// start the http server
	wg.Add(1)
	go httpserver.Serve(wg, httpAddr, ctx, logger)

	//SignalHandler (blocking operation)
	sighandler.SignalHandler(cancel, logger)

	// Wait for all services to cleanly shutdown
	wg.Wait()

}
