# gomall

hertz框架，官网：https://www.cloudwego.io/zh/docs/hertz/getting-started/

IDL使用thrift和protobuf协议

thrift资料：

https://thrift.apache.org/docs/idl.html

https://diwakergupta.github.io/thrift-missing-guide/

protobuf资料：

https://protobuf.dev/programming-guides/proto3/

在centos7完成项目，安装需要环境

## 环境安装
### go环境
```bash
# 管理员模式
root 
# 安装wegt
sudo yum install wget -y

# 安装go
wget https://go.dev/dl/go1.24.4.linux-amd64.tar.gz

# 执行tar解压到/usr/loacl目录下（官方推荐），得到go文件夹等
tar -C /usr/local -zxvf go1.24.4.linux-amd64.tar.gz

# 添加/usr/loacl/go/bin目录到PATH变量中。添加到/etc/profile 或$HOME/.profile都可以
# 用vim，没有的话可以用命令`sudo apt-get install vim`安装一个
vim /etc/profile
# 在最后一行添加
export GOROOT=/usr/local/go
export PATH=$PATH:$GOROOT/bin
# 保存退出(ESC+:wq+enter)后source一下（vim 的使用方法可以自己搜索一下）
source /etc/profile
# 验证是否成功
go version

```

### 代码位置创建
```bash
# 先创建你的工作空间(Workspaces)，官方建议目录$HOME/go。
mkdir $HOME/go
# 编辑 ~/.bash_profile 文件
vim ~/.bash_profile
# 在最后一行添加下面这句。$HOME/go 为你工作空间的路径
export GOPATH=$HOME/go
# 保存退出后source
source ~/.bash_profile
### 将代码放到$HOME/go/src下
```

### Hertz依赖安装
```bash
# 安装hertz
go get -u github.com/cloudwego/hertz
# 进入hello_world文件夹执行go run main判断8080端口是否启动成功即为安装成功
# go mod init github.com/cloudwego/biz-demo/gomall/hello_world
# go mod tidy
# go run main.go
```

### cwgo安装
```bash
# 安装cwgo
GOPROXY=https://goproxy.cn/,direct go install github.com/cloudwego/cwgo@latest
# 打开 .bashrc 或 .zshrc 文件
vim ~/.bashrc
# 添加到文件末尾
export PATH=$PATH:$HOME/go/bin
# 保存并关闭文件，立即生效
source ~/.bashrc
cwgo --version
```

### 安装thrift
```bash
# 安装thrift
GO111MODULE=on go install github.com/cloudwego/thriftgo@latest
```

### 验证cwgo和thrift
如果报错
`The module name given by the '-module' option ('gomall/demo/demo_thrift') is not consist with the name defined in go.mod ('gomall' from /root/go/src/gomall)`，检查CloudWeGo 的 cwgo 工具需要确保：-module 参数填写的模块名 必须与 go.mod 文件中定义的 module 名一致或是其子路径

```bash
mkdir -p demo/demo_thrift
cd demo/demo_thrift
# 生成代码
cwgo server --type RPC --module gomall/demo/demo_thrift --service demo_thrift --idl ../../idl/echo.thrift
go mod tidy
go run .
# 显示
# &{Env:test Kitex:{Service:demo_thrift Address::8888 LogLevel:info LogFileName:log/kitex.log LogMaxSize:10 LogMaxBackups:50 LogMaxAge:3} MySQL:{DSN:gorm:gorm@tcp(127.0.0.1:3306)/gorm?charset=utf8mb4&parseTime=True&loc=Local} Redis:{Address:127.0.0.1:6379 Username: Password: DB:0} Registry:{RegistryAddress:[127.0.0.1:2379] Username: Password:}}
# 即为成功
```


### 安装protobuf
```bash
# 安装protobuf
# 使用uname -m查看版本，在https://github.com/protocolbuffers/protobuf/releases找到相应链接进行下载
wget https://github.com/protocolbuffers/protobuf/releases/download/v31.1/protoc-31.1-linux-x86_64.zip
unzip ./protoc-31.1-linux-x86_64.zip
sudo cp ./bin/protoc /usr/local/bin
sudo cp -r ./include/google /usr/local/include/
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
# 验证protobuf（可能需要重新打开命令行）
protoc --version
```

### 验证cwgo和protobuf
```bash
mkdir -p demo/demo_proto
cd demo/demo_proto
cwgo server -I ../../idl --type RPC --module gomall/demo/demo_proto --service demo_proto --idl ../../idl/echo.proto
go mod tidy
# 将当前目录作为一个模块，添加到 go.work 工作区文件中
go work init .         # 或先 go work init 再 go work use .
go run .
# 显示
# &{Env:test Kitex:{Service:demo_proto Address::8888 LogLevel:info LogFileName:log/kitex.log LogMaxSize:10 LogMaxBackups:50 LogMaxAge:3} MySQL:{DSN:gorm:gorm@tcp(127.0.0.1:3306)/gorm?charset=utf8mb4&parseTime=True&loc=Local} Redis:{Address:127.0.0.1:6379 Username: Password: DB:0} Registry:{RegistryAddress:[127.0.0.1:2379] Username: Password:}}
# 即为成功
```

### 读取.env作为环境变量
```bash
go get github.com/joho/godotenv
```

### consul服务发现
#### 安装
```bash
go get github.com/kitex-contrib/registry-consul
# 使用docker拉取consul镜像
docker pull hashicorp/consul
```


#### 启动consul容器
```bash
# 修改gomall下docker-compose.yaml配置
docker compose up -d
# 老版本可能需要如下指令
# docker-compose up -d
```

#### 验证consul是否正常运行
```bash
# 服务端
cd /gomall/demo/demo_proto
go run .
# 客户端
cd /gomall/demo/demo_proto/cmd/client
go run .
# 出现hello即为成功
```
具体节点信息可以打开Consul提供的一个WEB界面查看所有的节点，通过8500端口访问

### 前端代码生成命令
```bash
cwgo server --type HTTP --service frontend -module gomall/app/frontend -I ../../idl --idl ../../idl/frontend/home.proto 
```

### 热加载Air工具安装
```bash
go install github.com/air-verse/air@latest
# 优先在当前路径查找 `.air.toml` 后缀的文件，如果没有找到，则使用默认的
air -c .air.toml
# 在这之后，只需执行 air 命令，无需额外参数，它就能使用 .air.toml 文件中的配置了。
air
```

### 安装session
```bash
go get github.com/hertz-contrib/sessions
```

### 安装redis
```bash
docker-compose up -d --build
```

启动docker
```bash
systemctl start docker
```

启动mysql
```bash
docker-compose exec mysql bash
# 建表
create database test;
```

linux测试go test
```bash
go test -v user_test.go
```

引用项目需要的依赖增加到go.mod文件, 去掉go.mod文件中项目不需要的依赖:
```
go mod tidy
```


参考文献：
- https://wangyi.one/cloudwego%E5%AD%A6%E4%B9%A0%E7%AC%94%E8%AE%B0/
- https://www.cloudwego.io/zh/docs/cwgo/
- https://www.cloudwego.io/zh/docs/kitex/
consul相关资料：
- https://developer.hashicorp.com/consul/tutorials/get-started-vms/virtual-machine-gs-deploy
- https://hub.docker.com/r/hashicorp/consul