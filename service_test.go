package pluto

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/golang/protobuf/jsonpb"
	"github.com/paulormart/assert"

	context "golang.org/x/net/context"

	"github.com/aukbit/pluto/v6/auth/jwt"
	"github.com/aukbit/pluto/v6/client"
	"github.com/aukbit/pluto/v6/reply"
	"github.com/aukbit/pluto/v6/server"
	"github.com/aukbit/pluto/v6/server/router"
	pb "github.com/aukbit/pluto/v6/test/proto"
	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

const serviceURL = "http://localhost:8081"
const serviceCallURL = "http://localhost:8081/call"
const serviceCallWithCredentialsURL = "http://localhost:8081/call-with-credentials"
const serviceCallWithCredentialsWrapErrorURL = "http://localhost:8081/call-with-credentials-wrap-error"
const healthURL = "http://localhost:9090/_health"
const healthzURL = "http://localhost:9090/healthz"

var serviceName = "gopher"

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	c, ok := FromContext(ctx).Client(serviceName)
	if !ok {
		panic("grpc client not available")
	}
	conn, err := c.Dial(client.Timeout(2 * time.Second))
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	res, err := c.Stub(conn).(pb.GreeterClient).SayHello(ctx, &pb.HelloRequest{Name: "World"})
	if err != nil {
		panic(err)
	}
	reply.Json(w, r, http.StatusOK, res.GetMessage())
}

func CallHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	in := &pb.HelloRequest{Name: "World"}
	response, err := Call(ctx, serviceName, "SayHello", in)
	if err != nil {
		panic(err)
	}
	reply.Json(w, r, http.StatusOK, response.(*pb.HelloReply).GetMessage())
}

func CallWithCredentialsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	in := &pb.GoodbyeRequest{Name: "Traffic!"}
	response, err := CallWithCredentials(ctx, serviceName, "SayGoodbye", in)
	if err != nil {
		panic(err)
	}
	reply.Json(w, r, http.StatusOK, response.(*pb.GoodbyeReply).GetMessage())
}

func CallWithCredentialsHandlerWrapError(w http.ResponseWriter, r *http.Request) *router.Err {
	ctx := r.Context()
	in := &pb.GoodbyeRequest{Name: "Car"}
	response, err := CallWithCredentials(ctx, serviceName, "SayGoodbye", in)
	if err != nil {
		panic(err)
	}
	reply.Jsonpb(w, r, http.StatusOK, &jsonpb.Marshaler{OrigName: true}, response.(*pb.GoodbyeReply))
	return nil
}

type greeter struct{}

// SayHello implements helloworld.GreeterServer
func (s *greeter) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: fmt.Sprintf("Hello %v", in.Name)}, nil
}

// SayGoodbye implements helloworld.GreeterServer
func (s *greeter) SayGoodbye(ctx context.Context, in *pb.GoodbyeRequest) (*pb.GoodbyeReply, error) {

	r, err := CallWithCredentials(ctx, serviceName, "SayGoodbyeAgain", &pb.GoodbyeRequest{Name: in.Name})
	if err != nil {
		panic(err)
	}
	fmt.Println(r)

	return &pb.GoodbyeReply{Message: fmt.Sprintf("Bye Bye %v", in.Name)}, nil
}

// SayGoodbye implements helloworld.GreeterServer
func (s *greeter) SayGoodbyeAgain(ctx context.Context, in *pb.GoodbyeRequest) (*pb.GoodbyeReply, error) {

	t, ok := jwt.TokenFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("token not available in context")
	}

	fmt.Println(t)

	return &pb.GoodbyeReply{Message: fmt.Sprintf("Bye Bye %v", in.Name)}, nil
}

func TestMain(m *testing.M) {
	// Create db client
	// cfg := gocql.NewCluster("localhost")
	// cfg.ProtoVersion = 3
	// cfg.Keyspace = "default"

	// Define router
	mux := router.New()
	mux.GET("/", IndexHandler)
	mux.GET("/call", CallHandler)
	mux.GET("/call-with-credentials", jwt.WrapBearerToken(CallWithCredentialsHandler))
	mux.Handle("GET", "/call-with-credentials-wrap-error", jwt.WrapBearerTokenErr(router.WrapErr(CallWithCredentialsHandlerWrapError)))
	// Create pluto server
	srvHTTP := server.New(
		server.Name(serviceName+"_http"),
		server.Description("gopher super server"),
		server.Addr(":8081"),
		server.Mux(mux),
		// server.Middlewares(ext.CassandraMiddleware("cassandra", cfg)),
	)
	// Create grpc pluto server
	srvGRPC := server.New(
		server.Name(serviceName+"_grpc"),
		server.Description("grpc super server"),
		server.Addr(":65060"),
		server.GRPCRegister(func(g *grpc.Server) {
			pb.RegisterGreeterServer(g, &greeter{})
		}),
		server.UnaryServerInterceptors(jwt.BearerTokenUnaryServerInterceptor()),
		// server.UnaryServerInterceptors(ext.CassandraUnaryServerInterceptor("cassandra", cfg)),
		// server.StreamServerInterceptors(ext.CassandraStreamServerInterceptor("cassandra", cfg)),
	)
	// Create grpc pluto client
	cltGRPC := client.New(
		client.Name(serviceName),
		client.Target("localhost:65060"),
		// client.TargetName("grpc"),
		client.GRPCRegister(func(cc *grpc.ClientConn) interface{} {
			return pb.NewGreeterClient(cc)
		}),
	)

	// Define service Discovery
	// d := discovery.NewDiscovery(discovery.Addr("192.168.99.100:8500"))
	// Hook functions
	fn1 := func(ctx context.Context) error {
		log.Print("first run after service starts")
		return nil
	}
	fn2 := func(ctx context.Context) error {
		log.Print("second run after service starts")
		return nil
	}
	// Define Pluto Service
	s := New(
		Name(serviceName),
		Servers(srvHTTP),
		Servers(srvGRPC),
		Clients(cltGRPC),
		HookAfterStart(fn1, fn2),
		HealthAddr(":9090"),
	)

	// if !testing.Short() {
	// Run Server
	go func() {
		if err := s.Run(); err != nil {
			log.Fatal(err)
		}
	}()
	time.Sleep(time.Second)
	// }
	result := m.Run()
	if !testing.Short() {
		// Stop Server
		s.Stop()
		time.Sleep(time.Millisecond * 100)
	}
	os.Exit(result)
}

func TestService(t *testing.T) {
	time.Sleep(time.Second)
	r, err := http.Get(serviceURL)
	if err != nil {
		t.Fatal(err)
	}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Fatal(err)
	}
	defer r.Body.Close()

	var message string
	if err := json.Unmarshal(b, &message); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
	assert.Equal(t, http.StatusOK, r.StatusCode)
	assert.Equal(t, "Hello World", message)
}

func TestCallMethod(t *testing.T) {
	time.Sleep(time.Second)
	r, err := http.Get(serviceCallURL)
	if err != nil {
		t.Fatal(err)
	}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Fatal(err)
	}
	defer r.Body.Close()

	var message string
	if err := json.Unmarshal(b, &message); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
	assert.Equal(t, http.StatusOK, r.StatusCode)
	assert.Equal(t, "Hello World", message)
}

func TestCredentials(t *testing.T) {
	time.Sleep(time.Second)
	req, err := http.NewRequest("GET", serviceCallWithCredentialsURL, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer my-token-123")
	c := &http.Client{}
	r, err := c.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Fatal(err)
	}
	defer r.Body.Close()

	var message string
	if err := json.Unmarshal(b, &message); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
	assert.Equal(t, http.StatusOK, r.StatusCode)
	assert.Equal(t, "Bye Bye Traffic!", message)
}

func TestWrapErr(t *testing.T) {
	time.Sleep(time.Second)
	req, err := http.NewRequest("GET", serviceCallWithCredentialsWrapErrorURL, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer my-token-123")
	c := &http.Client{}
	r, err := c.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	reply := &pb.HelloReply{}
	if err := jsonpb.Unmarshal(r.Body, reply); err != nil {
		t.Fatal(err)
	}
	defer r.Body.Close()

	assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
	assert.Equal(t, http.StatusOK, r.StatusCode)
	assert.Equal(t, "Bye Bye Car", reply.GetMessage())
}

func TestLiveHealth(t *testing.T) {
	r, err := http.Get(healthzURL + "/live")
	if err != nil {
		log.Fatal(err)
	}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()
	hcr := &healthpb.HealthCheckResponse{}
	if err := json.Unmarshal(b, hcr); err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
	assert.Equal(t, http.StatusOK, r.StatusCode)
	assert.Equal(t, "SERVING", hcr.Status.String())
	time.Sleep(time.Millisecond * 100)
}

func TestReadyHealth(t *testing.T) {
	r, err := http.Get(healthzURL + "/ready")
	if err != nil {
		log.Fatal(err)
	}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()
	hcr := &healthpb.HealthCheckResponse{}
	if err := json.Unmarshal(b, hcr); err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
	assert.Equal(t, http.StatusOK, r.StatusCode)
	assert.Equal(t, "SERVING", hcr.Status.String())
	time.Sleep(time.Millisecond * 100)
}

func TestHealth(t *testing.T) {

	time.Sleep(time.Second)

	var tests = []struct {
		Path         string
		Body         io.Reader
		BodyContains string
		Status       int
	}{
		{
			Path:         "/server/" + serviceName + "_pluto_health_server",
			BodyContains: `SERVING`,
			Status:       http.StatusOK,
		},
		{
			Path:         "/server/" + serviceName + "_http_server",
			BodyContains: `SERVING`,
			Status:       http.StatusOK,
		},
		{
			Path:         "/server/" + serviceName + "_grpc_server",
			BodyContains: `SERVING`,
			Status:       http.StatusOK,
		},
		{
			Path:         "/server/something_wrong",
			BodyContains: `UNKNOWN`,
			Status:       http.StatusNotFound,
		},
		{
			Path:         "/client/" + serviceName + "_client",
			BodyContains: `SERVING`,
			Status:       http.StatusOK,
			// if discovery is active when testing this handler is not serving
			// BodyContains: `NOT_SERVING`,
			// Status:       http.StatusTooManyRequests,
		},
		{
			Path:         "/client/something_wrong",
			BodyContains: `UNKNOWN`,
			Status:       http.StatusNotFound,
		},
		{
			Path:         "/pluto/gopher_pluto",
			BodyContains: `SERVING`,
			Status:       http.StatusOK,
		},
		{
			Path:         "/pluto/something_wrong",
			BodyContains: `UNKNOWN`,
			Status:       http.StatusNotFound,
		},
	}
	for _, test := range tests {

		r, err := http.Get(healthURL + test.Path)
		if err != nil {
			log.Fatal(err)
		}
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		defer r.Body.Close()
		hcr := &healthpb.HealthCheckResponse{}
		if err := json.Unmarshal(b, hcr); err != nil {
			log.Fatal(err)
		}
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, test.Status, r.StatusCode)
		assert.Equal(t, test.BodyContains, hcr.Status.String())
		time.Sleep(time.Millisecond * 100)
	}
}
