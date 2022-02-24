package client

import (
	"github.com/rock-go/rock/lua"
	"time"
)

type config struct {
	Resolve string `lua:"resolve" type:"string"`
	Timeout int    `lua:"timeout" type:"int"`
}

func newConfig(L *lua.LState) *config {
	switch L.GetTop() {
	case 0:
		return &config{Resolve: "114.114.114.114" , Timeout: 3}
	case 1:
		return &config{Resolve: L.CheckString(1) , Timeout: 3}

	default:
		return &config{Resolve: L.CheckString(1) , Timeout: L.IsInt(2)}

	}
}

func (cfg *config) verify() error {
	return nil
}

func (cfg *config) ToTimer() time.Duration {
	return time.Duration(cfg.Timeout) * time.Millisecond
}
