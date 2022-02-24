package client

import (
	"github.com/miekg/dns"
	"reflect"
	"sync"
)

var DnsCliType = reflect.TypeOf((*dnsClient)(nil)).String()

type dnsClient struct {
	cfg  *config
	pool sync.Pool
}

func newDnsClient(cfg *config) *dnsClient {
	cli := &dnsClient{cfg: cfg}

	cli.pool = sync.Pool{
		New: func() interface{} {
			return &dns.Client{Timeout: cfg.ToTimer()}
		},
	}

	return cli
}

func (dc *dnsClient) Client() *dns.Client {
	return dc.pool.Get().(*dns.Client)
}

func (dc *dnsClient) Free(cli *dns.Client) {
	dc.pool.Put(cli)
}
