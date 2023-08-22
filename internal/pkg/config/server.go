package config

// Server ..
type Server struct {
	GRPCPort  uint32
	HTTPPort  uint32
	DebugPort uint32
}

func serverConfig() Server {
	return Server{
		GRPCPort:  Uint32("service.ports.grpc"),
		HTTPPort:  Uint32("service.ports.http"),
		DebugPort: Uint32("service.ports.debug"),
	}
}
