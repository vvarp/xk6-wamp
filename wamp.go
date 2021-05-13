package wamp

import (
	"context"
	"log"

	"github.com/dop251/goja"
	"github.com/gammazero/nexus/v3/client"
	"github.com/gammazero/nexus/v3/wamp"
	"go.k6.io/k6/js/common"
	"go.k6.io/k6/js/modules"
)

func init() {
	modules.Register("k6/x/wamp", new(WAMP))
}

type WAMP struct{}

func (r *WAMP) XClient(ctxPtr *context.Context, addr string, opts *client.Config) interface{} {
	rt := common.GetRuntime(*ctxPtr)

	var c *client.Client
	c, err := client.ConnectNet(*ctxPtr, addr, *opts)

	if err != nil {
		log.Fatal(err)
	}

	return common.Bind(rt, &Client{
		client:        c,
		ctx:           *ctxPtr,
		eventHandlers: make(map[wamp.ID]goja.Callable),
	}, ctxPtr)
}

type Client struct {
	client        *client.Client
	ctx           context.Context
	eventHandlers map[wamp.ID]goja.Callable
}

func (c *Client) GetSessionID() uint64 {
	return uint64(c.client.ID())
}

func (c *Client) IsConnected() bool {
	return c.client.Connected()
}

func (c *Client) Subscribe(topic string, options wamp.Dict, handler goja.Value) uint64 {
	err := c.client.Subscribe(topic, c.handleSubscribeEvent, options)

	if err != nil {
		log.Print(err)
		return 0
	}

	subId, ok := c.client.SubscriptionID(topic)
	if ok {
		handlerFunc, isFunc := goja.AssertFunction(handler)
		if isFunc {
			c.eventHandlers[subId] = handlerFunc
		}
		return uint64(subId)
	}
	return 0
}

func (c *Client) Publish(topic string, options wamp.Dict, args wamp.List, kwargs wamp.Dict) {
	err := c.client.Publish(topic, options, args, kwargs)
	if err != nil {
		log.Print(err)
	}
}

func (c *Client) Disconnect() {
	err := c.client.Close()
	if err != nil {
		log.Printf("Error while disconnecting: %v", err)
	}
}

func (c *Client) handleSubscribeEvent(event *wamp.Event) {
	rt := common.GetRuntime(c.ctx)
	if handler, ok := c.eventHandlers[event.Subscription]; ok {
		if _, err := handler(goja.Undefined(), rt.ToValue(event.Arguments), rt.ToValue(event.ArgumentsKw)); err != nil {
			common.Throw(rt, err)
		}
	}
}
