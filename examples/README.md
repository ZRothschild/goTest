# gRPC Hello World

Follow these setup to run the [quick start][] example:

1. Get the code:

   ```console
   $ go get google.golang.org/grpc/examples/helloworld/greeter_client
   $ go get google.golang.org/grpc/examples/helloworld/greeter_server
   ```

2. Run the server:

   ```console
   $ $(go env GOPATH)/bin/greeter_server &
   ```

3. Run the client:

   ```console
   $ $(go env GOPATH)/bin/greeter_client
   Greeting: Hello world
   ```

For more details (including instructions for making a small change to the example code) or if you're having trouble
running this example, see [Quick Start][].

[quick start]: https://grpc.io/docs/languages/go/quickstart

### grpc-go 源码阅读

1. 监听`tcp`端口

```go

lis, err := net.Listen("tcp", ":5051")

```

2. 启动grpc服务，这个服务没有注册服务，没有接手任务请求.

```go

s := grpc.NewServer()

```

> 1. 服务接手参数`ServerOption`接口方法`apply(*serverOptions)`，为服务配置

```go

type serverOptions struct {
    creds                 credentials.TransportCredentials
    codec                 baseCodec
    cp                    Compressor
    dc                    Decompressor
    unaryInt              UnaryServerInterceptor
    streamInt             StreamServerInterceptor
    chainUnaryInts        []UnaryServerInterceptor
    chainStreamInts       []StreamServerInterceptor
    inTapHandle           tap.ServerInHandle
    statsHandler          stats.Handler
    maxConcurrentStreams  uint32
    maxReceiveMessageSize int
    maxSendMessageSize    int
    unknownStreamDesc     *StreamDesc
    keepaliveParams       keepalive.ServerParameters
    keepalivePolicy       keepalive.EnforcementPolicy
    initialWindowSize     int32
    initialConnWindowSize int32
    writeBufferSize       int
    readBufferSize        int
    connectionTimeout     time.Duration
    maxHeaderListSize     *uint32
    headerTableSize       *uint32
    numServerWorkers      uint32
}
```

> 2. `NewServer`主要实现

配置设置

```go

	opts := defaultServerOptions
	for _, o := range opt {
		o.apply(&opts)
	}
	
```

服务结构体初始化

```go

	s := &Server{
        lis:      make(map[net.Listener]bool),
        opts:     opts,
        conns:    make(map[string]map[transport.ServerTransport]bool),
        services: make(map[string]*serviceInfo),
        quit:     grpcsync.NewEvent(),
        done:     grpcsync.NewEvent(),
        czData:   new(channelzData),
    }
	
```

处理执行服务拦截器

```go

 // 一元服务拦截器
 chainUnaryServerInterceptors(s)
 // 流式服务蓝机器
 chainStreamServerInterceptors(s)
	
```

是否启用追踪包`s.events`是一个log接口

```go

   s.cv = sync.NewCond(&s.mu)
   if EnableTracing {
    _, file, line, _ := runtime.Caller(1)
    s.events = trace.NewEventLog("grpc.Server", fmt.Sprintf("%s:%d", file, line))
   }
	
```

是否开启多服务工作者`s.serverWorkerChannels[i]` 在`func (s *Server) serveStreams(st transport.ServerTransport)` 被设置

```go

    if s.opts.numServerWorkers > 0 {
        s.initServerWorkers()
    }
	
```

数据收集是否开启 `s.channelzID` 后面会赋值给 `ls.channelzID`

```go

	if channelz.IsOn() {
		s.channelzID = channelz.RegisterServer(&channelzServer{s}, "")
	}

```

2. `pb.RegisterGreeterServer(s, &server{})`注册服务

> 1. 执行 `s.RegisterService(&_Greeter_serviceDesc, srv)`判断`srv`是否实现`HandlerType`
> 1.1 执行 `s.register(sd, ss)` `sd`指定的服务模板为`grpc.ServiceDesc`类型结构体，`ss`为自定义实现接口服务结构体
```go
var _Greeter_serviceDesc = grpc.ServiceDesc{
	ServiceName: "helloworld.Greeter",
	HandlerType: (*GreeterServer)(nil),// 处理类型
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SayHello",
			Handler:    _Greeter_SayHello_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "helloworld/helloworld.proto",
}
```
1. 开启锁，支持并发
2. 判断是否服务已启动，才注册服务
3. 判断是否服务名重复
4. 使用`sd`, `ss`组装`serviceInfo`
```go

	info := &serviceInfo{
		serviceImpl: ss, // 服务结构体
		methods:     make(map[string]*MethodDesc), // 存放一元方法
		streams:     make(map[string]*StreamDesc),  // 存放流式方法
		mdata:       sd.Metadata, // 存放
	}
	for i := range sd.Methods {
		d := &sd.Methods[i]
		info.methods[d.MethodName] = d
	}
	for i := range sd.Streams {
		d := &sd.Streams[i]
		info.streams[d.StreamName] = d
	}
	s.services[sd.ServiceName] = info

```

3. `s.Serve(lis)` 启动服务

> 1. 设置服务开启，开启锁 
```go

s.serve = true
s.serveWG.Add(1)

```

>  2. 主装自定义`listenSocket`,设置`ls.channelzID`  打印行为 `channelz.Warningf(logger, s.channelzID, "grpc: Server.Serve failed to complete security handshake from %q: %v", rawConn.RemoteAddr(), err)`

```go

ls := &listenSocket{Listener: lis}
s.lis[ls] = true // 设置链接以开启

if channelz.IsOn() {
    ls.channelzID = channelz.RegisterListenSocket(ls, s.channelzID, lis.Addr().String())
}

```

>  3.  循环接听客户端请求
1. 判断是否错误，尝试多次链接
```go
    rawConn, err := lis.Accept()
	s.serveWG.Add(1)
	go func() {
		s.handleRawConn(lis.Addr().String(), rawConn)
		s.serveWG.Done()
	}()

```

详解  `s.handleRawConn(lis.Addr().String(), rawConn)`

1. 设置过期时间，服务器执行身份验证握手,握手这个有点复杂哦 返回的 `conn`它是一个原始包的接口`net.Conn`, `authInfo`

```go

	rawConn.SetDeadline(time.Now().Add(s.opts.connectionTimeout))
	conn, authInfo, err := s.useTransportAuthenticator(rawConn)
	
```

2. http2 协议转换

```go

 st := s.newHTTP2Transport(conn, authInfo)
 // 组装服务端传输层配置
 config := &transport.ServerConfig{
    MaxStreams:            s.opts.maxConcurrentStreams,
    AuthInfo:              authInfo,
    InTapHandle:           s.opts.inTapHandle,
    StatsHandler:          s.opts.statsHandler,
    KeepaliveParams:       s.opts.keepaliveParams,
    KeepalivePolicy:       s.opts.keepalivePolicy,
    InitialWindowSize:     s.opts.initialWindowSize,
    InitialConnWindowSize: s.opts.initialConnWindowSize,
    WriteBufferSize:       s.opts.writeBufferSize,
    ReadBufferSize:        s.opts.readBufferSize,
    ChannelzParentID:      s.channelzID,
    MaxHeaderListSize:     s.opts.maxHeaderListSize,
    HeaderTableSize:       s.opts.headerTableSize,
 }
 st, err := transport.NewServerTransport("http2", c, config)
 
 // 其实里面是执行了 
 func newHTTP2Server(conn net.Conn, config *ServerConfig) (_ ServerTransport, err error)

 // 组装请求如返回  r，w, framer 是 http2 包Framer 与 w 的组装
 func newFramer(conn net.Conn, writeBufferSize, readBufferSize int, maxHeaderListSize uint32) *framer
 
```

3. `s.addConn(lisAddr, st)` 再次判断是否已经有这个链接

```go

	if !s.addConn(lisAddr, st) {
		return
	}
	go func() {
		s.serveStreams(st)
		s.removeConn(lisAddr, st)
	}()

```
`func (s *Server) serveStreams(st transport.ServerTransport)`里面有`st.HandleStreams`


`st.HandleStreams`在`http2_server.go`469行，之后很多数据都交给`st transport.ServerTransport`处理