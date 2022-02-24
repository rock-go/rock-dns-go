package client

import (
	"github.com/miekg/dns"
	"github.com/rock-go/rock/lua"
)

type RR struct {
	data dns.RR
}

func (rr *RR) DisableReflect() {}

func (rr *RR) Get(L *lua.LState, key string) lua.LValue {
	switch key {
	case "name":
		return lua.S2L(rr.data.Header().Name)
	case "type":
		return lua.S2L(dns.Type(rr.data.Header().Rrtype).String())
	case "ttl":
		return lua.LNumber(rr.data.Header().Ttl)
	case "raw":
		return lua.S2L(rr.data.String())
	case "header":
		return lua.S2L(rr.data.Header().String())

	case "value":
		switch record := rr.data.(type) {
		case *dns.A:
			return lua.S2L(record.A.String())
		case *dns.CNAME:
			return lua.S2L(record.Target)
		case *dns.AAAA:
			return lua.S2L(record.AAAA.String())
		case *dns.DNAME:
			return lua.S2L(record.Target)
		default:
			return lua.S2L(rr.data.String())
		}
	}

	return lua.LNil
}
