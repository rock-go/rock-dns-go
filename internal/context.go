package internal

import (
	"fmt"
	"github.com/miekg/dns"
	"github.com/rock-go/rock/auxlib"
	"github.com/rock-go/rock/json"
	"github.com/rock-go/rock/logger"
	"github.com/rock-go/rock/node"
	"github.com/rock-go/rock/region"
	"strconv"
	"strings"
	"time"
)

type Context struct {
	request  *dns.Msg
	response *dns.Msg
	from     string
	w        dns.ResponseWriter
	err      error
	id       int

	Hit      bool
	Direct   *bool
	Region   *region.Region
}

func New(w dns.ResponseWriter, r *dns.Msg, from string) *Context {
	m := &dns.Msg{}
	m.SetReply(r)
	m.Compress = false

	return &Context{
		request:  r,
		response: m,
		from:     from,
		w:        w,
	}
}

func (ctx *Context) Deadline() (deadline time.Time, ok bool) {
	return
}

func (ctx *Context) Done() <-chan struct{} {
	ctx.w.WriteMsg(ctx.response)
	return nil
}

func (ctx *Context) Err() error {
	return ctx.err
}

func (ctx *Context) Value(key interface{}) interface{} {
	return nil
}

func (ctx *Context) Reply() {
	if ctx.err != nil {
		return
	}

	err := ctx.w.WriteMsg(ctx.response)
	if err != nil {
		ctx.err = err
	}
}

func (ctx *Context) Error(err error) {
	ctx.err = err
}

func (ctx *Context) Data() *dns.Msg {
	return ctx.response
}

func (ctx *Context) SetIndexId(id int) {
	ctx.id = id
}

func (ctx *Context) String() string {
	buff := json.NewEncoder()
	buff.Tab("")
	buff.KV("ID", node.ID())
	buff.KV("inet", node.LoadAddr())
	buff.KV("hdr_id", ctx.request.Id)
	buff.KV("response", ctx.request.Response)
	buff.KV("op_code", ctx.request.Opcode)
	buff.KV("zero", ctx.request.Zero)
	buff.KV("r_code", ctx.request.Rcode)
	buff.KV("compress", ctx.request.Compress)
	buff.KV("question", qqToStr(ctx.request.Question))
	buff.KV("answer", rrToStr(ctx.response.Answer))
	buff.KV("ns", rrToStr(ctx.request.Ns))
	buff.KV("extra", rrToStr(ctx.request.Extra))
	buff.KV("remote_addr", ctx.w.RemoteAddr().String())
	buff.KV("region", ctx.region())

	buff.KV("authoritative", ctx.request.Authoritative)
	buff.KV("recursion_desired", ctx.request.RecursionDesired)
	buff.KV("recursion_available", ctx.request.RecursionAvailable)
	buff.KV("authenticate_data", ctx.request.AuthenticatedData)
	buff.KV("error", ctx.err)
	buff.End("}")
	return auxlib.B2S(buff.Bytes())
}

func (ctx *Context) Q() dns.Question {
	return ctx.response.Question[ctx.id]
}

func (ctx *Context) Question() []dns.Question {
	return ctx.request.Question
}

func (ctx *Context) Say(val string) {
	q := ctx.Q()

	if t, ok := dns.TypeToString[q.Qtype]; !ok {
		ctx.err = fmt.Errorf("not found qtype %d", q.Qtype)
	} else {
		rr, err := dns.NewRR(q.Name + " " + t + " " + val)
		if err != nil {
			ctx.err = err
		}
		ctx.response.Answer = append(ctx.response.Answer, rr)
	}
}

func (ctx *Context) say(v string , t string) {
	if len(t) == 0 {
		return
	}
	q := ctx.Q()

	rr, err := dns.NewRR(q.Name + " " + t + " " + v)
	if err != nil {
		ctx.err = err
	}
	ctx.response.Answer = append(ctx.response.Answer, rr)

}

func (ctx *Context) Addr() string {
	addr := ctx.w.RemoteAddr().String()
	return strings.Split(addr, ":")[0]
}

func (ctx *Context) Port() int {
	addr := ctx.w.RemoteAddr().String()
	v, _ := strconv.Atoi(strings.Split(addr, ":")[1])
	return v
}

func (ctx *Context) region() string {
	if ctx.Region == nil {
		return ""
	}

	info, err := ctx.Region.Search(ctx.Addr())
	if err != nil {
		logger.Infof("dns got region error %v", err)
		return ""
	}

	return auxlib.B2S(info.Byte())

}

func (ctx *Context) QType() string {
	q := ctx.Q()

	if t, ok := dns.TypeToString[q.Qtype]; ok {
		return t
	}
	return ""
}

func (ctx *Context) QName() string {
	return ctx.Q().Name
}

func (ctx *Context) RR(addr string) (dns.RR, error) {
	t := ctx.QType()
	if t == "" {
		return nil, fmt.Errorf("invalid type")
	}

	if t == "CNAME" {
		return dns.NewRR(ctx.QName() + " IN CNAME " + addr + ".")
	}

	return dns.NewRR(ctx.QName() + " IN " + t + " " + addr)
}

