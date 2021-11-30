/*
	The serverstats package implements the grpc/stats.Handler interface and the related read-only functions
*/
package serverstats

import (
	"context"
	"strconv"

	"github.com/muzammilar/geomrpc/pkg/grpcmetrics"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/stats"
)

const method = "method"

// NOTE: This package has not been profiled or benchmarked for performance. The metrics collection can significantly decrease the application performance

// NOTE: Since the stats handler GRPCStats is a read only struct, pointer recievers for functions of the interface are NOT allowed.
// GRPCStats implements the grpc/stats.Handler interface: https://pkg.go.dev/google.golang.org/grpc@v1.42.0/stats#Handler
type GRPCStats struct {
	logger *logrus.Logger
}

// Create a new GRPCStats struct. Note that these fields can not be updated/modified in the stats.Handler interface functions
func NewGRPCStats(l *logrus.Logger) GRPCStats {
	return GRPCStats{
		logger: l,
	}
}

// TagRPC can attach some information to the given context.
// The context used for the rest lifetime of the RPC will be derived from
// the returned context.
func (g GRPCStats) TagRPC(ctx context.Context, rtag *stats.RPCTagInfo) context.Context {
	// Muz: This would be interesting to benchmark for performance
	c := context.WithValue(ctx, method, rtag.FullMethodName)
	return c
}

// HandleRPC processes the RPC stats.
func (g GRPCStats) HandleRPC(ctx context.Context, rpcStats stats.RPCStats) {
	// gRPC Method Name
	methodName := ctx.Value(method).(string)
	switch stat := rpcStats.(type) {
	case *stats.InPayload:
		grpcmetrics.PayloadBytes.WithLabelValues(methodName, "in", "compressed", strconv.FormatBool(stat.Client)).Observe(float64(stat.WireLength))
		grpcmetrics.PayloadBytes.WithLabelValues(methodName, "in", "uncompressed", strconv.FormatBool(stat.Client)).Observe(float64(stat.Length))
	case *stats.OutPayload:
		grpcmetrics.PayloadBytes.WithLabelValues(methodName, "out", "compressed", strconv.FormatBool(stat.Client)).Observe(float64(stat.WireLength))
		grpcmetrics.PayloadBytes.WithLabelValues(methodName, "out", "uncompressed", strconv.FormatBool(stat.Client)).Observe(float64(stat.Length))
	case *stats.InHeader:
		grpcmetrics.MethodCalls.WithLabelValues(stat.FullMethod).Inc()
		grpcmetrics.HeaderBytes.WithLabelValues(stat.FullMethod, "in", strconv.FormatBool(stat.Client)).Observe(float64(stat.WireLength))
	case *stats.OutTrailer:
		grpcmetrics.TrailerBytes.WithLabelValues(method, "out", strconv.FormatBool(stat.Client)).Observe(float64(stat.WireLength))
	default:
		// nothing
	}
}

// TagConn can attach some information to the given context.
// The returned context will be used for stats handling.
// For conn stats handling, the context used in HandleConn for this
// connection will be derived from the context returned.
// For RPC stats handling,
//  - On server side, the context used in HandleRPC for all RPCs on this
// connection will be derived from the context returned.
//  - On client side, the context is not derived from the context returned.
func (g GRPCStats) TagConn(ctx context.Context, ctag *stats.ConnTagInfo) context.Context {
	return ctx
}

// HandleConn processes the Conn stats.
func (g GRPCStats) HandleConn(ctx context.Context, connStats stats.ConnStats) {
	// gRPC Method Name
	switch stat := connStats.(type) {
	case *stats.ConnBegin:
		grpcmetrics.ConnectionsActive.WithLabelValues(strconv.FormatBool(stat.Client)).Inc()
		grpcmetrics.ConnectionsTotal.WithLabelValues(strconv.FormatBool(stat.Client)).Inc()
	case *stats.ConnEnd:
		grpcmetrics.ConnectionsActive.WithLabelValues(strconv.FormatBool(stat.Client)).Dec()
	default:
		// nothing
	}

}
