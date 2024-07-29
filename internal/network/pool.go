package network

import (
	"github.com/afex/hystrix-go/hystrix"
	"net/http"
	"sync"
	"time"
)

const (
	DEFAULT_MAX_CONCURRENCY  = 10
	DEFAULT_ERROR_THRESHOLD  = 25
	DEFAULT_VOLUME_THRESHOLD = 5
	DEFAULT_SLEEP_WINDOW     = 3000
)

type clientManager struct {
	sync.RWMutex
	pool map[string]*http.Client
}

type clientManagerInterface interface {
	getClient(client string) *http.Client
	setClient(client string, timeout int) *http.Client
}

func (c *clientManager) getClient(client string) *http.Client {
	c.RLock()
	defer c.RUnlock()
	return c.pool[client]
}

func (c *clientManager) setClient(client string, timeout int) *http.Client {
	c.Lock()
	defer c.Unlock()
	if c.pool[client] == nil {
		c.pool[client] = &http.Client{Transport: transport, Timeout: time.Duration(timeout) * time.Second}
	}
	hystrix.ConfigureCommand(client, hystrix.CommandConfig{
		Timeout:                timeout * 1000,
		ErrorPercentThreshold:  DEFAULT_ERROR_THRESHOLD,
		MaxConcurrentRequests:  DEFAULT_MAX_CONCURRENCY,
		SleepWindow:            DEFAULT_SLEEP_WINDOW,
		RequestVolumeThreshold: DEFAULT_VOLUME_THRESHOLD,
	})
	return c.pool[client]
}
