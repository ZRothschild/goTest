package main

//func main() {
//	s := rpc.NewServer()
//	mux := http.NewServeMux()
//	s.RegisterCodec(json.NewCodec(), "application/json")
//	s.RegisterService(new(HelloService), "")
//	mux.Handle("/rpc", s)
//	http.ListenAndServe(":1234", mux)
//}
//
//type HelloArgs struct {
//	Who string
//}
//
//type HelloReply struct {
//	Message string
//}
//
//type HelloService struct{}
//
//func (h *HelloService) Say(r *http.Request, args *HelloArgs, reply *HelloReply) error {
//	reply.Message = "Hello, " + args.Who + "!"
//	return nil
//}
