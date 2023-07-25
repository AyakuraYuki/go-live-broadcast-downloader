package wtmicro

import (
	"context"
	"fmt"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpcruntime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"go-live-broadcast-downloader/plugins/alert/dingding"
	ay_const "go-live-broadcast-downloader/plugins/ay-const/trace"
	cjson "go-live-broadcast-downloader/plugins/json"
	"go-live-broadcast-downloader/plugins/log"
	"go-live-broadcast-downloader/plugins/metric"
	"go-live-broadcast-downloader/plugins/trace/grpc_trace"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
	"net/http"
	"runtime"
	"strings"
	"sync"
	"time"
)

const (
	panicRecover       = "PanicRecover"
	microURIFullMethod = "uri.FullMethod"
)

var (
	_serviceName     string
	_env             string
	_printRequestLog bool
	once             sync.Once
)

func LoadServiceName(serviceName string, env string, printRequestLog bool) {
	once.Do(func() {
		_serviceName = metricName(serviceName)
		_env = metricName(env)
		_printRequestLog = printRequestLog
	})
}

func GatewayMux() *grpcruntime.ServeMux {
	opts := []grpcruntime.ServeMuxOption{
		grpcruntime.WithMarshalerOption(grpcruntime.MIMEWildcard, &grpcruntime.HTTPBodyMarshaler{
			Marshaler: &grpcruntime.JSONPb{
				MarshalOptions: protojson.MarshalOptions{
					EmitUnpopulated: true,
				},
				UnmarshalOptions: protojson.UnmarshalOptions{
					DiscardUnknown: false,
				},
			},
		}),
	}
	return grpcruntime.NewServeMux(opts...)
}

func GrpcDialOpts() []grpc.DialOption {
	credentials := insecure.NewCredentials()
	return []grpc.DialOption{grpc.WithTransportCredentials(credentials)}
}

func GrpcServerOpts(traceFilterMethodNames ...string) []grpc.ServerOption {
	return []grpc.ServerOption{
		grpc.UnaryInterceptor(grpcmiddleware.ChainUnaryServer(
			grpc_trace.UnaryServerInterceptor(opentracingOptions(traceFilterMethodNames...)...), // 链路追踪
			grpcrecovery.UnaryServerInterceptor(recoveryInterceptor()),
			handlerInterceptor,
		),
		),
	}
}

func handlerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {

	traceId := ""
	if span := opentracing.SpanFromContext(ctx); span != nil {
		if sc, ok := span.Context().(jaeger.SpanContext); ok {
			traceId = sc.TraceID().String()
			ctx = context.WithValue(ctx, ay_const.TraceCtxKey, traceId)
		}
	}

	fullMethod := metricName(info.FullMethod)
	if fullMethod != "" {
		fullMethod = fullMethod[strings.LastIndex(fullMethod, "_")+1:]
	}
	st := time.Now()
	resp, err := handler(ctx, req)

	metric.HistogramVec.Timing(_serviceName, []string{"env", _env, "method", fullMethod, "ret", metric.RetLabel(err)}, st)
	return resp, err
}

func metricName(s string) string {
	s = strings.ReplaceAll(s, "/", "_")
	s = strings.ReplaceAll(s, "-", "_")
	return strings.ReplaceAll(s, ".", "_")
}

// 链路追踪
func opentracingOptions(traceFilterMethodNames ...string) []grpc_trace.Option {
	return []grpc_trace.Option{
		grpc_trace.WithTracer(opentracing.GlobalTracer()),
		grpc_trace.WithOpName(func(method string) string {
			return strings.ReplaceAll(method, "/dodo.go.pbgen.service.", "")
		}),
		grpc_trace.WithFilterFunc(func(ctx context.Context, fullMethodName string) bool {
			// 过滤不需要投递的trace
			for _, name := range traceFilterMethodNames {
				if strings.HasSuffix(strings.ToLower(fullMethodName), "/"+strings.ToLower(name)) {
					return false
				}
			}
			return true
		}),
		grpc_trace.WithUnaryRequestHandlerFunc(func(span opentracing.Span, method string, begin time.Time, req, rsp interface{}, err error) {
			reqBs, _ := cjson.JSON.Marshal(req)
			respBs, _ := cjson.JSON.Marshal(rsp)
			traceId := ""
			if sc, ok := span.Context().(jaeger.SpanContext); ok {
				traceId = sc.TraceID().String()
			}
			span.LogKV("_request", string(reqBs), "_response", string(respBs))
			if _printRequestLog {
				log.Debug("handlerInterceptor").Msgf("traceId:[%s],uri:[%s],ut:[%s],  <<request:%s>>,  <<response:%s>>,  err:[%v]", traceId, method, time.Since(begin).String(), reqBs, respBs, err)
			}
		}),
	}
}

// panic recover
func recoveryInterceptor() grpcrecovery.Option {
	return grpcrecovery.WithRecoveryHandlerContext(func(ctx context.Context, p interface{}) (err error) {
		buf := make([]byte, 4<<10) //4k
		buf = buf[:runtime.Stack(buf, false)]
		str := fmt.Sprintf("%v", p)
		log.Error(panicRecover).Stack().Msgf("%s ______stack:___ %s", str, string(buf))

		dingding.AlertDingDing(&dingding.Message{
			Title:       panicRecover,
			Env:         _env,
			ServiceName: _serviceName,
			Text:        str + " ______stack:___ " + string(buf),
		})
		return err
	})
}

// GatewayHandlerFunc 根据请求头判断是 grpc 请求还是 grpc-gateway 请求
// grpc-http 共用同一个端口时使用
func GatewayHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHandler.ServeHTTP(w, r)
		}
	}), &http2.Server{})
}
