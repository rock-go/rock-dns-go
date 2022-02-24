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

func (s *server) startL(L *lua.LState) int {
	if s.CodeVM() != L.CodeVM() {
		L.RaiseError("%s must be %s run" , s.Name() , s.CodeVM())
		return 0
	}

	xEnv.Start(s , func(err error) {
		L.RaiseError("%v" , err)
	})

	return 0
}

func (s *server) Index(L *lua.LState , key string) lua.LValue {
	switch key {

	case "start":
		return L.NewFunction(s.startL)

	case "pipe":
		return L.NewFunction(s.pipeL)
	case "to":
		return L.NewFunction(s.toL)
	}
	return lua.LNil
}

func (s *server) NewIndex(L *lua.LState, key string, val lua.LValue) {
	s.addRoute(L, key, val)
}