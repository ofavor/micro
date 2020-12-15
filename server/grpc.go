package server

import (
	"context"
	"errors"
	"net"
	"reflect"
	"strings"

	"github.com/ofavor/micro-lite/internal/log"
	"github.com/ofavor/micro-lite/internal/transport"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

type methodType struct {
	method      reflect.Method
	argType     reflect.Type
	replyType   reflect.Type
	contextType reflect.Type
}

type receiver struct {
	name    string                 // name of receiver
	val     reflect.Value          // receiver of value
	typ     reflect.Type           // type of the receiver
	methods map[string]*methodType // registered methods
}

type grpcServer struct {
	opts     Options
	srv      *grpc.Server
	rcvrMap  map[string]*receiver
	handlers map[string]Handler
}

func newGRPCServer(opts ...Option) Server {
	options := defaultOptions()
	for _, o := range opts {
		o(&options)
	}
	return &grpcServer{
		opts:     options,
		rcvrMap:  make(map[string]*receiver),
		handlers: make(map[string]Handler),
	}
}

func (s *grpcServer) Init(opt Option) {
	opt(&s.opts)
}

func (s *grpcServer) Start() error {
	log.Infof("Trying to listen on TCP: %s", s.opts.Address)
	l, err := net.Listen("tcp", s.opts.Address)
	if err != nil {
		return err
	}

	s.srv = grpc.NewServer()
	transport.RegisterMicroServer(s.srv, s)

	go func() {
		if err := s.srv.Serve(l); err != nil {
			log.Error("gRPC server serve error: ", err)
		}
	}()
	return nil
}

func (s *grpcServer) Stop() error {
	return nil
}

func (s *grpcServer) Handle(h Handler) error {
	t := h.Target()
	rcvr := new(receiver)
	rcvr.val = reflect.ValueOf(t)
	rcvr.typ = reflect.TypeOf(t)
	rcvr.name = h.Name() // reflect.Indirect(rcvr.val).Type().Name()
	rcvr.methods = make(map[string]*methodType)

	log.Debug("Register handler: ", rcvr.name)
	// Install the methods
	for m := 0; m < rcvr.typ.NumMethod(); m++ {
		method := rcvr.typ.Method(m)
		if mt := prepareEndpoint(method); mt != nil {
			log.Debug("Endpoint prepared: ", method.Name)
			rcvr.methods[method.Name] = mt
		}
	}
	s.rcvrMap[rcvr.name] = rcvr

	s.handlers[h.Name()] = h
	return nil
}

var (
	// Precompute the reflect type for error. Can't use error directly
	// because Typeof takes an empty interface value. This is annoying.
	typeOfError = reflect.TypeOf((*error)(nil)).Elem()

	typeOfProtoMessage = reflect.TypeOf((*proto.Message)(nil)).Elem()
)

// prepareEndpoint() returns a methodType for the provided method or nil
// in case if the method was unsuitable.
func prepareEndpoint(method reflect.Method) *methodType {
	mtype := method.Type
	mname := method.Name
	var replyType, argType, contextType reflect.Type

	// Endpoint() must be exported.
	if method.PkgPath != "" {
		return nil
	}

	switch mtype.NumIn() {
	case 4:
		// method that takes a context
		argType = mtype.In(2)
		replyType = mtype.In(3)
		contextType = mtype.In(1)
	default:
		log.Errorf("method %v of %v has wrong number of ins: %v", mname, mtype, mtype.NumIn())
		return nil
	}

	// First arg need not be a pointer.
	if !argType.Implements(typeOfProtoMessage) {
		log.Errorf("%v argument type not proto.Message: %v", mname, argType)
		return nil
	}

	// Reply type must be exported.
	if !replyType.Implements(typeOfProtoMessage) {
		log.Errorf("method %v reply type not proto.Message: %v", mname, replyType)
		return nil
	}

	// Endpoint() needs one out.
	if mtype.NumOut() != 1 {
		log.Errorf("method %v has wrong number of outs: %v", mname, mtype.NumOut())
		return nil
	}
	// The return type of the method must be error.
	if returnType := mtype.Out(0); returnType != typeOfError {
		log.Errorf("method %v returns %v not error", mname, returnType.String())
		return nil
	}
	return &methodType{method: method, argType: argType, replyType: replyType, contextType: contextType}
}

func (s *grpcServer) HandleRequest(ctx context.Context, req *transport.Request) (*transport.Response, error) {
	log.Debug("Handling request: ", req.Endpoint)
	arr := strings.Split(req.Endpoint, ".")
	if len(arr) != 2 {
		log.Error("Bad request endpoint: ", req.Endpoint)
		return nil, errors.New("bad request endpoint")
	}

	rcvr, ok := s.rcvrMap[arr[0]]
	if !ok {
		log.Error("No handler found for endpoint: ", req.Endpoint)
		return nil, errors.New("no handler found")
	}
	method, ok := rcvr.methods[arr[1]]
	if !ok {
		log.Error("No method found for endpoint: ", req.Endpoint)
		return nil, errors.New("no method found")
	}
	argv := reflect.New(method.argType.Elem())
	err := proto.Unmarshal(req.Data, argv.Interface().(proto.Message))
	if err != nil {
		log.Error("Unmarshal request data error: ", err)
		return nil, err
	}
	replyv := reflect.New(method.replyType.Elem())
	function := method.method.Func
	returns := function.Call([]reflect.Value{rcvr.val, method.prepareContext(ctx), reflect.ValueOf(argv.Interface()), reflect.ValueOf(replyv.Interface())})
	if rerr := returns[0].Interface(); rerr != nil {
		return nil, rerr.(error)
	}
	data, err := proto.Marshal(replyv.Interface().(proto.Message))
	if err != nil {
		return nil, err
	}
	return &transport.Response{Data: data}, nil
}

func (m *methodType) prepareContext(ctx context.Context) reflect.Value {
	if contextv := reflect.ValueOf(ctx); contextv.IsValid() {
		return contextv
	}
	return reflect.Zero(m.contextType)
}
