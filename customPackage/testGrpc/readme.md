protoc --go_out=plugins=grpc:. helloworld.proto  这里用了plugins选项，提供对grpc的支持，否则不会生成Service的接口，方便编写服务器和客户端程序

protoc -I . helloworld.proto --go_out=plugins=grpc:.

protoc -I helloworld/ helloworld/helloworld.proto --go_out=plugins=grpc:helloworld


protoc --plugin=protoc-gen-go=C:\Users\zr\go\bin\protoc-gen-go.exe  --go_out=.  mes/message.proto


protoc --go_out=./name  --go_opt=paths=source_relative mes/message.proto


protoc --go_out=. --go_opt=paths=import mes/message.proto

