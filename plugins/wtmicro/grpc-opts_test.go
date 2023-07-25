package wtmicro

import (
	"testing"
)

func TestMixin(t *testing.T) {
	//ctx, cancel := context.WithCancel(context.Background())
	//defer cancel()
	//
	//conf := config.GlobalConfig
	//svr := service.NewHelloService(conf)
	//
	//gwMux := wtmicro.GatewayMux()
	//err := pb.RegisterGreeterHandlerFromEndpoint(ctx, gwMux, conf.GrpcPort, wtmicro.GrpcDialOpts())
	//if err != nil {
	//	log.Error(wtmicro.ServerStart).Msgf("%v", err)
	//	return
	//}
	//
	//rpcServer := grpc.NewServer(wtmicro.GrpcServerOpts()...)
	//defer func() {
	//	fmt.Println("aaaa")
	//	rpcServer.GracefulStop()
	//}()
	//
	//pb.RegisterGreeterServer(rpcServer, svr)
	//
	//hSvr := http.Server{
	//	Addr:    conf.GrpcPort,
	//	Handler: wtmicro.GatewayHandlerFunc(rpcServer, gwMux),
	//}
	//
	//ch := make(chan os.Signal, 1)
	//signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP, syscall.SIGQUIT)
	//go func() {
	//	s := <-ch
	//	fmt.Println("signal-down.Notify", time.Now(), s.String(), s)
	//	er := hSvr.Shutdown(ctx)
	//	fmt.Println(er)
	//	time.Sleep(time.Second * 4)
	//	fmt.Println("signal-down.Notify", time.Now())
	//
	//	//if i, ok := s.(syscall.Signal); ok {
	//	//	os.Exit(int(i))
	//	//} else {
	//	//	os.Exit(0)
	//	//}
	//}()
	//
	//err = hSvr.ListenAndServe()
	//if err != nil {
	//	log.Error(wtmicro.ServerStart).Msgf("%v", err)
	//}
}
