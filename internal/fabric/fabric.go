package fabric

import (
	"fmt"
	"github.com/unpoller/unifi"
	"sync"
)

type Type string

const (
	Unifi Type = "unifi"
)

type Device struct {
	MAC string
	IP  string
}

type Client interface {
	GetFabricDevices() ([]*Device, error)
}

type ClientConfig struct {
	Fabric   Type
	Username string
	Password string
	Url      string
}

func (c *ClientConfig) Key() string {
	return fmt.Sprintf("%s:%s:%s", c.Fabric, c.Username, c.Url)
}

type ClientGetter interface {
	GetClient(*ClientConfig) (Client, error)
}

func NewCachedFactory() *CachedFactory {
	return &CachedFactory{
		cache: make(map[string]Client),
	}
}

type Factory struct{}

func (f *Factory) GetClient(cfg *ClientConfig) (Client, error) {
	switch cfg.Fabric {
	case Unifi:
		client, err := unifi.NewUnifi(&unifi.Config{
			User: cfg.Username,
			Pass: cfg.Password,
			URL:  cfg.Url,
		})
		if err != nil {
			return nil, err
		}
		return &UnifiClient{
			client,
		}, nil
	default:
		return nil, fmt.Errorf("unknown fabric type: %v", cfg.Fabric)
	}
}

type CachedFactory struct {
	mut   sync.RWMutex
	cache map[string]Client
}

func (f *CachedFactory) GetClient(cfg *ClientConfig) (Client, error) {
	f.mut.RLock()
	c, ok := f.cache[cfg.Key()]
	f.mut.RUnlock()

	if ok {
		return c, nil
	}

	f.mut.Lock()
	defer f.mut.Unlock()

	c, err := (&Factory{}).GetClient(cfg)
	if err != nil {
		return nil, err
	}
	f.cache[cfg.Key()] = c

	return c, err
}
