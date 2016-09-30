package pluto_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"bitbucket.org/aukbit/pluto"
	"bitbucket.org/aukbit/pluto/reply"
	"bitbucket.org/aukbit/pluto/server"
	"bitbucket.org/aukbit/pluto/server/router"
	"github.com/paulormart/assert"
)

func Index(w http.ResponseWriter, r *http.Request) {
	reply.Json(w, r, http.StatusOK, "Hello World")
}

const URL = "http://localhost:8083"

func TestService(t *testing.T) {
	// 1. Config service
	s := pluto.NewService(
		pluto.Name("gopher"),
		pluto.Description("gopher super service"),
	)

	//assert.Equal(t, reflect.TypeOf(service.DefaultServer), reflect.TypeOf(s))
	cfg := s.Config()
	assert.Equal(t, true, len(cfg.ID) > 0)
	assert.Equal(t, "pluto_gopher", cfg.Name)
	assert.Equal(t, "gopher super service", cfg.Description)

	// 2. Set http server handlers
	mux := router.NewRouter()
	mux.GET("/", Index)
	// 3. Define server Router
	httpSrv := server.NewServer(server.Addr(":8083"), server.Mux(mux))

	// 4. Init service
	s.Init(pluto.Servers(httpSrv))

	// 5. Run service
	go func() {
		if err := s.Run(); err != nil {
			t.Fatal(err)
		}
	}()
	defer s.Stop()

	// Test
	r, err := http.Get(URL)
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

	// Stop service
	// syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
	// t.Logf("format %v", syscall.Getpid())
}
