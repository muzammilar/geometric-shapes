// The geomserver is TODO application
package main

import (
	"context"
	"crypto/x509"
	"flag"
	"io/ioutil"
	"sync"

	"github.com/muzammilar/geometric-shapes/internal/client"
	"github.com/muzammilar/geometric-shapes/internal/sighandler"
	"github.com/muzammilar/geometric-shapes/pkg/logs"
	"google.golang.org/grpc/credentials"
)

var (
	version string
	commit  string
)

func main() {

	// flags
	geomAddrPtr := flag.String("geomserver", "geomserver.grpc:8120", "The gRPC endpoint for connecting to the geometry server.") // skip validation
	dataAddrPtr := flag.String("dataserver", "dataserver.grpc:8120", "The gRPC endpoint for connecting to the geometry server.") // skip validation
	certFilePtr := flag.String("certfile", "/geometry/certs/server.grpc.crt", "The path of the cert file.")                      // skip path validation
	logPathPtr := flag.String("logpath", "/var/log/goclient.log", "The path of the logs file.")                                  // skip path validation
	logLevelPtr := flag.String("loglevel", "info", "The logging level for logrus.Logger.")                                       // skip path validation

	//parse flags
	flag.Parse()

	// post parsing
	geomAddr := *geomAddrPtr
	dataAddr := *dataAddrPtr
	certFile := *certFilePtr
	logPath := *logPathPtr
	logLevel := *logLevelPtr

	// context - the cancel function is called by the sighandler
	ctx, cancel := context.WithCancel(context.Background())

	// setup logger
	c := logs.NewConfiguration("", logLevel, logPath)
	logger, err := logs.InitLoggerWithFileOutput(c)
	if err != nil {
		panic(err)
	}

	// setup a wait group
	var wg *sync.WaitGroup = new(sync.WaitGroup)

	// setup the certificate

	// Read cert file
	serverCert, err := ioutil.ReadFile(certFile)
	if err != nil {
		panic(err)
	}

	// Create CertPool - This is only allowed cos it's a PoC
	rootCerts := x509.NewCertPool()
	rootCerts.AppendCertsFromPEM(serverCert)

	// Create credentials
	creds := credentials.NewClientTLSFromCert(rootCerts, "")

	// clients
	geomClient := client.NewServiceClient(wg, geomAddr, creds, logger, ctx)
	dataClient := client.NewServiceClient(wg, dataAddr, creds, logger, ctx)

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
	sighandler.SignalHandler(cancel, logger)

	// Wait for all services to cleanly shutdown
	wg.Wait()

}
