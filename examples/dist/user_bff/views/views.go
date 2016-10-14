package views

import (
	"errors"
	"net/http"

	"bitbucket.org/aukbit/pluto"
	pb "bitbucket.org/aukbit/pluto/examples/dist/user_bff/proto"
	"bitbucket.org/aukbit/pluto/reply"
	"github.com/golang/protobuf/jsonpb"
	"github.com/uber-go/zap"
)

var (
	errClientUserNotAvailable = errors.New("Client user not available")
)

func PostHandler(w http.ResponseWriter, r *http.Request) {
	// get context
	ctx := r.Context()
	// get logger from context
	log := ctx.Value("logger").(zap.Logger)
	// new user
	newUser := &pb.NewUser{}
	if err := jsonpb.Unmarshal(r.Body, newUser); err != nil {
		log.Error(err.Error())
		reply.Json(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	// get gRPC client from service
	c, ok := ctx.Value("pluto").(pluto.Service).Client("client")
	if !ok {
		log.Error(errClientUserNotAvailable.Error())
		reply.Json(w, r, http.StatusInternalServerError, errClientUserNotAvailable)
		return
	}
	// make a call the backend service
	user, err := c.Call().(pb.UserServiceClient).CreateUser(ctx, newUser)
	if err != nil {
		log.Error(err.Error())
		reply.Json(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	reply.Json(w, r, http.StatusCreated, user)
}

func GetHandlerDetail(w http.ResponseWriter, r *http.Request) {
	// get context
	ctx := r.Context()
	// get logger from context
	log := ctx.Value("logger").(zap.Logger)
	// get id context
	id := ctx.Value("id").(string)
	// set proto user
	user := &pb.User{Id: id}
	// get gRPC client from service
	c, ok := ctx.Value("pluto").(pluto.Service).Client("client")
	if !ok {
		log.Error(errClientUserNotAvailable.Error())
		reply.Json(w, r, http.StatusInternalServerError, errClientUserNotAvailable)
		return
	}
	// make a call the backend service
	user, err := c.Call().(pb.UserServiceClient).ReadUser(ctx, user)
	if err != nil {
		log.Error(err.Error())
		reply.Json(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	reply.Json(w, r, http.StatusOK, user)
}

func PutHandler(w http.ResponseWriter, r *http.Request) {
	// get context
	ctx := r.Context()
	// get logger from context
	log := ctx.Value("logger").(zap.Logger)
	// get id context
	id := ctx.Value("id").(string)
	// set proto user
	user := &pb.User{Id: id}
	// unmarshal body
	if err := jsonpb.Unmarshal(r.Body, user); err != nil {
		log.Error(err.Error())
		reply.Json(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	// get gRPC client from service
	c, ok := ctx.Value("pluto").(pluto.Service).Client("client")
	if !ok {
		log.Error(errClientUserNotAvailable.Error())
		reply.Json(w, r, http.StatusInternalServerError, errClientUserNotAvailable.Error())
		return
	}
	// make a call the backend service
	user, err := c.Call().(pb.UserServiceClient).UpdateUser(ctx, user)
	if err != nil {
		log.Error(err.Error())
		reply.Json(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	reply.Json(w, r, http.StatusOK, user)
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	// get context
	ctx := r.Context()
	// get logger from context
	log := ctx.Value("logger").(zap.Logger)
	// get id context
	id := ctx.Value("id").(string)
	// set proto user
	user := &pb.User{Id: id}
	// get gRPC client from service
	c, ok := ctx.Value("pluto").(pluto.Service).Client("client")
	if !ok {
		log.Error(errClientUserNotAvailable.Error())
		reply.Json(w, r, http.StatusInternalServerError, errClientUserNotAvailable)
		return
	}
	// make a call the backend service
	user, err := c.Call().(pb.UserServiceClient).DeleteUser(ctx, user)
	if err != nil {
		log.Error(err.Error())
		reply.Json(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	reply.Json(w, r, http.StatusOK, user)
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	// get context
	ctx := r.Context()
	// get logger from context
	log := ctx.Value("logger").(zap.Logger)
	// get parameters
	n := r.URL.Query().Get("name")
	// set proto filter
	filter := &pb.Filter{Name: n}
	// get gRPC client from service
	c, ok := ctx.Value("pluto").(pluto.Service).Client("client")
	if !ok {
		log.Error(errClientUserNotAvailable.Error())
		reply.Json(w, r, http.StatusInternalServerError, errClientUserNotAvailable.Error())
		return
	}
	// make a call the backend service
	users, err := c.Call().(pb.UserServiceClient).FilterUsers(ctx, filter)
	if err != nil {
		log.Error(err.Error())
		reply.Json(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	reply.Json(w, r, http.StatusOK, users)
}