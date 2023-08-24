package graceful

// IGRPCServer is an interface for shutting down *grpc.Server.
type IGRPCServer interface {
	GracefulStop()
}

// IGRPCClientConn is an interface for shutting down *grpc.ClientConn.
type IGRPCClientConn interface {
	Close() error
}

// GRPCServer returns a ShutdownErrorFunc that gracefully stops GRPC server.
func GRPCServer(server IGRPCServer) ShutdownErrorFunc {
	if server != nil {
		return makeErrorFunc(server.GracefulStop)
	}

	return nil
}

// GRPCClient returns a ShutdownErrorFunc that stops GRPC client.
func GRPCClient(client IGRPCClientConn) ShutdownErrorFunc {
	if client != nil {
		return client.Close
	}

	return nil
}
