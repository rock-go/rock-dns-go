package client

import (
	"github.com/miekg/dns"
	"github.com/rock-go/rock/lua"
	"time"
)

//func (c *Client) Exchange(m *Msg, address string) (r *Msg, rtt time.Duration, err error) {

type Reply struct {
	r   *dns.Msg
	rtt time.Duration
	err error
}


func (r *Reply) record(L *lua.LState) *lua.LTable {
	n := len(r.r.Answer)
	tab := L.CreateTable(n, 0)
	for i := 0; i < n; i++ {
		record := r.r.Answer[i].String()
		tab.RawSetInt(i, lua.S2L(record))
	}

	return tab
}

func (r *Reply) pairs(L *lua.LState) int {
	n := len(r.r.Answer)
	cp := xEnv.P(L.CheckFunction(1))
	co := xEnv.Clone(L)

	for i := 0; i < n; i++ {
		rr := &RR{data: r.r.Answer[i]}
		ud := L.NewAnyData(rr)
		err := co.CallByParam(cp, ud)
		if err != nil {
			L.RaiseError("%v", err)
			return 0
		}

		lv := co.Get(-1)
		switch lv.Type() {
		case lua.LTBool:
			if bool(lv.(lua.LBool)) {
				return 0
			}
			continue
		default:
			//todo
		}
		co.SetTop(0)
	}

	return 0
}

func (r *Reply) Index(L *lua.LState, key string) lua.LValue {

	switch key {
	case "record":
		return r.record(L)
	case "record_json":
		return lua.JsonMarshal(L, r.r.Answer)

	case "ERR":
		if r.err == nil {
			return lua.LNil
		}
		return lua.S2L(r.err.Error())

	case "pairs":
		return L.NewFunction(r.pairs)

	default:
		return lua.LNil
	}
}
