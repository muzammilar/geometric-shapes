/*
	Package grpcmetrics contains the prometheus handlers for grpc related metrics
*/

package grpcmetrics

import "github.com/prometheus/client_golang/prometheus"

/*
 * Public - Module Variables storing the metrics
 */

var (
	// RPC Stats
	PayloadBytes *prometheus.SummaryVec
	HeaderBytes  *prometheus.SummaryVec
	TrailerBytes *prometheus.SummaryVec

	// Conn Stats
	// Conn Stats only store the information about the underlying (HTTP/2) connections (not the number of gRPC clients connected)
	// A single connection can be used by multiple clients of multiple services, so method information is not part of connection
	ConnectionsActive *prometheus.GaugeVec
	ConnectionsTotal  *prometheus.CounterVec
)

/*
 * Init function
 */

// init is the packages init function that is called by default (the users don't need to call it)
func init() {
	namespace := "grpc"
	// Initialize the metrics
	initializeMetrics(namespace)
	// Register metrics
	registerMetrics()
}

// Initialize all supported the metrics
func initializeMetrics(namespace string) {
	// Create the required metrics
	PayloadBytes = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  namespace,
			Name:       "payload_bytes",
			Help:       "The amount of bytes recieved/sent by the node as part of the payload",
			Objectives: map[float64]float64{0.0: 0.1, 0.25: 0.1, 0.5: 0.1, 0.75: 0.1, 0.9: 0.01, 0.99: 0.001, 1.0: 0.0001},
		},
		[]string{"method", "direction", "compression", "client"},
		// method is method name
		// direction can be in our out, compression is generally 'compressed' or 'uncompressed'
	)
	HeaderBytes = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  namespace,
			Name:       "header_bytes",
			Help:       "The amount of bytes recieved/sent by the node as part of the message/stream header",
			Objectives: map[float64]float64{0.0: 0.1, 0.25: 0.1, 0.5: 0.1, 0.75: 0.1, 0.9: 0.01, 0.99: 0.001, 1.0: 0.0001},
		},
		[]string{"method", "direction", "client"},
		// method is method name
		// direction can be in our out, compression is generally 'compressed' or 'uncompressed'
	)
	TrailerBytes = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  namespace,
			Name:       "trailer_bytes",
			Help:       "The amount of bytes recieved/sent by the node as part of the message/stream trailer",
			Objectives: map[float64]float64{0.0: 0.1, 0.25: 0.1, 0.5: 0.1, 0.75: 0.1, 0.9: 0.01, 0.99: 0.001, 1.0: 0.0001},
		},
		[]string{"method", "direction", "client"},
		// method is method name
		// direction can be in our out, compression is generally 'compressed' or 'uncompressed'
	)

	// Gauges
	ConnectionsActive = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "connections_active",
			Help:      "The current number of active connections for a method.",
		},
		[]string{"client"},
	)

	// Vectors
	ConnectionsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "connections_total",
			Help:      "The total number of connections for a method.",
		},
		[]string{"client"},
	)

}

// Register the relavant metrics
func registerMetrics() {
	prometheus.MustRegister(PayloadBytes)
	prometheus.MustRegister(HeaderBytes)
	prometheus.MustRegister(TrailerBytes)

	prometheus.MustRegister(ConnectionsActive)
	prometheus.MustRegister(ConnectionsTotal)

	// register build info collector
	prometheus.MustRegister(prometheus.NewBuildInfoCollector())
}
