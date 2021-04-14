protoc --go_out=plugins=grpc:. helloworld.proto  这里用了plugins选项，提供对grpc的支持，否则不会生成Service的接口，方便编写服务器和客户端程序

protoc -I . helloworld.proto --go_out=plugins=grpc:.

protoc -I helloworld/ helloworld/helloworld.proto --go_out=plugins=grpc:helloworld