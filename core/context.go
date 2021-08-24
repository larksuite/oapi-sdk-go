package core

import (
	"context"
	"github.com/larksuite/oapi-sdk-go/core/constants"
	"sync"
	"time"
)

type Context struct {
	c  context.Context
	mu sync.RWMutex
	m  map[string]interface{}
}

// Deprecated, please use `WrapContext`
func WarpContext(c context.Context) *Context {
	return &Context{
		c: c,
	}
}

func WrapContext(c context.Context) *Context {
	return &Context{
		c: c,
	}
}

func (c *Context) Set(key string, value interface{}) {
	c.mu.Lock()
	if c.m == nil {
		c.m = make(map[string]interface{})
	}
	c.m[key] = value
	c.mu.Unlock()
}

// Get returns the value for the given key, ie: (value, true).
// If the value does not exists it returns (nil, false)
func (c *Context) Get(key string) (value interface{}, exists bool) {
	c.mu.RLock()
	value, exists = c.m[key]
	c.mu.RUnlock()
	return
}

func (c *Context) GetHeader() *OapiHeader {
	if header, ok := c.Get(constants.HTTPHeader); ok {
		return header.(*OapiHeader)
	}
	return NewOapiHeader(map[string][]string{})
}

func (c *Context) GetRequestID() string {
	header := c.GetHeader()
	logID := header.GetFirstValue(constants.HTTPHeaderKeyLogID)
	if logID != "" {
		return logID
	}
	return header.GetFirstValue(constants.HTTPHeaderKeyRequestID)
}

func (c *Context) GetHTTPStatusCode() int {
	if statusCode, ok := c.Get(constants.HTTPKeyStatusCode); ok {
		return statusCode.(int)
	}
	return 0
}

func (c *Context) Deadline() (deadline time.Time, ok bool) {
	return c.c.Deadline()
}

func (c *Context) Done() <-chan struct{} {
	return c.c.Done()
}

func (c *Context) Err() error {
	return c.c.Err()
}

func (c *Context) Value(key interface{}) interface{} {
	if keyAsString, ok := key.(string); ok {
		val, _ := c.Get(keyAsString)
		return val
	}
	return c.c.Value(key)
}
