package main

import (
	"context"
	"fmt"
	"time"

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
		for i := 0; i < 10; i++ {
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
	if err := service.Run(); err != nil {
		fmt.Println("Service running with error: ", err)
	}
}
