package dns

import (
	"github.com/rock-go/rock-dns-go/client"
	"github.com/rock-go/rock-dns-go/server"
	"github.com/rock-go/rock/lua"
	"github.com/rock-go/rock/xbase"
)

func LuaInjectApi(env *xbase.EnvT) {
	kv := lua.NewUserKV()
	client.LuaInjectApi(env, kv)
	server.LuaInjectApi(env, kv)
	env.Set("dns", kv)
}
