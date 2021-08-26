# Version 1.0 docker run --restart=always -p 8083:8083 -p 8081:8081 -p 8085:8085/tcp -p 3002:3002 --name goPay  -v D:\work:/data  -d centos:goPay
FROM centos:latest

MAINTAINER zrothschild

LABEL version="1.0" description="goPay" by="zrothschild"

VOLUME ["/data"]

RUN mkdir /usr/local/soft && cd /usr/local/soft && yum install -y wget  && wget https://studygolang.com/dl/golang/go1.16.5.linux-amd64.tar.gz && tar -zxvf go1.16.5.linux-amd64.tar.gz && /usr/local/soft/go/bin/go env -w GOPROXY=https://goproxy.cn,direct && /usr/local/soft/go/bin/go install github.com/go-delve/delve/cmd/dlv@latest && yum group install -y 'Development Tools'

ENV PATH /usr/local/soft/go/bin:/root/go/bin:$PATH


EXPOSE 8083 8085 3002 8081