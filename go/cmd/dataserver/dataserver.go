// The geomserver is TODO application
package main

import (
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
	var grpcPort int
	// flags
	httpPortPtr := flag.Int("httpport", 8123, "The port for the HTTP web server .") // skip validation
	flag.IntVar(&grpcPort, "port", 8120, "The port for the gRPC data application.") // skip validation

	// initialize global pseudo random generator
	rand.Seed(time.Now().Unix())

	// http addresses
	httpAddr := net.JoinHostPort("", strconv.Itoa(*httpPortPtr))
	for {
		fmt.Printf("Hello! This program was compiled on '%s' and has a port '%s'.\n", commit, httpAddr)
		time.Sleep(30 * time.Second)
	}

}
