### run
```
- run server
go run server/grpc/server.go

- run client(2つ立てる)
go run cmd/main.go
```

### othres(memo)
```
- サーバーサイド側
protoc -I./proto --go_out=. proto/*.proto

- クライアント側
protoc -I./proto --go-grpc_out=. proto/*.proto
```