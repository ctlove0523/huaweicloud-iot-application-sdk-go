# huaweicloud-iot-application-sdk-go

huaweicloud-iot-application-sdk-go封装了华为云IoT物联网平台提供的API，使用SDK可以减低API使用的复杂度和难度，更加快速开发应用平台。

支持如下功能：

* 产品管理
* 设备管理
* 设备消息
* 设备命令
* 设备属性
* AMQP队列管理
* 接入凭证管理
* 资源空间管理
* ......

## 设计理念

华为云IoT提供基础版、标准版和企业版三种类型的实例，可以购买一个或多个不同类型的实例，为了更加方便的使用API，在huaweicloud-iot-application-sdk-go中一个client对应一个实例，如果你有多个实例则需要构建多个client以用来访问API。client支持AK/SK以及Token进行鉴权，优先推荐使用AK/SK。

## 安装和构建

安装和构建的过程取决于是使用go的 [modules](https://golang.org/ref/mod)(推荐) 还是还是`GOPATH`

### Modules

如果你使用 [modules](https://golang.org/ref/mod) 只需要导入包"github.com/ctlove0523/[huaweicloud-iot-application-sdk-go"即可使用。当你使用go
build命令构建项目时，依赖的包会自动被下载。注意使用go build命令构建时会自动下载最新版本，最新版本还没有达到release的标准可能存在一些尚未修复的bug。如果想使用稳定的发布版本可以从[release](https://github.com/ctlove0523/huaweicloud-iot-application-sdk-go/releases)获取最新稳定的版本号，并在go.mod文件中指定版本号。

~~~go
module example

go 1.15

require github.com/ctlove0523/[huaweicloud-iot-application-sdk-go v0.0.1-alpha
~~~

### GOPATH

如果你使用GOPATH，下面的一条命令即可实现安装

~~~go
go get github.com/ctlove0523/[huaweicloud-iot-application-sdk-go
~~~

## 使用Client

### 创建同步Client

1、创建使用AK/SK鉴权的Client：

~~~go
options := iot.ApplicationOptions{
	ServerPort:    443,
	ServerAddress: "iotda.cn-north-4.myhuaweicloud.com",
	InstanceId:    "",
	ProjectId:     "25e1be7c374749e9b6a25bc4ad53393a",

	Credential: &iot.Credentials{
		Ak:      "xxx",
		Sk:      "xxx",
		UseAkSk: true,
	},
}

client := iot.CreateSyncIotApplicationClient(options)
~~~

2、如果不想使用AK/SK鉴权，还可以创建使用Token鉴权的Client

~~~go
options := iot.ApplicationOptions{
	ServerPort:    443,
	ServerAddress: "iotda.cn-north-4.myhuaweicloud.com",
	InstanceId:    "",
	ProjectId:     "25e1be7c374749e9b6a25bc4ad53393a",

	Credential: &iot.Credentials{
		UseAkSk: false,
		Token:   "xxx",
	},
}

client := iot.CreateSyncIotApplicationClient(options)
~~~

### 使用Client调用API

SDK中所有的方法返回值都为（x,y）格式，x根据不同的方法返回的对象不同，y都为Go的error，在使用结果x之前应当首先检查y是否为nil，也就是检查方法调用是否成功，只有方法调用成功时结果x才是可用的。下面以查询AMQP队列为例说明：

~~~go
queues, err := client.ListAmqpQueues(iot.ListAmqpQueuesRequest{})
if err != nil {  // 首先检查方法调用是否成功
	fmt.Println(err)
	panic(1)
}

fmt.Println(queues.Queues)  //方法调用成功，可以使用方法返回的结果
~~~



### 更多样例：

samples包中有更多使用样例。

## 报告bugs

如果你在使用过程中遇到任何问题或bugs，请通过issue的方式上报问题或bug，我们将会在第一时间内答复。上报问题或bugs时请尽量提供以下内容：

* 使用的版本
* 使用场景
* 重现问题或bug的样例代码
* 错误信息
* ······

## 贡献

该项目欢迎来自所有人的pull request。
