package server

import (
	"github.com/rock-go/rock/lua"
	"github.com/rock-go/rock/pipe"
)

func (s *server) pipeL(L *lua.LState) int {
	pp := pipe.Check(L)
	if len(pp) > 0 {
		s.cfg.pipe = append(s.cfg.pipe , pp...)
	}
	return 0
}

func (s *server) toL(L *lua.LState) int {
	s.cfg.sdk = pipe.LValue(L.Get(1))
	return 0
}

func (s *server) Index(L *lua.LState , key string) lua.LValue {
	switch key {
	case "on":
		return L.NewFunction(s.pipeL)
	case "to":
		return L.NewFunction(s.toL)
	}
	return lua.LNil
}

func (s *server) NewIndex(L *lua.LState, key string, val lua.LValue) {
	s.addRoute(L, key, val)
}