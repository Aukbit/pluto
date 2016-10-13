package pluto

import (
	"bitbucket.org/aukbit/pluto/client"
	"bitbucket.org/aukbit/pluto/server"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

// Service is the basic interface that defines what to expect from any server.
type Service interface {
	Run() error
	Stop()
	Client(string) (client.Client, bool)
	Server(string) (server.Server, bool)
	Config() *Config
	Health() *healthpb.HealthCheckResponse
}

var (
	defaultName    = "pluto"
	defaultVersion = "v1.0.0"
)

// NewService returns a new service with cfg passed in
func NewService(cfgs ...ConfigFunc) Service {
	return newService(cfgs...)
}
