package pluto

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/aukbit/fibonacci"
	"github.com/aukbit/pluto/client"
	"github.com/aukbit/pluto/common"
	"github.com/aukbit/pluto/datastore"
	"github.com/aukbit/pluto/server"
	"github.com/aukbit/pluto/server/router"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

const (
	defaultName       = "pluto"
	defaultVersion    = "1.0.0"
	defaultHealthAddr = ":9090"
)

var (
	ErrDatastoreNotInitialized = errors.New("datastore not initialized")
)

// Key ...
type Key string

// Service ...
type Service struct {
	cfg    Config
	close  chan bool
	wg     *sync.WaitGroup
	health *health.Server
	logger *zap.Logger
}

// New returns a new pluto service with Options passed in
func New(opts ...Option) *Service {
	return newService(opts...)
}

func newService(opts ...Option) *Service {
	s := &Service{
		cfg:    newConfig(),
		close:  make(chan bool),
		wg:     &sync.WaitGroup{},
		health: health.NewServer(),
	}
	s.logger, _ = zap.NewProduction()
	if len(opts) > 0 {
		s = s.WithOptions(opts...)
	}
	return s
}

// WithOptions clones the current Service, applies the supplied Options, and
// returns the resulting Service. It's safe to use concurrently.
func (s *Service) WithOptions(opts ...Option) *Service {
	c := s.clone()
	for _, opt := range opts {
		opt.apply(c)
	}
	return c
}

// clone creates a shallow copy service
func (s *Service) clone() *Service {
	copy := *s
	return &copy
}

// Run starts service
func (s *Service) Run() error {
	// set logger
	s.logger = s.logger.With(
		zap.String("id", s.cfg.ID),
		zap.String("name", s.cfg.Name),
	)
	// set health server
	s.setHealthServer()
	// start service
	if err := s.start(); err != nil {
		return err
	}
	// hook run after start
	s.hookAfterStart()
	// wait for all go routines to finish
	s.wg.Wait()
	s.logger.Info("exit")
	return nil
}

// Stop stops service
func (s *Service) Stop() {
	s.logger.Info("stop")
	s.close <- true
}

// Config service configration options
func (s *Service) Config() Config {
	return s.cfg
}

// Push allows to start additional options while service is running
func (s *Service) Push(opts ...Option) {
	for _, opt := range opts {
		opt.apply(s)
	}
}

// Server returns a server instance by name if initialized in service
func (s *Service) Server(name string) (srv *server.Server, ok bool) {
	name = common.SafeName(name, server.DefaultName)
	if srv, ok = s.cfg.Servers[name]; !ok {
		return
	}
	return srv, true
}

// Client returns a client instance by name if initialized in service
func (s *Service) Client(name string) (clt *client.Client, ok bool) {
	name = common.SafeName(name, client.DefaultName)
	if clt, ok = s.cfg.Clients[name]; !ok {
		return
	}
	return clt, true
}

// Datastore returns the datastore instance in initialize in service
func (s *Service) Datastore() (*datastore.Datastore, error) {
	if s.cfg.Datastore != nil {
		return s.cfg.Datastore, nil
	}
	return nil, ErrDatastoreNotInitialized
}

// Health ...
func (s *Service) Health() *healthpb.HealthCheckResponse {
	hcr, err := s.health.Check(
		context.Background(), &healthpb.HealthCheckRequest{Service: s.cfg.ID})
	if err != nil {
		s.logger.Error("Health", zap.String("err", err.Error()))
	}
	return hcr
}

func (s *Service) setHealthServer() {
	s.health.SetServingStatus(s.cfg.ID, 1)
	// Define Router
	mux := router.New()
	mux.GET("/_health/:module/:name", healthHandler)
	// Define server
	srv := server.New(
		server.Name(s.cfg.Name+"_health"),
		server.Addr(s.cfg.HealthAddr),
		server.Mux(mux),
		server.Logger(s.logger),
	)
	s.cfg.Servers[srv.Config().Name] = srv
}

func (s *Service) start() error {
	s.logger.Info("start",
		zap.String("ip4", common.IPaddress()),
		zap.Int("servers", len(s.cfg.Servers)),
		zap.Int("clients", len(s.cfg.Clients)))

	// connect to db
	err := s.initDatastore()
	if err != nil {
		return err
	}
	// run servers
	s.startServers()
	// dial clients
	s.startClients()
	// add go routine to WaitGroup
	s.wg.Add(1)
	go s.waitUntilStopOrSig()
	return nil
}

func (s *Service) hookAfterStart() {
	hooks, ok := s.cfg.Hooks["after_start"]
	if !ok {
		return
	}
	ctx := context.Background()
	for _, h := range hooks {
		h(ctx)
	}
}

func (s *Service) initDatastore() error {
	db, err := s.Datastore()
	if err == ErrDatastoreNotInitialized {
		return nil
	}
	err = db.Init(
		datastore.Logger(s.logger),
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) startServers() {
	for _, srv := range s.cfg.Servers {
		// add go routine to WaitGroup
		s.wg.Add(1)
		go func(srv *server.Server) {
			defer s.wg.Done()
			f := fibonacci.F()
			for {
				err := srv.Run(
					server.Middlewares(
						datastoreContextMiddleware(s),
						serviceContextMiddleware(s),
					),
					server.UnaryServerInterceptors(
						datastoreContextUnaryServerInterceptor(s),
						serviceContextUnaryServerInterceptor(s),
					),
					server.Logger(s.logger),
				)
				if err == nil {
					return
				}
				s.logger.Error(fmt.Sprintf("run failed on server: %v - error: %v", srv.Config().Name, err.Error()))
				time.Sleep(time.Duration(f()) * time.Second)
			}
		}(srv)
	}
}

// startClients listen to the clientsCh
func (s *Service) startClients() {
	go func() {
		for {
			select {
			case clt, ok := <-s.cfg.clientsCh:
				if !ok {
					break
				}
				s.startClient(clt)
			default:
				time.Sleep(500 * time.Millisecond)
				continue
			}
		}
	}()
}

func (s *Service) startClient(clt *client.Client) {
	go func(clt *client.Client) {
		clt.Init()
	}(clt)
}

// waitUntilStopOrSig waits for close channel or syscall Signal
func (s *Service) waitUntilStopOrSig() {
	defer s.wg.Done()
	//  Stop also in case of any host signal
	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, syscall.SIGTERM, syscall.SIGINT)

outer:
	for {
		select {
		case <-s.close:
			// Waits for call to stop
			s.health.SetServingStatus(s.cfg.ID, 2)
			s.unregister()
			s.closeClients()
			s.stopServers()
			break outer
		case sig := <-sigch:
			// Waits for signal to stop
			s.logger.Info("signal received",
				zap.String("signal", sig.String()))
			s.health.SetServingStatus(s.cfg.ID, 2)
			s.unregister()
			s.closeClients()
			s.stopServers()
			break outer
		default:
			// keep on looping, non-blocking channel operations
			time.Sleep(50 * time.Millisecond)
			continue
		}
	}
}

func (s *Service) closeClients() {
	close(s.cfg.clientsCh)
	// for _, clt := range s.cfg.Clients {
	// 	// add go routine to WaitGroup
	// 	s.wg.Add(1)
	// 	go func(clt *client.Client) {
	// 		defer s.wg.Done()
	// 		clt.Close()
	// 	}(clt)
	// }
}

func (s *Service) stopServers() {
	for _, srv := range s.cfg.Servers {
		// add go routine to WaitGroup
		s.wg.Add(1)
		go func(srv *server.Server) {
			defer s.wg.Done()
			srv.Stop()
		}(srv)
	}
}
