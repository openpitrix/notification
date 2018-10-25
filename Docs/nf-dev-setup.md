#ubuntu 配置grpc helloworld环境 
https://grpc.io/docs/quickstart/go.html#download-the-example
```
cd $GOPATH/src/google.golang.org/grpc/examples/helloworld
$ protoc -I helloworld/ helloworld/helloworld.proto --go_out=plugins=grpc:helloworld

$ go run greeter_server/main.go 
$ go run greeter_client/main.go
```
 


#Windows开发环境中配置grpc 
```
#https://studygolang.com/articles/13224?fr=sidebar  
cd /root/goworkspace/src/notification/pkg/pb
protoc -I ../pb --go_out=plugins=grpc:../pb ../pb/notification.proto
 
进入greeter_server下执行 
go run server.go

进入greeter_client下执行 
go run client.go

go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
 
 ```
#ubuntu下配置grpc，grpc-gateway 和swaggerUI

##0.0 ubuntu 配置grpc环境  
参考上一节
https://www.cnblogs.com/lienhua34/p/6285829.html  

##0.1 修改proto文件 
**增加对http的扩展配置**

```
option (google.api.http) = {
        post: "/v1/mail"
        body: "*"
    };
``` 

##1.notification.pb.go 
**notification.pb.go  这个go文件是grpc server服务需要的**

``` 
cd /root/goworkspace/src/notification/pkg/ 
 
protoc -I/usr/local/include -I. \
-I$GOPATH/src \
-I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
-I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway \
--go_out=plugins=grpc:. \
pb/notification.proto
```
>注意：github上给的文档的命令多带了一行参数Mgoogle...，导致出错。
no Go files in /root/goworkspace/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis/google/api

>--go_out=Mgoogle/api/annotations.proto=github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis/google/api,plugins=grpc:. \

##2.notification.pb.gw.go 

**notification.pb.gw.go 需要使用protoc生成gateway需要的go文件
这个文件就是gateway用来的协议文件，用来做grpc和http的协议转换**

``` 
protoc -I/usr/local/include -I. \
-I$GOPATH/src \
-I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
-I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway \
--grpc-gateway_out=logtostderr=true:. \
pb/notification.proto
 

```


##3.编写gateway代码 
gateway代码位于pkg/apigateway/apigateway.go
 
**首先echoEndpoint存储了需要连接的server信息，**  
**然后将这些信息和新建的server用gw.go中的RegisterGreeterHandlerFromEndpoint进行一个注册和绑定，** 
**这时底层就会连接echoEndpoint提供的远程server地址，这样gateway就作为客户端和远程server建立了连接，**
**之后用http启动新建的server，gateway就作为服务器端对外提供http的服务了。** 
**先启动api-server服务，再启动api-gateway，这时api-gateway连接上api-server后，对外建立http的监听。**

``` 
cd /root/goworkspace/src/notification/cmd/server
go run server_main.go
 
cd /root/goworkspace/src/notification/cmd/client
go run client_main.go

cd  /root/goworkspace/src/notification/cmd/gateway
go run gateway_main.go
 
curl -X POST "http://localhost:8080/v1/hello" -H "accept: application/json" -H "Cson" -d "{ \"name\": \"hello JoJo\"}"
curl -X POST -k http://localhost:8080/v1/hello -d "{ \"name\": \"hello JoJo\"}"
curl -X POST -k http://localhost:8080/v1/hello -d '{"name": "world"}'

http://192.168.0.3:8080/swagger.json

#使用浏览器访问http://localhost:8080/swagger.json

``` 

##4.集成swagger-ui  

###4.1.生成RESTful JSON API的Swagger说明 

**对应生成notification.swagger.json文件**   

```  
cd /root/goworkspace/src/notification/pkg/ 

protoc -I/usr/local/include -I. \
-I$GOPATH/src  -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
-I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway \
--swagger_out=logtostderr=true:. \
pb/notification.proto

```


###4.2.notification.swagger.go

**notification.swagger.go**  
在proto文件的目录下新建notification.swagger.go文件,定义常量为Swagger,常量值为notification.swagger.json内容
添加方法返回swagger.json的内容
 
###4.3 修改proxy.go代码

 
```
vim  /root/goworkspace/src/notification/pkg/apigateway/apigateway.go
     
func run() error {
    ctx := context.Background()
    ctx, cancel := context.WithCancel(ctx)
    defer cancel()

    mux := http.NewServeMux()
    mux.HandleFunc("/swagger.json", func(w http.ResponseWriter, req *http.Request) {
        io.Copy(w, strings.NewReader(gw.Swagger))
    })

    gwmux := runtime.NewServeMux()
    opts := []grpc.DialOption{grpc.WithInsecure()}
    err := gw.RegisterGreeterHandlerFromEndpoint(ctx, gwmux, *greeterEndpoint, opts)
    if err != nil {
        return err
    }

    log.Print("Greeter gRPC Server gateway start at port 8080...")
    http.ListenAndServe(":8080", mux)
    return nil
}
```

重新编译并运行proxy。
使用浏览器访问http://localhost:8080/swagger.json，即可得到hello world RESful API的具体说明。
 
 
###4.3下载swagger github源码
```
cd /root/joworkspace/
git clone https://github.com/swagger-api/swagger-ui.git
将dist目录下的所有文件拷贝到项目目录pkg/apigateway/third_party/swagger_ui里面(已经copy过来了)
cd swagger-ui
cp -rf dist ~/goworkspace/src/notification/apigateway/third_party
cd ~/goworkspace/src/notification/apigateway/third_party
mv dist swagger-ui
 
修改swagger-ui/index.html,替换
//url: "https://petstore.swagger.io/v2/swagger.json",
        url:"http://192.168.0.3:8080/swagger.json",

```
cp -rf /root/joworkspace/swagger-ui/dist ~/goworkspace/src/notification/apigateway/third_party
mv dist swagger-ui
 
 
###4.4将swagger-ui文件编译成go的内置文件: 
1.安装go-bindata工具
go get -u github.com/jteeuwen/go-bindata/...
 
2.制作成go的内置数据文件  
``` 
cd /root/goworkspace/src/notification/apigateway
go-bindata --nocompress -pkg swagger -o pkg/ui/data/swagger/datafile.go third_party/swagger-ui/...

cd /root/goworkspace/src/notification/pkg/apigateway
go-bindata --nocompress -pkg swagger -o pkg/ui/data/swagger/datafile.go third_party/swagger-ui/...

``` 
最终生成的文件 pkg/ui/data/swagger/datafile.go

 ###4.5 swagger-ui的文件服务器  
1.elazarl/go-bindata-assetf将内置的数据文件对外提供http服务 
go get github.com/elazarl/go-bindata-assetfs/...
 
2.修改proxy.go ---> apigateway/main.go代码，添加swagger函数 
Docs/Backup/main-http.go修改之前的代码
Docs/Backup/main-swagger.go修改之后的代码

 
3)重新编译proxy
apigateway/main.go
go build main.go
  
nf_proto/notification.swagger.go
go build main.go  

4)重新启动gateway服务
使用浏览器查看swagger-ui
浏览器输入
http://192.168.0.3:8080/swagger-ui

输入框输入
http://192.168.0.3:8080/swagger.json
 
#refs  
>grpc-gateway set up Steps
https://blog.csdn.net/dapangzi88/article/details/63686334 

>grpc-gateway和swagger-ui界面 集群Set up Steps
https://blog.csdn.net/StephenLu0422/article/details/82757905#2helloworldproto_109

>gRPC helloworld service, restful gateway and swagger 基于python
https://github.com/lienhua34/notes/blob/master/grpc/helloworld_restful_swagger/README.md
 
 
#如何生成gomock的代码
```
go get github.com/golang/mock/gomock
go get github.com/golang/mock/mockgen
cd /notification/pkg/pb
mockgen -source=notification.pb.go > ../mock/mockgen/nfmock.go

```

cd /notification/pkg/pb
mockgen -source=notification.pb.go > ../services/mock/mockgen/nf_handler_mock.go