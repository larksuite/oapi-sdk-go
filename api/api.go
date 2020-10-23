package api

import (
	"github.com/larksuite/oapi-sdk-go/api/core/handlers"
	"github.com/larksuite/oapi-sdk-go/api/core/request"
	"github.com/larksuite/oapi-sdk-go/api/core/response"
	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/core/config"
	"sync"
)

func Send(ctx *core.Context, conf *config.Config, req *request.Request) error {
	conf.WithContext(ctx)
	req.WithContext(ctx)
	handlers.Handle(ctx, req)
	if req.Err == nil {
		return nil
	}
	return response.ToError(req.Err)
}

type ErrorCallback = func(reqCall ReqCaller, err error)

type ReqCallDone struct {
	ctx    *core.Context
	Result interface{}
	Err    interface{}
}

type ReqCaller interface {
	Ctx() *core.Context
	Do() (interface{}, error)
}

type BatchReqCall struct {
	wg            sync.WaitGroup
	reqCalls      []ReqCaller
	ReqCallDos    []*ReqCallDone
	errorCallback ErrorCallback
}

func NewBatchReqCall(errorCallback ErrorCallback, reqCalls ...ReqCaller) *BatchReqCall {
	reqCallDos := make([]*ReqCallDone, 0, len(reqCalls))
	for _, reqCall := range reqCalls {
		reqCallDos = append(reqCallDos, &ReqCallDone{
			ctx: reqCall.Ctx(),
		})
	}
	return &BatchReqCall{
		reqCalls:      reqCalls,
		ReqCallDos:    reqCallDos,
		errorCallback: errorCallback,
	}
}

func (c *BatchReqCall) Do() *BatchReqCall {
	c.wg.Add(len(c.reqCalls))
	for i, reqCaller := range c.reqCalls {
		go func(i int, rc ReqCaller) {
			defer c.wg.Done()
			result, err := rc.Do()
			c.ReqCallDos[i].Result = result
			c.ReqCallDos[i].Err = err
			if err != nil {
				c.errorCallback(rc, err)
			}
		}(i, reqCaller)
	}
	c.wg.Wait()
	return c
}
