package a_service

import (
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
)

type A_service struct {
	Tracer opentracing.Tracer
}

func New_a_service() *A_service {
	ptr_config := &config.Configuration{
		ServiceName: "v1.api.a_service",
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},

		Reporter: &config.ReporterConfig{
			LocalAgentHostPort: "worker01.dev.com:6831",
			LogSpans:           false,
		},
	}

	tracer, _, err := ptr_config.NewTracer(config.Gen128Bit(true))
	if err != nil {
		panic(err)
	}

	return &A_service{Tracer: tracer}
}
