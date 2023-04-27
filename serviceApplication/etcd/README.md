# ETCD
1. dats_dir /usr/local/soft/etcd/data.etcd 必须以.etcd结尾
2. clent3 需要配置 
>  wget https://pkg.cfssl.org/R1.2/cfssljson_linux-amd64  wget https://pkg.cfssl.org/R1.2/cfssl_linux-amd64
>  mv cfssl_linux-amd64 /usr/bin/cfssl mv cfssljson_linux-amd64 /usr/bin/cfssljson　chmod +x /usr/bin/{cfssl,cfssljson}
>  建立一个目录参数 证书  etcd-ca     http://play.etcd.io/install

>  etcd-root-ca-csr.json
```json
{
  "key": {
    "algo": "rsa",
    "size": 2048
  },
  "names": [
    {
      "O": "etcd",
      "OU": "etcd Security",
      "L": "San Francisco",
      "ST": "California",
      "C": "USA"
    }
  ],
  "CN": "etcd-root-ca"
}
```
> cfssl gencert --initca=true /usr/local/soft/etcd/etcd-ca/etcd-root-ca-csr.json | cfssljson --bare /usr/local/soft/etcd/etcd-ca/etcd-root-c
```shell script
etcd-root-ca-csr.json  
etcd-root-c-key.pem    
etcd-root-c.csr        
etcd-root-c.pem  
```
> 验证 openssl x509 -in /usr/local/soft/etcd/etcd-ca/etcd-root-c.pem -text -noout

>  ca-csr.json
```json
{
    "CN": "WAE Etcd CA",
    "key": {
        "algo": "rsa",
        "size": 2048
    },
    "names": [
        {
            "C": "CN",
            "L": "SZ",
            "ST": "ZRothschild"
        }
    ]
}
```

>  服务端  etcd-gencert.json

```json
{
  "signing": {
    "default": {
        "usages": [
          "signing",
          "key encipherment",
          "server auth",
          "client auth"
        ],
        "expiry": "87600h"
    }
  }
}
```
>  s1-ca-csr.json

```json
{
  "key": {
    "algo": "rsa",
    "size": 2048
  },
  "names": [
    {
      "O": "etcd",
      "OU": "etcd Security",
      "L": "San Francisco",
      "ST": "California",
      "C": "USA"
    }
  ],
  "CN": "s1",
  "hosts": [
    "127.0.0.1",
    "localhost"
  ]
}
```

```shell script
cfssl gencert \
  --ca /usr/local/soft/etcd/etcd-ca/etcd-root-c.pem \
  --ca-key /usr/local/soft/etcd/etcd-ca/etcd-root-c-key.pem \
  --config /usr/local/soft/etcd/etcd-ca/etcd-gencert.json \
  /usr/local/soft/etcd/etcd-ca/s1-ca-csr.json | cfssljson --bare /usr/local/soft/etcd/etcd-ca/s1
  
/usr/local/soft/etcd/etcd-ca/s1-ca-csr.json
/usr/local/soft/etcd/etcd-ca/s1.csr
/usr/local/soft/etcd/etcd-ca/s1-key.pem
/usr/local/soft/etcd/etcd-ca/s1.pem
```

cfssl gencert \
  --ca /usr/local/soft/etcd/etcd-ca/etcd-root-c.pem \
  --ca-key /usr/local/soft/etcd/etcd-ca/etcd-root-c-key.pem \
  --config /usr/local/soft/etcd/etcd-ca/etcd-gencert.json \
  /usr/local/soft/etcd/etcd-ca/s3-ca-csr.json | cfssljson --bare /usr/local/soft/etcd/etcd-ca/s3

```shell
[Unit]
Description=etcd
Documentation=https://github.com/coreos/etcd
Conflicts=etcd.service
Conflicts=etcd2.service

[Service]
Type=notify
Restart=always
RestartSec=5s
LimitNOFILE=40000
TimeoutStartSec=0

ExecStart=/usr/local/soft/etcd/bin/etcd --name s3 \
  --data-dir /usr/local/soft/etcd/s3 \
  --listen-client-urls https://localhost:32379 \
  --advertise-client-urls https://localhost:32379 \
  --listen-peer-urls https://localhost:32380 \
  --initial-advertise-peer-urls https://localhost:32380 \
  --initial-cluster s1=https://localhost:2380,s2=https://localhost:22380,s3=https://localhost:32380 \
  --initial-cluster-token tkn \
  --initial-cluster-state new \
  --client-cert-auth \
  --trusted-ca-file /usr/local/soft/etcd/etcd-ca/etcd-root-c.pem \
  --cert-file /usr/local/soft/etcd/etcd-ca/s3.pem \
  --key-file /usr/local/soft/etcd/etcd-ca/s3-key.pem \
  --peer-client-cert-auth \
  --peer-trusted-ca-file /usr/local/soft/etcd/etcd-ca/etcd-root-c.pem \
  --peer-cert-file /usr/local/soft/etcd/etcd-ca/s3.pem \
  --peer-key-file /usr/local/soft/etcd/etcd-ca/s3-key.pem \

[Install]
WantedBy=multi-user.target
```

> sudo mv /tmp/s1.service /etc/systemd/system/s1.service
```bash
# to start service
sudo systemctl daemon-reload
sudo systemctl cat s1.service
sudo systemctl enable s1.service
sudo systemctl start s1.service

# to get logs from service
sudo systemctl status s1.service -l --no-pager
sudo journalctl -u s1.service -l --no-pager|less
sudo journalctl -f -u s1.service

# to stop service
sudo systemctl stop s1.service
sudo systemctl disable s1.service


ETCDCTL_API=3 /usr/local/soft/etcd/bin/etcdctl \
  --endpoints localhost:2379,localhost:22379,localhost:32379 \
  --cacert /usr/local/soft/etcd/etcd-ca/etcd-root-c.pem \
  --cert /usr/local/soft/etcd/etcd-ca/s1.pem \
  --key /usr/local/soft/etcd/etcd-ca/s1-key.pem \
  endpoint health

ETCDCTL_API=3 etcdctl put foo bar --lease=1234
```


ENDPOINTS="127.0.0.1:2379,127.0.0.1:22379,127.0.0.1:32379"

value("user1") = "bad"
2. nohup etcd --config-file etcd.conf > nohup.log 2>&1 &
3. 状态检测 
> ./bin/etcd
> etcdctl member list 
> etcdctl cluster-health
> etcdctl --endpoints=127.0.0.1:2379 put key value

etcdctl --endpoints=127.0.0.1:2379  --cacert /usr/local/soft/etcd/etcd-ca/etcd-root-c.pem --cert /usr/local/soft/etcd/etcd-ca/s1.pem --key /usr/local/soft/etcd/etcd-ca/s1-key.pem put key value

```shell script

goreman -f Procfile start

export ETCDCTL_API=3
HOST=127.0.0.1
ENDPOINTS=$HOST:2379,$HOST:22379,$HOST:32379

export ETCDCTL_DIAL_TIMEOUT=3s
export ETCDCTL_CACERT=/usr/local/soft/etcd/etcd-ca/ca.pem
export ETCDCTL_CERT=/usr/local/soft/etcd/etcd-ca/client.pem
export ETCDCTL_KEY=/usr/local/soft/etcd/etcd-ca/client-key.pem

export ETCDCTL_CA_FILE=/usr/local/soft/etcd/etcd-ca/ca.pem
export ETCDCTL_KEY_FILE=/usr/local/soft/etcd/etcd-ca/client-key.pem
export ETCDCTL_CERT_FILE=/usr/local/soft/etcd/etcd-ca/client.pem
```


### 快速安装

etcd docker run -p 2379:2379 -p 2380:2380 --mount type=bind,source=D:/soft/dockerData/data/etcd-data.tmp,destination=/etcd-data --name etcd-gcr-v3.5.8 gcr.io/etcd-development/etcd:v3.5.8 /usr/local/bin/etcd --name s1 --data-dir /etcd-data --listen-client-urls http://0.0.0.0:2379  --advertise-client-urls http://0.0.0.0:2379 --listen-peer-urls http://0.0.0.0:2380 --initial-advertise-peer-urls http://0.0.0.0:2380 --initial-cluster s1=http://0.0.0.0:2380 --initial-cluster-token tkn --initial-cluster-state new --log-level info --logger zap --log-outputs stderr