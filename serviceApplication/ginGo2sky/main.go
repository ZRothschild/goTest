package main

import (
	"fmt"
	"github.com/SkyAPM/go2sky"
	v3 "github.com/SkyAPM/go2sky-plugins/gin/v3"
	"github.com/SkyAPM/go2sky/reporter"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	agentv3 "skywalking.apache.org/repo/goapi/collect/language/agent/v3"
	"strconv"
	"time"
)

//func main() {
//	// Use gRPC reporter for production
//	re, err := reporter.NewGRPCReporter("127.0.0.1:11800")
//	if err != nil {
//		log.Fatalf("new reporter error %v \n", err)
//	}
//	defer re.Close()
//
//	hostname, _ := os.Hostname()
//	tracer, err := go2sky.NewTracer("gin-server", go2sky.WithInstance(hostname), go2sky.WithReporter(re))
//	if err != nil {
//		log.Fatalf("create tracer error %v \n", err)
//	}
//
//	gin.SetMode(gin.ReleaseMode)
//	r := gin.New()
//
//	//Use go2sky middleware with tracing
//	r.Use(v3.Middleware(r, tracer))
//
//	r.GET("/user/:name", func(c *gin.Context) {
//		name := c.Param("name")
//		c.String(200, "Hello %s", name)
//	})
//
//	go func() {
//		if err := http.ListenAndServe(":7070", r); err != nil {
//			panic(err)
//		}
//	}()
//	// Wait for the server to start
//	time.Sleep(time.Second)
//
//	wg := sync.WaitGroup{}
//	wg.Add(1)
//
//	go func() {
//		defer wg.Done()
//		request(tracer)
//	}()
//	wg.Wait()
//	// Output:
//}
//
//func request(tracer *go2sky.Tracer, _ ...h.ClientOption) {
//	//NewClient returns an HTTP Client with tracer
//	client, err := h.NewClient(tracer)
//	if err != nil {
//		log.Fatalf("create client error %v \n", err)
//	}
//
//	request, err := http.NewRequest("GET", fmt.Sprintf("%s/user/gin", "http://127.0.0.1:8080"), nil)
//	if err != nil {
//		log.Fatalf("unable to create http request: %+v\n", err)
//	}
//
//	res, err := client.Do(request)
//	if err != nil {
//		log.Fatalf("unable to do http request: %+v\n", err)
//	}
//
//	_ = res.Body.Close()
//}

func main() {
	// Use gRPC reporter for production
	rp, err := reporter.NewGRPCReporter(":11800")
	if err != nil {
		log.Fatalf("new reporter error %v \n", err)
	}
	defer rp.Close()

	hostname, _ := os.Hostname()
	tracer, err := go2sky.NewTracer("gin-server",go2sky.WithInstance(hostname), go2sky.WithReporter(rp))
	if err != nil {
		log.Fatalf("create tracer error %v \n", err)
	}

	gin.SetMode(gin.DebugMode)
	r := gin.New()
	//Use go2sky middleware with tracing
	r.Use(v3.Middleware(r, tracer))

	r.GET("/test", func(c *gin.Context) {
		span, ctx, err := tracer.CreateEntrySpan(c.Request.Context(), getOperationName(c), func(key string) (string, error) {
			return c.Request.Header.Get(key), nil
		})
		if err != nil {
			c.JSON(http.StatusOK, map[string]string{"aaa": "aaa"})
			return
		}
		span.SetOperationName("方法里面哦")
		span.SetComponent(5007)
		span.Tag(go2sky.TagHTTPMethod, c.Request.Method)
		span.Tag(go2sky.TagURL, c.Request.Host+c.Request.URL.Path)
		span.SetSpanLayer(agentv3.SpanLayer_Http)
		span.Log(time.Now(), "inner data","hellpw")
		c.Request = c.Request.WithContext(ctx)
		if len(c.Errors) > 0 {
			span.Error(time.Now(), c.Errors.String())
		}
		span.Tag(go2sky.TagStatusCode, strconv.Itoa(c.Writer.Status()))
		span.End()




		c.Request = c.Request.WithContext(ctx)
		c.JSON(http.StatusOK, map[string]string{"aaa": "aaa"})
	})

	// do something
	if err = r.Run(":7070"); err != nil {
		return
	}
}

func getOperationName(c *gin.Context) string {
	return fmt.Sprintf("/%s%s", c.Request.Method, c.FullPath())
}
