# docker

```sh

docker pull imageid  # 拉取镜像

docker run -it ubuntu /bin/bash ## 进入容器

docker ps -a # 查看所有容器

docker start container_id ## 启动镜像

docker stop container_id ## 停止镜像

docker restart container_id ## 重启镜像 

# d 表示后台运行
docker run -itd --name ubuntu-test ubuntu /bin/bash 


# 后台运行的镜像，需要进入容器，如果从这里退出容器，会导致容器退出运行
docker attach container_id 

# 进入容器，相当于重新开启一个终端，从这里退出容器，容器不会停止
docker exec -it container_id /bin/bash 

# 导出本地某个容器
docker export container_id > container_name.tar 

# 导入容器
docker import container_name.tar test/ubuntu:v1
cat container_name.tar | docker import - test/ubuntu:v1

# 删除容器
docker rm -f container_id 

# 载入镜像 
docker pull training/webapp 

docker run -d -P traning/webapp python app.py
-d 后台运行
-P 将内部使用的网络端口随机映射到我们使用的主机上

docker run -d -p 5000:4999 training/webapp python app.py
# 本机上5000端口映射到容器里的4999

docker ps 
docker port container_id 

# 查看应用日志
docker logs -f container_id

docker tp container_id


```