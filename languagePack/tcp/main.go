package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	//conn, err := net.Dial("tcp", "baidu.com:80")
	//if err != nil {
	//	// handle error
	//}
	//fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
	//status, err := bufio.NewReader(conn).ReadString('\n')
	//fmt.Println(status)

	la, err := net.ResolveUDPAddr("udp4", "127.0.0.1:0")
	if err != nil {
		fmt.Println(err)
	}
	c, err := net.ListenUDP("udp4", la)
	if err != nil {
		fmt.Println(err)
	}
	defer c.Close()
	c.LocalAddr()
	c.RemoteAddr()
	c.SetDeadline(time.Now().Add(10))
	c.SetReadDeadline(time.Now().Add(10))
	c.SetWriteDeadline(time.Now().Add(10))
	c.SetReadBuffer(2048)
	c.SetWriteBuffer(2048)



	wb := []byte("UDPCONN TEST")
	rb := make([]byte, 128)
	if _, err := c.WriteToUDP(wb, c.LocalAddr().(*net.UDPAddr)); err != nil {
		fmt.Println(err)
	}
	if _, _, err := c.ReadFromUDP(rb); err != nil {
		fmt.Println(err)
	}
	if _, _, err := c.WriteMsgUDP(wb, nil, c.LocalAddr().(*net.UDPAddr)); err != nil {
		fmt.Println(err)
	}
	if _, _, _, _, err := c.ReadMsgUDP(rb, nil); err != nil {
		fmt.Println(err)
	}

	if f, err := c.File(); err != nil {
		fmt.Println(err)
	} else {
		f.Close()
	}

	defer func() {
		if p := recover(); p != nil {
			fmt.Println(err)
		}
	}()

	c.WriteToUDP(wb, nil)
	c.WriteMsgUDP(wb, nil, nil)


}

--logtostderr=false \
--v=2 \
--log-dir=/opt/kubernetes/logs \
--etcd-servers=https://192.168.11.128:2379,https://192.168.11.129:2379,https://192.168.11.130:2379 \
--bind-address=0.0.0.0 \
--secure-port=6443 \
--advertise-address=192.168.11.128 \
--allow-privileged=true \
--service-cluster-ip-range=10.0.0.0/24 \
--enable-admission-plugins=NamespaceLifecycle,LimitRanger,ServiceAccount,ResourceQuota,NodeRestriction \
--authorization-mode=RBAC,Node \
--enable-bootstrap-token-auth=true \
--token-auth-file=/opt/kubernetes/cfg/token.csv \
--service-node-port-range=30000-32767 \
--kubelet-certificate-authority=/opt/kubernetes/ssl/ca.pem \
--kubelet-client-certificate=/opt/kubernetes/ssl/server.pem \
--kubelet-client-key=/opt/kubernetes/ssl/server-key.pem \
--tls-cert-file=/opt/kubernetes/ssl/server.pem  \
--tls-private-key-file=/opt/kubernetes/ssl/server-key.pem \
--client-ca-file=/opt/kubernetes/ssl/ca.pem \
--service-account-key-file=/opt/kubernetes/ssl/ca.pem \
--service-account-signing-key-file=/opt/kubernetes/ssl/ca-key.pem \
--service-account-issuer=https://kubernetes.default.svc.cluster.local \
--etcd-cafile=/opt/etcd/ssl/ca.pem \
--etcd-certfile=/opt/etcd/ssl/server.pem \
--etcd-keyfile=/opt/etcd/ssl/server-key.pem \
--audit-log-maxage=30 \
--audit-log-maxbackup=3 \
--audit-log-maxsize=100 \
--audit-log-path=/opt/kubernetes/logs/k8s-audit.log



--logtostderr=false \
--v=2 \
--log-dir=/opt/kubernetes/logs \
--leader-elect=true \
--master=https://192.168.11.128:6443 \
--bind-address=127.0.0.1 \
--allocate-node-cidrs=true \
--cluster-cidr=10.244.0.0/16 \
--client-ca-file=/opt/kubernetes/ssl/ca.pem \
--service-cluster-ip-range=10.0.0.0/24 \
--use-service-account-credentials=true \
--controllers=*,bootstrapsigner,tokencleaner \
--cluster-signing-cert-file=/opt/kubernetes/ssl/ca.pem \
--cluster-signing-key-file=/opt/kubernetes/ssl/ca-key.pem \
--root-ca-file=/opt/kubernetes/ssl/ca.pem \
--service-account-private-key-file=/opt/kubernetes/ssl/ca-key.pem




--logtostderr=false \\
--v=2 \\
--log-dir=/opt/kubernetes/logs \\
--leader-elect=true \\
--master=https://192.168.11.128:6443 \\
--bind-address=0.0.0.0


curl --cacert ./ssl/ca.pem --cert ./ssl/admin.pem --key ./ssl/admin-key.pem https://192.168.11.128:6443
