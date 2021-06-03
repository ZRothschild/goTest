package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io"
	"log"
	"net/http"
	"time"
)

// WebRequestTotal 初始化 web_reqeust_total， counter类型指标， 表示接收http请求总次数
var WebRequestTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "aa_reqeust_total",
		Help: "Number of hello requests in total",
	},
	// 设置两个标签 请求方法和 路径 对请求总次数在两个
	[]string{"method", "endpoint"},
)

// WebRequestDuration web_request_duration_seconds，Histogram类型指标，bucket代表duration的分布区间
var WebRequestDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "aa_request_duration_seconds",
		Help:    "web request duration distribution",
		Buckets: []float64{0.1, 0.3, 0.5, 0.7, 0.9, 1},
	},
	[]string{"method", "endpoint"},
)

func init() {
	// 注册监控指标
	err := prometheus.Register(WebRequestTotal)
	if err != nil {
		fmt.Println("========", err)
		return
	}
	err = prometheus.Register(WebRequestDuration)
	if err != nil {
		fmt.Println("========", err)
		return
	}
}

// https://github.com/prometheus/prometheus
// docker run --name prometheus -d -p  9090:9090 -v D:\work\goTest\serviceApplication\prometheus:/etc/prometheus  prom/prometheus
// docker run --name prometheus -d -p  9090:9090 prom/prometheus
// http://localhost:9090/
// /etc/prometheus
// /metrics 接口用于数据监听或展示
// grafana 数据管理工具
func main() {
	// expose prometheus metrics接口
	var t = Test{
		CounterVec:   WebRequestTotal,
		HistogramVec: WebRequestDuration,
	}
	server := http.NewServeMux()                    // create a new mux server
	server.Handle("/metrics", promhttp.Handler())   // register a new handler for the /metrics endpoint
	server.Handle("/add", t)                        // register a new handler for the /metrics endpoint
	log.Fatal(http.ListenAndServe(":8756", server)) // start an http server using the mux server
}

type Test struct {
	*prometheus.CounterVec
	*prometheus.HistogramVec
}

func (t Test) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	t.CounterVec.With(prometheus.Labels{"method": r.Method, "endpoint": r.URL.Path}).Inc()
	_, err := io.WriteString(w, "大家好呀")
	//_, err := w.Write([]byte("大家好呀"))
	if err != nil {
		return
	}
	t.HistogramVec.With(prometheus.Labels{"method": r.Method, "endpoint": r.URL.Path}).Observe(time.Since(start).Seconds())
}
