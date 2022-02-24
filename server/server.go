package server

import (
	"github.com/miekg/dns"
	"github.com/rock-go/rock-dns-go/internal"
	"github.com/rock-go/rock/logger"
	"github.com/rock-go/rock/lua"
	"github.com/rock-go/rock/pipe"
	"reflect"
)

var (
	dnsServerTypeOf = reflect.TypeOf((*server)(nil)).String()
)

type server struct {
	lua.Super

	cfg *config
	fd  *dns.Server
}

func newDnsServer(cfg *config) *server {
	s := &server{cfg: cfg}
	s.V(lua.INIT, dnsServerTypeOf)
	return s
}

func (s *server) Name() string {
	return s.cfg.Name
}

func (s *server) Start() error {
	fd := &dns.Server{Addr: s.cfg.Bind.Hostname(), Net: s.cfg.Bind.Scheme(), Handler: s}
	var err error
	xEnv.Spawn(100, func() {
		err = fd.ListenAndServe()
	})

	if err != nil {
		return err
	}
	s.fd = fd
	logger.Errorf("%s dns server start succeed", s.Name())
	return nil
}

func (s *server) Close() error {
	if s.fd != nil {
		return s.fd.Shutdown()
	}

	return nil
}

func (s *server) call(val lua.LValue, ctx *internal.Context) {
	switch val.Type() {
	case lua.LTString:
		ctx.Say(val.String())

	case lua.LTFunction:
		cp := xEnv.P(val.(*lua.LFunction))
		cp.NRet = 1
		co := xEnv.Clone(s.cfg.co)
		defer xEnv.Free(co)

		err := co.CallByParam(cp, ctx)
		if err != nil {
			ctx.Error(err)
			return
		}
	}
}

func (s *server) router(name string , ctx *internal.Context) {
	n := s.Len()
	if n == 0 {
		return
	}

	for i := 0 ; i < n ; i++ {
		r := s.cfg.route[i]
		if !r.match(name) {
			continue
		}

		ctx.Hit = true
		r.handle(ctx)
		return
	}
}

func (s *server) pipe(ctx *internal.Context) {
	pipe.Do(s.cfg.pipe , ctx , s.cfg.co , func(err error) {
		xEnv.Errorf("%s pipe call fail %v" , s.Name() , err)
	})
}

func (s *server) Query(ctx *internal.Context) {
	questions := ctx.Question()
	s.cfg.mutex.RLock()
	defer s.cfg.mutex.RUnlock()

	for id, q := range questions {
		ctx.SetIndexId(id)

		//判断直接处理结果
		s.router(q.Name , ctx)

		//pipe
		s.pipe( ctx )
	}
}

func (s *server) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	ctx := internal.New(w, r, s.cfg.co.CodeVM())
	ctx.Region = s.cfg.Region

	switch r.Opcode {
	case dns.OpcodeQuery:
		s.Query(ctx)
	case dns.OpcodeNotify:
	case dns.OpcodeStatus:
	case dns.OpcodeUpdate:

	}

	ctx.Reply()
	s.log(ctx)
}

func (s *server) log(ctx *internal.Context) {
	if s.cfg.sdk == nil {
		return
	}
	s.cfg.sdk(lua.S2L(ctx.String()) , s.cfg.co)
}

