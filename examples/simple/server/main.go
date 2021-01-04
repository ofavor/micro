package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/ofavor/micro-lite/examples/simple/server/toto"
	"github.com/ofavor/micro-lite/internal/transport"
	"github.com/ofavor/micro-lite/server"

	"github.com/ofavor/micro-lite/examples/simple/server/foo"

	"github.com/ofavor/micro-lite/examples/simple/server/srv"

	"github.com/ofavor/micro-lite"
)

type myFoo struct {
}

func (s *myFoo) Bar(ctx context.Context, req *foo.Request, rsp *foo.Response) error {
	fmt.Println("myFoo.Bar is handling", req.Name, req.Age)
	rsp.Name = req.Name
	rsp.Age = req.Age
	rsp.Adult = req.Age > 4
	return nil
}

func logWrapper(f server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req *transport.Request, rsp interface{}) error {
		fmt.Println("This is a logger wrapper:", req.Endpoint)
		fmt.Println(">>>>>Before call")
		f(ctx, req, rsp)
		fmt.Println(">>>>>After call")
		return nil
	}
}

func main() {
	fmt.Println(">>>>>>>>>>>>>", os.Args)
	regAddrs := flag.String("registry_addrs", "127.0.0.1:2379", "registry addresses, splitted by ','")
	flag.Parse()

	fmt.Println(">>>>>>>>>>>>>", *regAddrs)
	service := micro.NewService(
		micro.LogLevel("debug"),
		micro.Name("simple.server"),
		micro.ID("srv1"),
		// micro.Version("latest"),
		// micro.Address(":8888"),
		micro.RegistryAddrs(strings.Split(*regAddrs, ",")),
		micro.WrapHandler(logWrapper),
	)
	foo.RegisterFooHandler(service.Server(), &myFoo{})
	toto.RegisterTotoHandler(service.Server(), &srv.TotoHandler{})
	if err := service.Run(); err != nil {
		fmt.Println("Service running with error: ", err)
	}
}
