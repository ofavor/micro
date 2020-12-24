package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ofavor/micro-lite/client/selector"
	"github.com/ofavor/micro-lite/examles/simple/server/toto"

	"github.com/ofavor/micro-lite"
	"github.com/ofavor/micro-lite/client"
)

func main() {
	// if true {
	// 	arr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	// 	for i, a := range arr {
	// 		fmt.Println("Before:", arr, i, a)
	// 		if a == 1 {
	// 			arr = append(arr[0:i], arr[i+1:]...)
	// 		}
	// 		fmt.Println("After: ", arr, i, a)
	// 	}
	// 	return
	// }
	service := micro.NewService(
		micro.LogLevel("debug"),
		micro.Name("simple.client"),
		// micro.Version("latest"),
		micro.Address(":8889"),
	)
	// go func() {
	// 	for i := 0; ; i++ {
	// 		time.Sleep(2 * time.Second)
	// 		f := foo.NewFooService(service.Client())
	// 		req := &foo.Request{
	// 			Name: "Bob ",
	// 			Age:  uint32(i) + 1,
	// 		}
	// 		rsp, err := f.Bar(context.Background(), req)
	// 		if err != nil {
	// 			fmt.Println("Error:", err)
	// 		} else {
	// 			fmt.Println("Response:", rsp)
	// 		}
	// 	}
	// }()
	go func() {
		for {
			time.Sleep(2 * time.Second)
			t := toto.NewTotoService("simple.server", service.Client())
			req := &toto.Request{
				Val1: 30,
				Val2: 40,
			}
			rsp, err := t.Multiply(context.Background(), req,
				// client.WithSelectOption(selector.WithAddressFilter([]string{"172.20.10.2:8888"})), // test selector with ip address
				client.WithSelectOption(selector.WithIDFilter([]string{"srv1"})),
			)

			if err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println(req.Val1, " x ", req.Val2, " = ", rsp.Result)
			}
		}
	}()
	// go func() {
	// 	for {
	// 		time.Sleep(2 * time.Second)
	// 		t := toto.NewTotoService(service.Client())
	// 		req := &toto.Request{
	// 			Val1: 30,
	// 			Val2: 40,
	// 		}
	// 		rsp, err := t.Multiply(context.Background(), req)

	// 		if err != nil {
	// 			fmt.Println("Error:", err)
	// 		} else {
	// 			fmt.Println(req.Val1, " x ", req.Val2, " = ", rsp.Result)
	// 		}
	// 	}
	// }()
	if err := service.Run(); err != nil {
		fmt.Println("Service running with error: ", err)
	}
}
