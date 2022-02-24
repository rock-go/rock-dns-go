package client

import (
	"github.com/miekg/dns"
	"github.com/rock-go/rock/lua"
	"github.com/rock-go/rock/xbase"
)

var xEnv *xbase.EnvT

func (dc *dnsClient) Query(L *lua.LState) lua.LValue {

	str := L.CheckString(1) + "."
	msg := dns.Msg{}
	msg.SetQuestion(str, dns.TypeANY)
	cli := dc.Client()

	r, rtt, err := cli.Exchange(&msg, dc.cfg.Resolve)
	return L.NewAnyData(&Reply{r, rtt, err})
}

func newLuaDnsClient(L *lua.LState) int {
	cfg := newConfig(L)
	cli := newDnsClient(cfg)
	L.Push(L.NewAnyData(cli))
	return 1
}

func LuaInjectApi(env *xbase.EnvT, kv lua.UserKV) {
	kv.Set("client", lua.NewFunction(newLuaDnsClient))
	xEnv = env
}
