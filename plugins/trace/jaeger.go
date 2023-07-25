package trace

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-client-go/transport"
	"go-live-broadcast-downloader/plugins/log"
	"io"
	"runtime"
	"strings"
	"time"
)

const (
	_httpURLInternal = "" // 内网
	_httpURLPublic   = "" // 公网
)

// NewJaegerTracer http trace
func NewJaegerTracer(serviceName, env string) (opentracing.Tracer, io.Closer, error) {
	cfg := jaegercfg.Configuration{
		ServiceName: fmt.Sprintf("%s-%s", env, serviceName),
		Tags:        []opentracing.Tag{{Key: "env", Value: env}},
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
		},
	}
	var (
		_httpUrl = _httpURLInternal
	)
	//_httpUrl = _httpURLPublic
	//if strings.ToLower(env) == "pro" || strings.ToLower(env) == "dev" {
	//_httpUrl = _httpURLInternal
	//}

	sender := transport.NewHTTPTransport(_httpUrl)
	reporter := jaeger.NewRemoteReporter(sender)
	tracer, closer, err := cfg.NewTracer(
		jaegercfg.Reporter(reporter),
	)
	if err != nil {
		log.Error("NewTracer").Msgf("%v", err)
	}
	if tracer != nil {
		opentracing.SetGlobalTracer(tracer)
	}

	return tracer, closer, err
}

// StartSpanFromContext 新建span
func StartSpanFromContext(ctx context.Context, opName ...string) (opentracing.Span, context.Context) {
	if len(opName) > 0 {
		return opentracing.StartSpanFromContext(ctx, opName[0])
	}

	pc, _, _, _ := runtime.Caller(1)
	if f := runtime.FuncForPC(pc); f != nil {
		if arr := strings.Split(f.Name(), "/"); len(arr) > 0 {
			return opentracing.StartSpanFromContext(ctx, arr[len(arr)-1])
		}
	}

	return opentracing.StartSpanFromContext(ctx, "unknow")
}

// StartSpanFromContextWithSt 新建span
func StartSpanFromContextWithSt(ctx context.Context, opName string, startTime time.Time) (opentracing.Span, context.Context) {
	ops := opentracing.StartTime(startTime)
	if len(opName) > 0 {
		return opentracing.StartSpanFromContext(ctx, opName, ops)
	}

	pc, _, _, _ := runtime.Caller(1)
	if f := runtime.FuncForPC(pc); f != nil {
		if arr := strings.Split(f.Name(), "/"); len(arr) > 0 {
			return opentracing.StartSpanFromContext(ctx, arr[len(arr)-1], ops)
		}
	}

	return opentracing.StartSpanFromContext(ctx, "unknow", ops)
}
