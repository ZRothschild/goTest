## openssl 国家 SM 系列加密

1. 下载安装最新的 [openssl地址  https://www.openssl.org/source/openssl-3.0.0-alpha16.tar.gz](https://www.openssl.org/source/)
2. 解压安装包 tar -zxvf openssl-3.0.0-alpha16.tar.gz -C /usr/local

> 在解压包里面执行

3. 清空之前的安装数据

```shell
sudo make clean
```  

4. 安装前的配置 自己查看配置

```shell
sudo ./config  --prefix=/usr/local
// 或者执行 
sudo ./config  share
```   

5. 安装执行

```shell
sudo make && make install
```  

6. 设置链接 libssl.so.3 libcrypto.so.3 可以先查看是否有这个文件

```shell
ln -s /usr/local/lib/libssl.so.3 /usr/lib/libssl.so.3
ln -s /usr/local/lib/libcrypto.so.3 /usr/lib/libcrypto.so.3
```  

7. 查看加密方法 openssl 简单命令

```shell
// 查看版本号
openssl versio

// 检查椭圆曲线是否包含SM2
openssl ecparam -list_curves | grep SM2

// 检查对称算法
openssl enc -ciphers

// 检查 SM3 哈希校验和
echo -n "abc" | ./apps/openssl dgst -SM3

```  
