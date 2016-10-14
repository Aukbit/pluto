package backend

import (
	"flag"

	"github.com/uber-go/zap"

	"bitbucket.org/aukbit/pluto"
	pba "bitbucket.org/aukbit/pluto/auth/proto"
	"bitbucket.org/aukbit/pluto/client"
	"bitbucket.org/aukbit/pluto/examples/auth/backend/views"
	pbu "bitbucket.org/aukbit/pluto/examples/user/proto"
	"bitbucket.org/aukbit/pluto/server"
	"google.golang.org/grpc"
)

var (
	userTarget = flag.String("user_target", "127.0.0.1:65080", "user backend address")
	grpcPort   = flag.String("grpc_port", ":65081", "grpc listening port")
	logger     = zap.New(zap.NewJSONEncoder())
)

// Run runs auth backend service
func Run() error {
	flag.Parse()

	// Define user Client
	clt := client.NewClient(
		client.Name("user"),
		client.GRPCRegister(func(cc *grpc.ClientConn) interface{} {
			return pbu.NewUserServiceClient(cc)
		}),
		client.Target(*userTarget))

	// Define Pluto Server
	srv := server.NewServer(
		server.Addr(*grpcPort),
		server.GRPCRegister(func(g *grpc.Server) {
			pba.RegisterAuthServiceServer(g, &backend.AuthViews{})
		}))

	// Define Pluto Service
	s := pluto.NewService(
		pluto.Name("auth_backend"),
		pluto.Description("Backend service issuing access tokens to the client after successfully authenticating the resource owner and obtaining authorization"),
		pluto.Servers(srv),
		pluto.Clients(clt))

	// Run service
	if err := s.Run(); err != nil {
		return err
	}
	return nil
}