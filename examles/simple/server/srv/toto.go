package srv

import (
	"context"
	"time"

	"github.com/ofavor/micro-lite/examles/simple/server/toto"
)

type TotoHandler struct {
}

func (h *TotoHandler) Multiply(ctx context.Context, req *toto.Request, rsp *toto.Response) error {
	time.Sleep(6 * time.Second)
	rsp.Result = req.Val1 * req.Val2
	return nil
}
