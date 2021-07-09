package main
// docker run  -p 5775:5775/udp    -p 16686:16686    -p 6831:6831/udp   -p 6832:6832/udp   -p 5778:5778   -p 14268:14268    jaegertracing/all-in-one:latest
import (
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-lib/metrics"

	"github.com/uber/jaeger-client-go/config"
)

var (
	tracerServer opentracing.Tracer
)

func TraceInit(serviceName string, samplerType string, samplerParam float64) (opentracing.Tracer, io.Closer) {
	cfg := &config.Configuration{
		ServiceName: serviceName,
		Sampler: &config.SamplerConfig{
			Type:  samplerType,
			Param: samplerParam,
		},
		Reporter: &config.ReporterConfig{
			LocalAgentHostPort: "0.0.0.0:6831",
			LogSpans:           true,
		},
	}
	tracer, closer, err := cfg.NewTracer(config.Logger(jaeger.StdLogger), config.Metrics(metrics.NullFactory))
	if err != nil {
		panic(fmt.Sprintf("Init failed: %v\n", err))
	}

	return tracer, closer
}

func GetListProc(w http.ResponseWriter, req *http.Request) {
	spanCtx, _ := tracerServer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))
	span := tracerServer.StartSpan("GetListProc", ext.RPCServerOption(spanCtx))
	defer span.Finish()
	span.SetTag("zhao","赵钱孙")
	fmt.Println("Get request getList")
	respList := []string{"l1", "l2", "l3", "l4", "l5"}
	respString := ""

	for _, v := range respList {
		respString += v + ","
	}

	fmt.Println(respString)
	_, err := io.WriteString(w, respString)
	if err != nil {
		return
	}
}

func sendRequest(req *http.Request) {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Do send requst failed(%s)\n", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("ReadAll error(%s)\n", err)
		return
	}
	if resp.StatusCode != 200 {
		return
	}
	fmt.Printf("Response:%s\n", string(body))
}

func main() {
	// server
	var closerServer io.Closer
	tracerServer, closerServer = TraceInit("Trace-Server", "const", 1)
	defer closerServer.Close()


	http.HandleFunc("/getList", GetListProc)

	go http.ListenAndServe(":9909", nil)

	//client
	tracerClient, closerClient := TraceInit("CS-tracing", "const", 1)
	defer closerClient.Close()

	opentracing.SetGlobalTracer(tracerClient)

	span := tracerClient.StartSpan("getlist trace")
	span.SetTag("trace to", "getlist")
	defer span.Finish()

	reqURL := "http://localhost:9909/getList"
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = span.Tracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Header),
	)
	if err != nil {
		return
	}
	sendRequest(req)
}