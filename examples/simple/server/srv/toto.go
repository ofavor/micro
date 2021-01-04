package srv

import (
	"context"
	"fmt"

	"github.com/ofavor/micro-lite/examples/simple/server/toto"
)

type TotoHandler struct {
}

func (h *TotoHandler) Multiply(ctx context.Context, req *toto.Request, rsp *toto.Response) error {
	fmt.Println("toto multiply is doing ....")
	// time.Sleep(6 * time.Second)
	rsp.Result = req.Val1 * req.Val2
	return nil
}
