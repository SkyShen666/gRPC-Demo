## 1. RPC

**RPC（Remote Procedure Call Protocol）**——远程过程调用协议，它是一种通过网络从远程计算机程序上请求服务，而不需要了解底层网络技术的协议. 

简单来说，就是跟远程访问或者web请求差不多，都是一个client向远端服务器请求服务返回结果，但是web请求使用的网络协议是http高层协议，而rpc所使用的协议多为TCP，是网络层协议，减少了信息的包装，加快了处理速度.

**所谓RPC框架实际是提供了一套机制，使得应用程序之间可以进行通信，而且也遵从server/client模型。使用的时候客户端调用server端提供的接口就像是调用本地的函数一样。**如下图所示就是一个典型的RPC结构图。

![img](https://upload-images.jianshu.io/upload_images/3959253-76284b64125a8673.png?imageMogr2/auto-orient/strip|imageView2/2/format/webp)



golang本身有rpc包，可以方便的使用，来构建自己的rpc服务，下边是一个简单是实例，可以加深我们的理解.

![1604992692504](C:\Users\Shen\AppData\Roaming\Typora\typora-user-images\1604992692504.png)

1.调用客户端句柄；执行传送参数 

2.调用本地系统内核发送网络消息 

3.消息传送到远程主机 

4.服务器句柄得到消息并取得参数 

5.执行远程过程 

6.执行的过程将结果返回服务器句柄 

7.服务器句柄返回结果，调用远程系统内核 

8.消息传回本地主机 

9.客户句柄由内核接收消息 

10.客户接收句柄返回的数据 



### 服务端

```go
package main

import (
	"io"
	"net"
	"net/http"
	"net/rpc"

	"github.com/astaxie/beego"
)

//- 方法是导出的
//- 方法有两个参数，都是导出类型或内建类型
//- 方法的第二个参数是指针
//- 方法只有一个error接口类型的返回值
//
//func (t *T) MethodName(argType T1, replyType *T2) error

type Panda int

func (this *Panda) Getinfo(argType int, replyType *int) error {
	beego.Info(argType)
	*replyType = 1 + argType
	return nil
}

func main() {
	//注册1个页面请求
	http.HandleFunc("/panda", pandatext)

	//new 一个对象
	pd := new(Panda)
	//注册服务
	//Register在默认服务中注册并公布 接收服务 pd对象 的方法
	rpc.Register(pd)

	rpc.HandleHTTP()
	//建立网络监听
	ln, err := net.Listen("tcp", "127.0.0.1:10086")
	if err != nil {
		beego.Info("网络连接失败")
	}

	beego.Info("正在监听10086")
	
	//service接受侦听器l上传入的HTTP连接，
	http.Serve(ln, nil)
}

//用来现实网页的web函数
func pandatext(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "panda")
}
```



### 客户端

```go
package main

import (
	"net/rpc"

	"github.com/astaxie/beego"
)

func main() {
	//rpc的与服务端建立网络连接
	cli, err := rpc.DialHTTP("tcp", "127.0.0.1:10086")
	if err != nil {
		beego.Info("网络连接失败")
	}

	var val int

	//远程调用函数（被调用的方法，传入的参数 ，返回的参数）
	err = cli.Call("Panda.Getinfo", 123, &val)
	if err != nil {
		beego.Info("打call失败")
	}

	beego.Info("返回结果", val)
}
```



## 2. gRPC是什么？

在 gRPC里客户端应用可以像调用本地对象一样直接调用另一台不同的机器上服务端应用的方法，使得您能够更容易地创建分布式应用和服务。与许多 RPC系统类似， gRPC也是基于以下理念： 

* 定义一个服务，指定其能够被远程调用的方法（包含参数和返回类型）。 

* 在服务端实现这个接口，并运行一个 gRPC服务器来处理客户端调用。 

* 在客户端拥有一个存根能够像服务端一样的方法。 gRPC客户端和服务端可以在多种环境中运行和交互 -从 google内部的服务器到你自己的笔记本，并且可以用任何 gRPC支持的语言 来编写。 

所以，你可以很容易地用 Java创建一个 gRPC服务端，用 Go、 Python、Ruby来创建客户端。此外， Google最新API将有 gRPC版本的接口，使你很容易地将 Google的功能集成到你的应用里。



## 3. gRPC有什么好处以及在什么场景下需要用gRPC

既然是server/client模型，那么我们直接用restful api不是也可以满足吗，为什么还需要RPC呢？下面我们就来看看RPC到底有哪些优势



### 3.1 gRPC vs. RESTful API

gRPC和restful API都提供了一套通信机制，用于server/client模型通信，而且它们都使用http作为底层的传输协议(严格地说, gRPC使用的http2.0，而restful api则不一定)。不过gRPC还是有些特有的优势，如下：

- gRPC可以通过protobuf来定义接口，从而可以有更加严格的接口约束条件。关于protobuf可以参见笔者之前的小文[Google Protobuf简明教程](https://www.jianshu.com/p/b723053a86a6) 
- 另外，通过protobuf可以将数据序列化为二进制编码，这会大幅减少需要传输的数据量，从而大幅提高性能。
- gRPC可以方便地支持流式通信(理论上通过http2.0就可以使用streaming模式, 但是通常web服务的restful api似乎很少这么用，通常的流式数据应用如视频流，一般都会使用专门的协议如HLS，RTMP等，这些就不是我们通常web服务了，而是有专门的服务器应用。）



### 3.2 使用场景

- 需要对接口进行严格约束的情况，比如我们提供了一个公共的服务，很多人，甚至公司外部的人也可以访问这个服务，这时对于接口我们希望有更加严格的约束，我们不希望客户端给我们传递任意的数据，尤其是考虑到安全性的因素，我们通常需要对接口进行更加严格的约束。这时gRPC就可以通过protobuf来提供严格的接口约束。
- 对于性能有更高的要求时。有时我们的服务需要传递大量的数据，而又希望不影响我们的性能，这个时候也可以考虑gRPC服务，因为通过protobuf我们可以将数据压缩编码转化为二进制格式，通常传递的数据量要小得多，而且通过http2我们可以实现异步的请求，从而大大提高了通信效率。

但是，通常我们不会去单独使用gRPC，而是将gRPC作为一个部件进行使用，这是因为在生产环境，我们面对大并发的情况下，需要使用分布式系统来去处理，<font color = 'red'>而gRPC并没有提供分布式系统相关的一些必要组件。</font>而且，真正的线上服务还需要提供包括负载均衡，限流熔断，监控报警，服务注册和发现等等必要的组件。不过，这就不属于本篇文章讨论的主题了，我们还是先继续看下如何使用gRPC。



## 4. gRPC HelloWorld实例详解

gRPC的使用通常包括如下几个步骤：

1. 通过protobuf来定义接口和数据类型
2. 编写gRPC server端代码
3. 编写gRPC client端代码
    下面来通过一个实例来详细讲解上述的三步。
    下边的hello world实例完成之后，其目录结果如下：

![1599546368203](C:\Users\Shen\AppData\Roaming\Typora\typora-user-images\1599546368203.png)



### 3.1 定义接口和数据类型

#### 3.1.1 通过protobuf定义接口和数据类型

```protobuf
syntax = "proto3";

package rpc_package;

// define a service
service HelloWorldService {
    // define the interface and data type
    rpc SayHello (HelloRequest) returns (HelloReply) {}
}

// define the data type of request
message HelloRequest {
    string name = 1;
}

// define the data type of response
message HelloReply {
    string message = 1;
}
```

#### 3.1.2 使用gRPC protobuf生成工具生成对应语言的库函数

```powershell
python -m grpc_tools.protoc -I./protos --python_out=./rpc_package --grpc_python_out=./rpc_package ./protos/helloworld.proto
```

这个指令会自动生成rpc_package文件夹中的`helloworld_pb2.py`和`helloworld_pb2_grpc.py`

#### 3.1.3 gRPC server端代码

```python
#!/usr/bin/env python
# -*-coding: utf-8 -*-

from concurrent import futures
import grpc
import logging
import time

from helloworld_pb2 import HelloRequest, HelloReply
from helloworld_pb2_grpc import add_HelloWorldServiceServicer_to_server, HelloWorldServiceServicer



class Hello(HelloWorldServiceServicer):

    # 这里实现我们定义的接口
    def SayHello(self, request, context):
        return HelloReply(message='Hello, %s!' % request.name)


def serve():
    # 这里通过thread pool来并发处理server的任务
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))

    # 将对应的任务处理函数添加到rpc server中
    add_HelloWorldServiceServicer_to_server(Hello(), server)

    # 这里使用的非安全接口，世界gRPC支持TLS/SSL安全连接，以及各种鉴权机制
    server.add_insecure_port('[::]:50000')
    server.start()
    try:
        while True:
            time.sleep(60 * 60 * 24)
    except KeyboardInterrupt:
        server.stop(0)


if __name__ == "__main__":
    logging.basicConfig()
    serve()
```

#### 3.1.4 gRPC client端代码

```python
#!/usr/bin/env python
# -*- coding: utf-8 -*-

from __future__ import print_function
import logging

import grpc
from helloworld_pb2 import HelloRequest, HelloReply
from helloworld_pb2_grpc import HelloWorldServiceStub

def run():
    # 使用with语法保证channel自动close
    with grpc.insecure_channel('localhost:50000') as channel:
        # 客户端通过stub来实现rpc通信
        stub = HelloWorldServiceStub(channel)

        # 客户端必须使用定义好的类型，这里是HelloRequest类型
        response = stub.SayHello(HelloRequest(name='shen'))
    print ("hello client received: " + response.message)

if __name__ == "__main__":
    logging.basicConfig()
    run()
```

#### 3.1.5 演示

先执行server端代码

```powershell
python hello_server.py
```

接着执行client端代码如下：

```powershell
➜  grpc_test python hello_client.py
hello client received: Hello, shen!
```

![1599547222840](C:\Users\Shen\AppData\Roaming\Typora\typora-user-images\1599547222840.png)

代码位置：D:\develop\Python_WorkSpace\grpc_test



### References

[简书gRPC](https://www.jianshu.com/p/9c947d98e192)

[gRPC quick start](https://grpc.io/docs/languages/python/quickstart/)





<hr>

## 5. 跨语言实例:

[参考博客](https://github.com/leisurelicht/grpc-demo)

Go做服务端

Python做客户端

```shell
$ python -m grpc_tools.protoc -I ./ --python_out=./pythonProtobuf/ --grpc_python_out=./pythonProtobuf ./auth.proto 
```

![1605167968533](C:\Users\Shen\AppData\Roaming\Typora\typora-user-images\1605167968533.png)

若想客户端与服务端彻底分离,只需要将在服务端自动生成的.py文件,复制到客户端使用即可.

![1605168144189](C:\Users\Shen\AppData\Roaming\Typora\typora-user-images\1605168144189.png)

