package xclient

import (
	"context"
	"geeRpc"
	"io"
	"reflect"
	"sync"
)

type XClient struct {
	d      Discovery
	mode   SelectMode
	opt    *geeRpc.Option
	mu     sync.Mutex
	client map[string]*geeRpc.Client
}

var _ io.Closer = (*XClient)(nil)

func NewXClient(d Discovery, mode SelectMode, opt *geeRpc.Option) *XClient {
	return &XClient{
		d:      d,
		mode:   mode,
		opt:    opt,
		client: make(map[string]*geeRpc.Client),
	}
}

func (xc *XClient) Close() error {
	xc.mu.Lock()
	defer xc.mu.Unlock()
	for key, client := range xc.client {
		_ = client.Close()
		delete(xc.client, key)
	}
	return nil
}

func (xc *XClient) dial(rpcAddr string) (*geeRpc.Client, error) {
	xc.mu.Lock()
	defer xc.mu.Unlock()
	client, ok := xc.client[rpcAddr]
	if ok && !client.IsAvailable() {
		_ = client.Close()
		delete(xc.client, rpcAddr)
		client = nil
	}
	if client == nil {
		var err error
		client, err = geeRpc.XDial(rpcAddr, xc.opt)
		if err != nil {
			return nil, err
		}
		xc.client[rpcAddr] = client
	}
	return client, nil
}

func (xc *XClient) call(rpcAddr string, ctx context.Context, serviceMethod string, args, reply interface{}) error {
	client, err := xc.dial(rpcAddr)
	if err != nil {
		return err
	}
	return client.Call(ctx, serviceMethod, args, reply)
}

func (xc *XClient) Call(ctx context.Context, serviceMethod string, args, reply interface{}) error {
	rpcAddr, err := xc.d.Get(xc.mode)
	if err != nil {
		return err
	}
	return xc.call(rpcAddr, ctx, serviceMethod, args, reply)
}

func (xc *XClient) Broadcast(ctx context.Context, serviceMethod string, args, reply interface{}) error {
	services, err := xc.d.GetAll()
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	var mu sync.Mutex
	var e error
	replyDone := reply == nil
	ctx, cancel := context.WithCancel(ctx)
	for _, rpcAddr := range services {
		wg.Add(1)
		go func(rpcAddr string) {
			defer wg.Done()
			var clondReply interface{}
			if reply != nil {
				clondReply = reflect.New(reflect.ValueOf(reply).Elem().Type()).Interface()
			}
			err := xc.call(rpcAddr, ctx, serviceMethod, args, clondReply)
			mu.Lock()
			defer mu.Unlock()
			if err != nil && e == nil {
				e = err
				cancel()
			}
			if err == nil && !replyDone {
				reflect.ValueOf(reply).Elem().Set(reflect.ValueOf(clondReply).Elem())
				replyDone = true
			}
		}(rpcAddr)
	}
	wg.Wait()
	return e
}
