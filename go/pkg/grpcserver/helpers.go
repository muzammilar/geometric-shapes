/*
	The package grpcserver provides useful set of function to create grpc servers
*/
package grpcserver

import (
	"context"

	"google.golang.org/grpc/peer"
)

// GetRemoteHostFromContext returns the host address of the remote peer or returns `unknown:0`
func GetRemoteHostFromContext(ctx context.Context) string {
	remotePeer, ok := peer.FromContext(ctx)
	if !ok {
		return "unknown:0"
	}
	return remotePeer.Addr.String()
}
