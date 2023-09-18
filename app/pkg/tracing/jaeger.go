package tracing

import (
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"io"
)

type Config struct {
	ServiceName string `mapstructure:"serviceName"`
	HostPort    string `mapstructure:"hostPort"`
	Enable      bool   `mapstructure:"enable"`
	LogSpans    bool   `mapstructure:"logSpans"`
}

func NewJaegerTracer(jaegerConfig *Config) (opentracing.Tracer, io.Closer, error) {
	cfg := &config.Configuration{
		ServiceName: jaegerConfig.ServiceName,

		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},

		Reporter: &config.ReporterConfig{
			LogSpans:           jaegerConfig.LogSpans,
			LocalAgentHostPort: jaegerConfig.HostPort,
		},
	}

	return cfg.NewTracer(config.Logger(jaeger.StdLogger))
}
