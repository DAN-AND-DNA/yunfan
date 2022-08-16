package rpc

import (
	"context"
	"io"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go/config"
)

var (
	default_rpc_tracer_mgr = New_rpc_tracer_mgr()
)

type Rpc_tracer_mgr struct {
	tracer        opentracing.Tracer
	tracer_closer io.Closer
	is_set_tracer bool
}

func New_rpc_tracer_mgr() *Rpc_tracer_mgr {
	return &Rpc_tracer_mgr{}
}

func (this *Rpc_tracer_mgr) Set_remote_tracer(service_name, hostname, port string) error {
	if this.is_set_tracer {
		return nil
	}

	if service_name == "" || hostname == "" || port == "" {
		this.is_set_tracer = false
		return nil
	}

	ptr_config := &config.Configuration{
		ServiceName: service_name,
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},

		Reporter: &config.ReporterConfig{
			LocalAgentHostPort: hostname + ":" + port,
			LogSpans:           false,
		},
	}

	var err error
	this.tracer, this.tracer_closer, err = ptr_config.NewTracer(config.Gen128Bit(true))
	if err != nil {
		this.is_set_tracer = false
		return err
	}

	opentracing.SetGlobalTracer(this.tracer)
	this.is_set_tracer = true
	return nil
}

func (this *Rpc_tracer_mgr) Shutdown() {
	if this == nil || !this.is_set_tracer || this.tracer_closer == nil {
		return
	}

	this.tracer_closer.Close()
}

func (this *Rpc_tracer_mgr) New_rpc_span() *Rpc_span {
	if this == nil || !this.is_set_tracer {
		return nil
	}

	new_rpc_track := &Rpc_span{
		mgr_tracer: this.tracer,
		Carrier:    map[string]string{},
	}

	return new_rpc_track
}

type Rpc_span struct {
	mgr_tracer opentracing.Tracer
	span       opentracing.Span
	Carrier    map[string]string
}

func (this *Rpc_span) Start(ctx context.Context, service_method string) (context.Context, error) {
	if this == nil || service_method == "" || this.mgr_tracer == nil {
		return ctx, nil
	}

	span, new_ctx := opentracing.StartSpanFromContext(ctx, service_method)
	this.span = span

	ext.SpanKindRPCClient.Set(this.span)
	if err := this.mgr_tracer.Inject(this.span.Context(), opentracing.TextMap, opentracing.TextMapCarrier(this.Carrier)); err != nil {
		return ctx, err
	}

	return new_ctx, nil
}

func (this *Rpc_span) Finish() {
	if this == nil || this.span == nil {
		return
	}

	this.span.Finish()
}

func Set_remote_tracer(service_name, hostname, port string) error {
	return default_rpc_tracer_mgr.Set_remote_tracer(service_name, hostname, port)
}

func New_rpc_span() *Rpc_span {
	return default_rpc_tracer_mgr.New_rpc_span()
}

func Shutdown_tracer() {
	default_rpc_tracer_mgr.Shutdown()
}
