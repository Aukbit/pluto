package server

import (
	"log"
	"net"
	"syscall"
	"os/signal"
	"os"
	"google.golang.org/grpc"
)

// A Server defines parameters for running an HTTP server.
// The zero value for Server is a valid configuration.
type gRPCServer struct {
	cfg 			*Config
	// close chan for graceful shutdown
	close 			chan bool
}

// newGRPCServer will instantiate a new Server with the given config
func newGRPCServer(cfgs ...ConfigFunc) Server {
	c := newConfig(cfgs...)
	c.Format = "grpc"
	return &gRPCServer{cfg: c, close: make(chan bool)}
}

func (s *gRPCServer) Init(cfgs ...ConfigFunc) error {
	for _, c := range cfgs {
		c(s.cfg)
	}
	return nil
}

func (s *gRPCServer) Config() *Config {
	cfg := s.cfg
	return cfg
}

// Run
func (s *gRPCServer) Run() error {
	if err := s.start(); err != nil {
		return err
	}
	// parse address for host, port
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
	sig := <-ch
	log.Printf("----- %s signal %v received ", s.cfg.Name, sig)
	return s.Stop()
}

// Stop sends message to close the listener via channel
func (s *gRPCServer) Stop() error {
	s.close <-true
	return nil
}

// start start the Server
func (s *gRPCServer) start() error {
	log.Printf("START %s %s \t%s", s.cfg.Format, s.cfg.Name, s.cfg.Id)
	if err := s.listenAndServe(); err != nil{
		log.Fatalf("ERROR %s s.listenAndServe() %v", s.cfg.Name, err)
	}
	return nil
}

func (s *gRPCServer) listenAndServe() (err error) {

	addr := s.cfg.Addr
	if addr == "" {
		addr, err = getNewAddr()
		if err != nil {
			return err
		}
	}
	// set cfg.Addr
	s.cfg.Addr = addr

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	ln = net.Listener(TcpKeepAliveListener{ln.(*net.TCPListener)})

	// new gRPC server
	g := grpc.NewServer()

	// pb.RegisterServerFunc
	s.cfg.RegisterServerFunc(g)

	go func() {
		if err := g.Serve(ln); err != nil {
			log.Fatalf("ERROR %s g.Serve(lis) %v", s.cfg.Name, err)
		}
	}()
	//
	log.Printf("----- %s %s listening on %s", s.cfg.Format, s.cfg.Name, ln.Addr().String())
	//
	go func() {
		// Waits for call to stop
		<-s.close
		log.Printf("CLOSE %s received", s.cfg.Name)
		// close listener
		if err := ln.Close(); err != nil {
			log.Fatalf("ERROR %s ln.Close() %v", s.cfg.Name, err)
		}
		log.Printf("----- %s listener closed", s.cfg.Name)
	}()

	return nil
}