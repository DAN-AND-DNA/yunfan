package a_service

import (
	"context"
	"io"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

type Ping_args struct {
	Carrier map[string]string
}

func (this *A_service) Ping(rwc io.ReadWriteCloser, args *Ping_args, reply *string) error {
	if args.Carrier != nil {
		wireContext, err := this.Tracer.Extract(opentracing.TextMap, opentracing.TextMapCarrier(args.Carrier))
		if err == nil {
			span := this.Tracer.StartSpan("Ping", ext.RPCServerOption(wireContext))
			defer span.Finish()
		} else {
			return err
		}
	}

	raw_msg, _ := endpoint_ping(context.Background(), args)
	*reply = raw_msg.(string)
	return nil
}

// XXX handle request
func endpoint_ping(c context.Context, req_or_err interface{}) (interface{}, error) {
	_ = c
	_ = req_or_err
	return "pong", nil
}
