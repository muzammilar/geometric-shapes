// The geomserver is TODO application
package main

import (
	"context"
	"flag"
	"sync"

	"github.com/muzammilar/geomrpc/internal/client"
	"github.com/muzammilar/geomrpc/internal/tlsconf"
	"github.com/muzammilar/geomrpc/pkg/logs"
)

var (
	version string
	commit  string
)

func main() {

	// flags
	geomAddrPtr := flag.String("geomserver", "geomserver.geometry:8120", "The gRPC endpoint for connecting to the geometry server.") // skip validation
	dataAddrPtr := flag.String("dataserver", "dataserver.geometry:8120", "The gRPC endpoint for connecting to the geometry server.") // skip validation
	certFilePtr := flag.String("certfile", "/geometry/certs/server.grpc.crt", "The path of the cert file.")                          // skip path validation
	insecureConnPtr := flag.Bool("insecure", false, "Use insecure gRPC connection instead of TLS connection.")                       // skip validation
	logPathPtr := flag.String("logpath", "/var/log/goclient.log", "The path of the logs file.")                                      // skip path validation
	logLevelPtr := flag.String("loglevel", "info", "The logging level for logrus.Logger.")                                           // skip validation
	logStdOutPtr := flag.Bool("logstdout", false, "The logging level for logrus.Logger.")                                            // skip validation

	//parse flags
	flag.Parse()

	// post parsing
	geomAddr := *geomAddrPtr
	dataAddr := *dataAddrPtr
	certFile := *certFilePtr
	insecureConn := *insecureConnPtr
	logPath := *logPathPtr
	logLevel := *logLevelPtr
	logStdOut := *logStdOutPtr

	// context - the cancel function is called by the sighandler
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// setup logger
	c := logs.NewConfiguration("", logLevel, logPath, logStdOut)
	logger, err := logs.InitLoggerWithFileOutput(c)
	if err != nil {
		panic(err)
	}

	// setup a wait group
	var wg *sync.WaitGroup = new(sync.WaitGroup)

	// setup the certificate (if needed)
	tlsconf, err := tlsconf.ClientTLSConfigWithCustomRootCA(certFile, logger)
	if err != nil {
		panic(err)
	}

	// clients
	geomClient := client.NewServiceClient(wg, geomAddr, tlsconf, insecureConn, logger, ctx)
	dataClient := client.NewServiceClient(wg, dataAddr, tlsconf, insecureConn, logger, ctx)

	// start the shapecalc - geometry client
	wg.Add(1)
	go client.GeometryClient(geomClient)

	// start shapecalc - info client
	wg.Add(1)
	go client.InfoClient(geomClient)

	// start the shapestore - generator client
	wg.Add(1)
	go client.GeneratorClient(dataClient)

	// start shapestore - store client
	wg.Add(1)
	go client.StoreClient(dataClient)

	//SignalHandler (blocking operation)
	//sighandler.SignalHandler(cancel, logger)

	// Wait for all services to cleanly shutdown
	wg.Wait()

}
