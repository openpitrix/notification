# openpitrix All-in-One 模式
https://docs.openpitrix.io/v1.0/zh-CN/allinone/

第0步：
 wget -qO- https://get.docker.com/ | sh
 sudo service docker start
 docker run hello-world
  
 sudo curl -L "https://github.com/docker/compose/releases/download/1.22.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
 chmod +x /usr/local/bin/docker-compose
 docker-compose --version
 
 apt-get install git
 
 apt get make 
 

第一步: 准备环境
All-in-One 模式需要依赖以下软件:
第二步: 准备 OpenPitrix 源码文件
$ git clone https://github.com/openpitrix/openpitrix
第三步: 部署 OpenPitrix
$ cd openpitrix
$ make build
$ make compose-up
第四步: 验证

查看服务情况:
$ docker ps
查看界面:
访问 http://localhost:8000/ 查看 OpenPitrix 管理界面。

















