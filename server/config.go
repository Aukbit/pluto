package server

import (
	"crypto/tls"
	"fmt"
	"log"
	"regexp"
	"strings"

	"bitbucket.org/aukbit/pluto/server/router"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type Config struct {
	Id          string
	Name        string
	Description string
	Version     string
	Addr        string // TCP address (e.g. localhost:8000) to listen on, ":http" if empty
	Format      string
	Mux         router.Mux
	TLSConfig   *tls.Config // optional TLS config, used by ListenAndServeTLS
	GRPCServer  *grpc.Server
}

type ConfigFunc func(*Config)

var (
	defaultAddr   = ":8080"
	defaultFormat = "http"
)

func newConfig(cfgs ...ConfigFunc) *Config {

	cfg := &Config{Addr: defaultAddr, Format: defaultFormat}

	for _, c := range cfgs {
		c(cfg)
	}

	if len(cfg.Id) == 0 {
		cfg.Id = uuid.New().String()
	}

	if len(cfg.Name) == 0 {
		cfg.Name = DefaultName
	}

	if len(cfg.Version) == 0 {
		cfg.Version = DefaultVersion
	}

	return cfg
}

// Id server id
func Id(id string) ConfigFunc {
	return func(cfg *Config) {
		cfg.Id = id
	}
}

// Name server name
func Name(n string) ConfigFunc {
	return func(cfg *Config) {
		// support only alphanumeric and underscore characters
		reg, err := regexp.Compile("[^A-Za-z0-9_]+")
		if err != nil {
			log.Fatal(err)
		}
		safe := reg.ReplaceAllString(n, "_")
		cfg.Name = fmt.Sprintf("%s_%s", DefaultName, strings.ToLower(safe))
	}
}

// Description server description
func Description(d string) ConfigFunc {
	return func(cfg *Config) {
		cfg.Description = d
	}
}

// Addr server address
func Addr(a string) ConfigFunc {
	return func(cfg *Config) {
		cfg.Addr = a
	}
}

// Mux server multiplexer
func Mux(m router.Mux) ConfigFunc {
	return func(cfg *Config) {
		cfg.Mux = m
	}
}

// TLSConfig server multiplexer
func TLSConfig(certFile, keyFile string) ConfigFunc {
	return func(cfg *Config) {
		cer, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			log.Printf("ERROR tls.LoadX509KeyPair %v", err)
			return
		}
		cfg.TLSConfig = &tls.Config{
			MinVersion:               tls.VersionTLS12,
			CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
			PreferServerCipherSuites: true,
			Certificates:             []tls.Certificate{cer},
		}
		cfg.Format = "https"
	}
}

func GRPCServer(s *grpc.Server) ConfigFunc {
	return func(cfg *Config) {
		cfg.GRPCServer = s
		cfg.Format = "grpc"
	}
}
