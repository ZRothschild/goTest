## k8s部署

### Linux内核版本
1. 查看内核版本
```
[zr@localhost ~]$ uname -r
5.14.0-390.el9.aarch64
```


### 修改主机名

> ***注意*** 每一个节点设置相应的名称

```shell

sudo hostnamectl set-hostname master01
sudo hostnamectl set-hostname master02
sudo hostnamectl set-hostname node01
sudo hostnamectl set-hostname node02
sudo hostnamectl set-hostname node03

```



### 添加hosts

```shell

sudo bash -c "cat >> /etc/hosts" << EOF
192.168.103.129   master01
192.168.103.130   master02

192.168.103.131   node01
192.168.103.132   node02
192.168.103.133   node03
EOF

```

EOF

# 配置etcd请求csr文件

```shell
bash -c "cat > ~/work/etcd/etcd-csr.json" << EOF
{
    "CN": "kubernetes",
    "hosts": [
        "10.0.0.1",
        "127.0.0.1",
        "192.168.103.121",
        "192.168.103.122",
        "192.168.103.123",
        "192.168.103.124",
        "192.168.103.125",
        "192.168.103.126",
        "192.168.103.127",
        "192.168.103.128",
        "192.168.103.129",
        "192.168.103.130",
        "192.168.103.131",
        "192.168.103.132",
        "192.168.103.133",
        "192.168.103.134",
        "192.168.103.135",
        "192.168.103.136",
        "192.168.103.137",
        "192.168.103.138",
        "192.168.103.139",
        "192.168.103.140",
        "192.168.103.141",
        "192.168.103.142",
        "192.168.103.143",
        "192.168.103.144",
        "192.168.103.145",
        "192.168.103.146",
        "192.168.103.147",
        "192.168.103.148",
        "192.168.103.149",
        "192.168.103.150",
        "192.168.103.151",
        "192.168.103.152",
        "192.168.103.153",
        "192.168.103.154",
        "192.168.103.155",
        "192.168.103.156",
        "192.168.103.157",
        "192.168.103.158",
        "192.168.103.159",
        "192.168.103.160",
        "kubernetes",
        "kubernetes.default",
        "kubernetes.default.svc",
        "kubernetes.default.svc.cluster",
        "kubernetes.default.svc.cluster.local"
    ],
    "key": {
        "algo": "rsa",
        "size": 2048
    },
    "names": [
        {
            "C": "CN",
            "L": "BeiJing",
            "ST": "BeiJing"
        }
    ]
}
EOF
```

```shell

bash -c "cat > ~/work/k8s/kube-apiserver-csr.json" << EOF
{
    "CN": "kubernetes",
    "hosts": [
        "10.0.0.1",
        "127.0.0.1",
        "192.168.103.121",
        "192.168.103.122",
        "192.168.103.123",
        "192.168.103.124",
        "192.168.103.125",
        "192.168.103.126",
        "192.168.103.127",
        "192.168.103.128",
        "192.168.103.129",
        "192.168.103.130",
        "192.168.103.131",
        "192.168.103.132",
        "192.168.103.133",
        "192.168.103.134",
        "192.168.103.135",
        "192.168.103.136",
        "192.168.103.137",
        "192.168.103.138",
        "192.168.103.139",
        "192.168.103.140",
        "192.168.103.141",
        "192.168.103.142",
        "192.168.103.143",
        "192.168.103.144",
        "192.168.103.145",
        "192.168.103.146",
        "192.168.103.147",
        "192.168.103.148",
        "192.168.103.149",
        "192.168.103.150",
        "192.168.103.151",
        "192.168.103.152",
        "192.168.103.153",
        "192.168.103.154",
        "192.168.103.155",
        "192.168.103.156",
        "192.168.103.157",
        "192.168.103.158",
        "192.168.103.159",
        "192.168.103.160",
        "kubernetes",
        "kubernetes.default",
        "kubernetes.default.svc",
        "kubernetes.default.svc.cluster",
        "kubernetes.default.svc.cluster.local"
    ],
    "key": {
        "algo": "rsa",
        "size": 2048
    },
    "names": [
        {
            "C": "CN",
            "L": "BeiJing",
            "ST": "BeiJing",
            "O": "k8s",
            "OU": "system"
        }
    ]
}
EOF
```

```shell

sudo bash -c "cat > /usr/local/k8s/cfg/etcd.conf" << EOF
#[Member]
ETCD_NAME="cmaster01"
ETCD_DATA_DIR="/var/lib/etcd/default.etcd"
ETCD_LISTEN_PEER_URLS="https://192.168.103.129:2380"
ETCD_LISTEN_CLIENT_URLS="https://192.168.103.129:2379"

#[Clustering]
ETCD_INITIAL_ADVERTISE_PEER_URLS="https://192.168.103.129:2380"
ETCD_ADVERTISE_CLIENT_URLS="https://192.168.103.129:2379"
ETCD_INITIAL_CLUSTER="cnode03=https://192.168.103.133:2380,cnode01=https://192.168.103.131:2380,cmaster01=https://192.168.103.129:2380"
ETCD_INITIAL_CLUSTER_TOKEN="etcd-cluster"
ETCD_INITIAL_CLUSTER_STATE="new"
EOF


sudo bash -c "cat > /usr/local/k8s/cfg/etcd.conf" << EOF
#[Member]
ETCD_NAME="cnode03"
ETCD_DATA_DIR="/var/lib/etcd/default.etcd"
ETCD_LISTEN_PEER_URLS="https://192.168.103.133:2380"
ETCD_LISTEN_CLIENT_URLS="https://192.168.103.133:2379"

#[Clustering]
ETCD_INITIAL_ADVERTISE_PEER_URLS="https://192.168.103.133:2380"
ETCD_ADVERTISE_CLIENT_URLS="https://192.168.103.133:2379"
ETCD_INITIAL_CLUSTER="cnode03=https://192.168.103.133:2380,cnode01=https://192.168.103.131:2380,cmaster01=https://192.168.103.129:2380"
ETCD_INITIAL_CLUSTER_TOKEN="etcd-cluster"
ETCD_INITIAL_CLUSTER_STATE="new"
EOF



sudo bash -c "cat > /usr/local/k8s/cfg/etcd.conf" << EOF
#[Member]
ETCD_NAME="cnode01"
ETCD_DATA_DIR="/var/lib/etcd/default.etcd"
ETCD_LISTEN_PEER_URLS="https://192.168.103.131:2380"
ETCD_LISTEN_CLIENT_URLS="https://192.168.103.131:2379"

#[Clustering]
ETCD_INITIAL_ADVERTISE_PEER_URLS="https://192.168.103.131:2380"
ETCD_ADVERTISE_CLIENT_URLS="https://192.168.103.131:2379"
ETCD_INITIAL_CLUSTER="cnode03=https://192.168.103.133:2380,cnode01=https://192.168.103.131:2380,cmaster01=https://192.168.103.129:2380"
ETCD_INITIAL_CLUSTER_TOKEN="etcd-cluster"
ETCD_INITIAL_CLUSTER_STATE="new"
EOF
```