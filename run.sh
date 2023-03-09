protoc calculator/calculatorpb/calculator.proto --go_out=:. --go-grpc_out=:.
go run server/server.go
#go run client/client.go
