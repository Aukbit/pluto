package router

import (
	"testing"

	"github.com/paulormart/assert"
)

func TestTransformPath(t *testing.T) {
	var tests = []struct {
		Path   string
		Key    string
		Value  string
		Prefix string
		Params []string
	}{
		{
			Path:   "/",
			Key:    "/",
			Value:  "/",
			Prefix: "",
			Params: []string{},
		},
		{
			Path:   "/home",
			Key:    "/home",
			Value:  "/home",
			Prefix: "",
			Params: []string{},
		},
		{
			Path:   "/home/something",
			Key:    "/home/something",
			Value:  "/something",
			Prefix: "/home",
			Params: []string{},
		},
		{
			Path:   "/home/something/great",
			Key:    "/home/something/great",
			Value:  "/great",
			Prefix: "/home/something",
			Params: []string{},
		},
		{
			Path:   "/:a",
			Key:    "/:",
			Value:  "/:",
			Prefix: "",
			Params: []string{"a"},
		},
		{
			Path:   "/:a/:b",
			Key:    "/:/:",
			Value:  "/:",
			Prefix: "/:",
			Params: []string{"a", "b"},
		},
		{
			Path:   "/home/:a/room",
			Key:    "/home/:/room",
			Value:  "/room",
			Prefix: "/home/:",
			Params: []string{"a"},
		},
		{
			Path:   "/home/:a/room/:b",
			Key:    "/home/:/room/:",
			Value:  "/:",
			Prefix: "/home/:/room",
			Params: []string{"a", "b"},
		},
	}
	for _, test := range tests {
		Key, Value, Prefix, Params := transformPath(test.Path)
		assert.Equal(t, test.Key, Key)
		assert.Equal(t, test.Value, Value)
		assert.Equal(t, test.Prefix, Prefix)
		assert.Equal(t, test.Params, Params)
	}
}

func TestValidPaths(t *testing.T) {
	var tests = []struct {
		Path  string
		Paths map[string][]string
	}{
		{
			Path:  "/",
			Paths: map[string][]string{"/": {}},
		},
		{
			Path:  "/a",
			Paths: map[string][]string{"/a": {}, "/:": {"a"}},
		},
		{
			Path:  "/a/",
			Paths: map[string][]string{"/a": {}, "/:": {"a"}},
		},
		{
			Path:  "/a/b",
			Paths: map[string][]string{"/a/b": {}, "/a/:": {"b"}, "/:/b": {"a"}, "/:/:": {"a", "b"}},
		},
	}
	for _, test := range tests {
		out := make(map[string][]string)
		paths := validPaths(test.Path, "", "", "", nil, out)
		assert.Equal(t, test.Paths, paths)
	}
}

func BenchmarkValidPaths(b *testing.B) {
	d := newData()
	d.value = "home"
	path := "/a/b/c"

	// run the Put function b.N times
	for n := 0; n < b.N; n++ {
		validPaths(path, "", "", "", nil, nil)
	}
}

func TestFindData(t *testing.T) {
	r := NewRouter()
	r.GET("/home", DefaultRootHandler)
	var tests = []struct {
		Method string
		Path   string
		Data   *data
		Values []string
	}{
		{
			Method: "GET",
			Path:   "/home",
			Data: &data{
				value:   "/home",
				prefix:  "",
				vars:    []string{},
				methods: map[string]HandlerFunc{"GET": DefaultRootHandler},
			},
			Values: []string{},
		},
	}
	for _, test := range tests {
		data, values := findData(r, test.Method, test.Path, "", "", "", []string{})
		if data != nil {
			assert.Equal(t, test.Data.Value(), data.Value())
			assert.Equal(t, test.Values, values)
		} else {
			assert.Equal(t, test.Data, data)
		}
	}
}
