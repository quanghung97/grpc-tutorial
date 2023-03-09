install golang in disk /d/Go

source code: /d/Go/src/quanghung97

install protoc https://github.com/google/protobuf/releases

Extract all to D:/Go/protoc 
Your directory structure should now be
D:\protoc\bin 
D:\protoc\include
Finally, add D:\protoc\bin to your PATH:


updated .bashrc

export GOPATH='/d/Go/bin'
export PATH=$PATH:'/d/Go/bin'
export PATH=$PATH:'/d/Go/protoc/bin'
export PATH="$PATH:$GOPATH/bin"

install protoc-gen-go, protoc-gen-go-grpc

go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

////////////////////////////////////////////



syntax kieu du lieu pho bien
dia chi, con tro
cac framework beego, gin
ORM (mo hinh map kien truc backend)
implement kien truc vao golang

MVC (model-view-controller)
implement kien truc vao golang

socket.io golang

queue job (redis, database)

nang cao:
dependency injection
repository pattern (SOLID) 

grpc (tools for microservices) (streaming API (client, server) http/2)
