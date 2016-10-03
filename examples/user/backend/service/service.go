package backend

import (
	"flag"

	"bitbucket.org/aukbit/pluto"
	"bitbucket.org/aukbit/pluto/examples/user/backend/views"
	pb "bitbucket.org/aukbit/pluto/examples/user/proto"
	"bitbucket.org/aukbit/pluto/server"
	"google.golang.org/grpc"
)

var db_addr = flag.String("db_addr", "127.0.0.1", "cassandra address")
var grpc_port = flag.String("grpc_port", ":65060", "grpc listening port")

func Run() error {
	flag.Parse()

	// GRPC server
	// Define gRPC server and register
	grpcServer := grpc.NewServer()

	// Define Pluto Service
	s := pluto.NewService(
		pluto.Name("backend"),
		pluto.Description("Backend service is responsible for persist data"),
		pluto.Datastore(*db_addr),
	)
	// Register grpc Server
	pb.RegisterUserServiceServer(grpcServer, &backend.User{Cluster: s.Config().Datastore})

	// Define Pluto Server
	grpcSrv := server.NewServer(server.Addr(*grpc_port), server.GRPCServer(grpcServer))

	// 5. Init service
	// TODO remove init method redundant
	s.Init(pluto.Servers(grpcSrv))
	// 6. Run service
	if err := s.Run(); err != nil {
		return err
	}
	return nil
}
