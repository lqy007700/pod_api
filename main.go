package main

import (
	"github.com/afex/hystrix-go/hystrix"
	"github.com/asim/go-micro/plugins/registry/consul/v3"
	ratelimit "github.com/asim/go-micro/plugins/wrapper/ratelimiter/uber/v3"
	"github.com/asim/go-micro/plugins/wrapper/select/roundrobin/v3"
	opentracing2 "github.com/asim/go-micro/plugins/wrapper/trace/opentracing/v3"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/registry"
	"github.com/asim/go-micro/v3/server"
	"github.com/opentracing/opentracing-go"
	"github.com/zxnlx/common"
	"github.com/zxnlx/pod/proto/pod"
	handler2 "github.com/zxnlx/pod_api/handler"
	hystrix2 "github.com/zxnlx/pod_api/plugin/hystrix"
	"github.com/zxnlx/pod_api/proto/pod_api"
	"net"
	"net/http"

	"strconv"
)

var (
	serviceHost = "192.168.55.2"
	servicePort = "8081"

	// 注册中心配置
	consulHost       = serviceHost
	consulPort int64 = 8500

	// 链路
	tracerHost = "192.168.55.2"
	tracerPort = 6381

	// 熔断
	hystrixPort = 9092

	//监控
	prometheusPort = 9192
)

// 注册中心
func initRegistry() registry.Registry {
	return consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			consulHost + ":" + strconv.FormatInt(consulPort, 10),
		}
	})
}

func initTracer() {
	// 链路追踪
	// jaeger
	tracer, i, err := common.NewTracer("base", tracerHost+":"+strconv.Itoa(tracerPort))
	if err != nil {
		common.Fatal(err)
		return
	}
	defer i.Close()
	opentracing.SetGlobalTracer(tracer)
}

func initHystrix() {
	// 熔断
	hystrixHandler := hystrix.NewStreamHandler()
	hystrixHandler.Start()
	go func() {
		//http://192.168.0.112:9092/turbine/turbine.stream
		//看板访问地址 http://127.0.0.1:9002/hystrix，url后面一定要带 /hystrix
		err := http.ListenAndServe(net.JoinHostPort("0.0.0.0", strconv.Itoa(hystrixPort)), hystrixHandler)
		if err != nil {
			common.Fatal(err)
		}
	}()
}

func main() {
	c := initRegistry()
	initTracer()
	initHystrix()

	common.PrometheusBoot(prometheusPort)

	service := micro.NewService(
		micro.Server(server.NewServer(func(options *server.Options) {
			options.Advertise = serviceHost + ":" + servicePort
		})),
		micro.Name("go.micro.api.pod"),
		micro.Version("latest"),
		micro.Address(":"+servicePort),

		micro.Registry(c),
		// 链路
		micro.WrapHandler(opentracing2.NewHandlerWrapper(opentracing.GlobalTracer())),
		micro.WrapClient(opentracing2.NewClientWrapper(opentracing.GlobalTracer())),
		// 熔断
		micro.WrapClient(hystrix2.NewClientHystrixWrapper()),
		// 限流
		micro.WrapHandler(ratelimit.NewHandlerWrapper(1000)),

		// 负载均衡
		micro.WrapClient(roundrobin.NewClientWrapper()),
	)

	service.Init()

	podService := pod.NewPodService("go.micro.service.pod", service.Client())

	err := pod_api.RegisterPodApiHandler(service.Server(), &handler2.PodApi{
		PodService: podService,
	})
	if err != nil {
		return
	}

	err = service.Run()
	if err != nil {
		common.Fatal(err)
		return
	}
}
