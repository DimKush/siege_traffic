package dns

import (
	"context"
	"github.com/rs/dnscache"
	"sort"
	"sync"
	"time"
)

type Lookup func(string) string

type DNSResolver struct {
	done     chan struct{}
	resolver *dnscache.Resolver
	wg       sync.WaitGroup
}

func NewDnsResolver() *DNSResolver {
	r := &DNSResolver{
		done:     make(chan struct{}, 1),
		resolver: &dnscache.Resolver{},
	}
	r.start()

	return r
}

func (c *DNSResolver) start() {
	c.wg.Add(1)
	defer c.wg.Done()

	go func() {
		t := time.NewTicker(1 * time.Minute)
		defer t.Stop()

		for {
			select {
			case <-c.done:
				return
			case <-t.C:
				c.resolver.Refresh(true)
			}
		}
	}()
}

func (c *DNSResolver) Close() {
	c.done <- struct{}{}
	c.wg.Wait()
}

func (c *DNSResolver) Lookup(ip string) string {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	addrs, err := c.resolver.LookupAddr(ctx, ip)
	if err != nil {
		return ip
	}

	if len(addrs) == 0 {
		return ip
	}

	if len(addrs) == 1 {
		return addrs[0]
	}

	sort.Strings(addrs)
	return addrs[0]
}
