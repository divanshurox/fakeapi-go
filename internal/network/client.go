package network

import (
	"FakeAPI/internal/logger"
	"context"
	"errors"
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	"go.uber.org/zap"
	"io"
	"net/http"
	"sync"
	"time"
)

var (
	DefaultTimeout = 10
	once           sync.Once
	transport      *http.Transport
	clientMap      clientManager
)

type client struct {
	headers http.Header
	body    io.Reader
	timeout int
	client  string
	ctx     *context.Context
}

type ClientInterface interface {
	Headers(headers http.Header) *client
	Body(body io.Reader) *client
	Timeout(timeout int) *client
	Client(client string) *client
	WithContext(ctx *context.Context) *client
	Get(url string) (*http.Response, error)
	Put(url string) (*http.Response, error)
	Post(url string) (*http.Response, error)
	Do(method, url string) (*http.Response, error)
}

func NewClient() *client {
	return &client{client: "default", timeout: DefaultTimeout}
}

func init() {
	once.Do(func() {
		defaultTransport, _ := http.DefaultTransport.(*http.Transport)
		transport = defaultTransport
		transport.MaxIdleConns = 100
		transport.MaxIdleConnsPerHost = 10
		clientMap = clientManager{
			pool: make(map[string]*http.Client),
		}
	})
}

func (c *client) Headers(headers http.Header) *client {
	if len(c.headers) == 0 {
		c.headers = make(http.Header)
	}
	for k, v := range headers {
		c.headers[k] = v
	}
	return c
}

func (c *client) Body(body io.Reader) *client {
	c.body = body
	return c
}

func (c *client) Client(client string) *client {
	c.client = client
	return c
}
func (c *client) Timeout(timeout int) *client {
	c.timeout = timeout
	return c
}

func (c *client) WithContext(ctx *context.Context) *client {
	c.ctx = ctx
	return c
}

func (c *client) Get(url string) (*http.Response, error) {
	return c.Do("GET", url)
}

func (c *client) Put(url string) (*http.Response, error) {
	return c.Do("PUT", url)
}

func (c *client) Post(url string) (*http.Response, error) {
	return c.Do("POST", url)
}

func (c *client) Do(method, url string) (*http.Response, error) {
	clientInstance := clientMap.getClient(c.client)
	if clientInstance == nil {
		clientInstance = clientMap.setClient(c.client, c.timeout)
	}
	req, err := http.NewRequest(method, url, c.body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(*getTimedContext(c.timeout))
	for k, v := range c.headers {
		req.Header[k] = v
	}
	var res *http.Response
	var apiErr error
	err = hystrix.Do(c.client, func() error {
		res, apiErr = clientInstance.Do(req)
		if apiErr != nil {
			return apiErr
		}
		if res != nil && (res.StatusCode < 200 || res.StatusCode > 299) {
			return fmt.Errorf("non success status code found - %d", res.StatusCode)
		}
		return nil
	}, func(err error) error {
		logger.GetLogger().Error("fallback for circuit break", zap.Error(err))
		return errors.New("fallback for circuit break")
	})
	return res, err
}

func getTimedContext(timeout int) *context.Context {
	ctx, cancel := context.WithCancel(context.TODO())
	time.AfterFunc(time.Duration(timeout)*time.Second, func() {
		cancel()
	})
	return &ctx
}
