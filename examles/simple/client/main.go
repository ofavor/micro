package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ofavor/micro-lite/examles/simple/server/toto"

	"github.com/ofavor/micro-lite"
	"github.com/ofavor/micro-lite/examles/simple/server/foo"
)

func main() {
	service := micro.NewService(
		micro.LogLevel("debug"),
		micro.Name("simple.client"),
		// micro.Version("latest"),
		micro.Address(":8889"),
	)
	go func() {
		for i := 4; i < 6; i++ {
			time.Sleep(2 * time.Second)
			f := foo.NewFooService(service.Client())
			req := &foo.Request{
				Name: "Bob ",
				Age:  uint32(i) + 1,
			}
			rsp, err := f.Bar(context.Background(), req)
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println("Response:", rsp)
			}
		}
	}()
	go func() {
		t := toto.NewTotoService(service.Client())
		req := &toto.Request{
			Val1: 30,
			Val2: 40,
		}
		rsp, err := t.Multiply(context.Background(), req)

		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println(req.Val1, " x ", req.Val2, " = ", rsp.Result)
		}
	}()
	if err := service.Run(); err != nil {
		fmt.Println("Service running with error: ", err)
	}
}
