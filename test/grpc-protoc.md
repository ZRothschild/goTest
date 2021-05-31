## grpc 与 protoc

protoc -I . --go_out=plugins=grpc:. --go_opt=paths=source_relative mes/message.proto name/message.proto

protoc -I . --grpc-gateway_out=. --grpc-gateway_opt logtostderr=true --grpc-gateway_opt paths=source_relative
--grpc-gateway_opt generate_unbound_methods=true mes/message.proto name/message.proto
./third_party/googleapis/google/api/annotations.proto

1. 下载protoc 执行命令 https://github.com/protocolbuffers/protobuf 下载对应的os版本解压，把bin目录加入环境变量，或者把protoc.exe 移入已经加入环境变量的目录
2. grpc-gateway 文件只能对于单独的与应用proto 文件对应，必须都是同一个目录
3. 应该找对 annotations.proto 文件或复制到工作目录 google/api/
4. grpc 使用http2 所以必须使用http2
5. 可以用多个 -I 可以寻找依赖 比如说内嵌 的proto 文件 或 annotations.proto


protoc -I ./proto -I .  --go-grpc_out=./pb   doubleoffline.proto


protoc -I ./proto -I . --go_out=./pb --go_opt=paths=source_relative  --go-grpc_out=./pb --go-grpc_opt=paths=source_relative 

protoc -I ./proto -I .   --go_out=plugins=grpc:./pb  --go_opt=paths=source_relative  common.proto doubleoffline.proto google/api/annotations.proto

protoc -I ./proto -I . --go_out=./pb --go_opt=paths=source_relative  --go-grpc_out=./pb --go-grpc_opt=paths=source_relative --grpc-gateway_out=./pb  --grpc-gateway_opt generate_unbound_methods=true --grpc-gateway_opt paths=source_relative  doubleoffline.proto common.proto