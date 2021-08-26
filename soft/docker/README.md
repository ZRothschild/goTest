### docker

#### [win10修改docker镜像的存储位置](https://blog.csdn.net/keyiis_sh/article/details/115273037?utm_medium=distribute.pc_relevant.none-task-blog-2%7Edefault%7EBlogCommendFromMachineLearnPai2%7Edefault-16.base&depth_1-utm_source=distribute.pc_relevant.none-task-blog-2%7Edefault%7EBlogCommendFromMachineLearnPai2%7Edefault-16.base)

将win10更新至最新，目前最新版本是20H2，在安装docker desktop(基于wsl2)后，直接下载镜像默认是安装在C盘的，如果C盘空间不足，可按如下方式转移至其它盘

docker启动后会自动运行下面两个镜像，docker-desktop是存放程序的，docker-desktop-data是存放镜像的，这两个wsl子系统都是默认放在系统盘的

1. 查看存在的服务

```shell

PS C:\Users\Administrator> wsl -l -v --all
  NAME                   STATE           VERSION
* docker-desktop         Running         2
  docker-desktop-data    Running         2
  
```

2、停止wsl

```shell

wsl --shutdown
  
```

3、停止wsl

```shell

wsl --export docker-desktop D:\docker\docker-desktop.tar
wsl --export docker-desktop-data D:\docker\docker-desktop-data.tar
  
```

4、取消注册

```shell

wsl --unregister docker-desktop
wsl --unregister docker-desktop-data
 
```

5、将导出的镜像导入到你要转移的目标位置

```shell

// 用于存储镜像目录 D:\docker\image 
wsl --import docker-desktop-data D:\docker\image D:\docker\docker-desktop-data.tar
// 运行docker主程序目录 D:\docker\bin 
wsl --import docker-desktop D:\docker\bin D:\docker\docker-desktop.tar
 
```

6、确认是否转移成功

```shell

PS C:\Users\Administrator> wsl -l -v --all
  NAME                   STATE           VERSION
* docker-desktop         Stopped         2
  docker-desktop-data    Stopped         2
 
```