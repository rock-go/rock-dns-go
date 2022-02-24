package internal

import (
	"github.com/miekg/dns"
	"github.com/rock-go/rock/auxlib"
	"github.com/rock-go/rock/lua"
	"strings"
)

func (ctx *Context) Type() lua.LValueType                   { return lua.LTObject }
func (ctx *Context) AssertFloat64() (float64, bool)         { return 0, false }
func (ctx *Context) AssertString() (string, bool)           { return "", false }
func (ctx *Context) AssertFunction() (*lua.LFunction, bool) { return nil, false }
func (ctx *Context) Peek() lua.LValue                       { return ctx }

func (ctx *Context) sayL(L *lua.LState) int {
	addr := L.CheckString(1)

	rr, err := ctx.RR(addr)
	if err != nil {
		L.Pushf("response dns %v", err)
		return 1
	}
	ctx.response.Answer = append(ctx.response.Answer, rr)
	return 0
}

func (ctx *Context) sayAL(L *lua.LState) int {
	addr := L.CheckString(1)

	if !auxlib.Ipv4(addr) {
		L.Push(lua.S2L("invalid ipv4"))
		return 1
	}

	ctx.say(addr , "A")
	return 0
}

func (ctx *Context) sayAAAAL(L *lua.LState) int {
	addr := L.CheckString(1)

	if !auxlib.Ipv6(addr) {
		L.Push(lua.S2L("invalid ipv6"))
		return 1
	}

	ctx.say(addr , "AAAA")
	return 0
}

func (ctx *Context) cnameL(L *lua.LState) int {
	cname := L.CheckString(1)
	ctx.say(cname , "CNAME")
	return 0
}

func (ctx *Context) suffixTrimL(L *lua.LState) int {
	str := L.CheckString(1)
	name := ctx.QName()
	nl := len(name)
	sl := len(str)

	if nl > sl && name[nl-sl:] == str {
		L.Push(lua.S2L(name[:nl-sl]))
		return 1
	}
	L.Push(lua.LNil)
	return 1
}

func (ctx *Context) prefixTrimL(L *lua.LState) int {
	str := L.CheckString(1)

	name := ctx.QName()
	nl := len(name)
	sl := len(str)

	if nl > sl && name[0:sl] == str {
		L.Push(lua.S2L(name[sl:]))
		return 1
	}

	L.Push(lua.LNil)
	return 1
}

func (ctx *Context) suffixL(L *lua.LState) int {
	return helper(L , ctx.QName() , strings.HasSuffix)

}

func (ctx *Context) prefixL(L *lua.LState) int {
	return helper(L , ctx.QName() , strings.HasPrefix)
}


func (ctx *Context) Index(L *lua.LState, key string) lua.LValue {
	switch key {
	case "say":
		return lua.NewFunction(ctx.sayL)
	case "A":
		return lua.NewFunction(ctx.sayAL)
	case "AAAA":
		return lua.NewFunction(ctx.sayAAAAL)
	case "CNAME":
		return lua.NewFunction(ctx.cnameL)
	case "suffix_trim":
		return lua.NewFunction(ctx.suffixTrimL)
	case "prefix_trim":
		return lua.NewFunction(ctx.prefixTrimL)
	case "suffix":
		return lua.NewFunction(ctx.suffixL)
	case "prefix":
		return lua.NewFunction(ctx.prefixL)
	case "class":
		return lua.LString(dns.Class(ctx.Q().Qclass).String())
	case "type":
		return lua.LString(ctx.QType())
	case "dns":
		return lua.LString(ctx.request.Question[ctx.id].Name)
	case "op_code":
		return lua.LInt(ctx.request.Opcode)
	case "extra":
		return lua.LString(rrToStr(ctx.request.Extra))
	case "ns":
		return lua.LString(rrToStr(ctx.request.Ns))
	case "answer":
		return lua.LString(rrToStr(ctx.response.Answer))
	case "hit":
		return lua.LBool(ctx.Hit)
	case "region":
		return lua.LString(ctx.region())
	case "addr":
		return lua.LString(ctx.Addr())
	case "port":
		return lua.LInt(ctx.Port())
	case "hdr_raw":
		return lua.LString(ctx.request.MsgHdr.String())
	case "question_raw":
		return lua.LString(qqToStr(ctx.request.Question))
	}

	return lua.LNil
}
