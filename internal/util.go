package internal

import (
	"bytes"
	"github.com/miekg/dns"
	"github.com/rock-go/rock/lua"
)

func rrToStr(rrs []dns.RR) []byte {
	n := len(rrs)
	if n == 0 {
		return nil
	}

	var buff bytes.Buffer
	for i := 0; i < n; i++ {
		rr := rrs[i]
		if rr == nil {
			continue
		}

		buff.WriteString(rr.String())
		buff.WriteByte(',')
	}

	return buff.Bytes()
}

func qqToStr(qqs []dns.Question) []byte {
	n := len(qqs)
	if n == 0 {
		return nil
	}

	var buff bytes.Buffer
	for i := 0; i < n; i++ {
		buff.WriteString(qqs[i].String()[1:])
		buff.WriteByte('\n')
	}

	return buff.Bytes()
}

func helper(L *lua.LState, name string, fn func(string, string) bool) int {
	n := L.GetTop()
	if n < 2 {
		L.Push(lua.LFalse)
		return 1
	}

	for i := 2; i <= n; i++ {
		val := L.CheckString(i)
		if fn(name, val) {
			L.Push(lua.LTrue)
			return 1
		}
	}

	L.Push(lua.LFalse)
	return 1
}
