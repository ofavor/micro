package main

import (
	"context"
	"fmt"

	"github.com/ofavor/micro-lite"
	"github.com/ofavor/micro-lite/examles/simple/server/foo"
)

type myFoo struct {
}

func (s *myFoo) Bar(ctx context.Context, req *foo.Request, rsp *foo.Response) error {
	rsp.Name = req.Name
	rsp.Age = req.Age
	rsp.Adult = req.Age > 4
	return nil
}

func main() {
	service := micro.NewService(
		micro.LogLevel("debug"),
		micro.Name("simple.server"),
		// micro.Version("latest"),
		// micro.Address(":8888"),
	)
	foo.RegisterFooHandler(service.Server(), &myFoo{})
	if err := service.Run(); err != nil {
		fmt.Println("Service running with error: ", err)
	}
}
