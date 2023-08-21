# gracefulshut

graceful shutdown

### Features

- Support for the wrap HTTP/GRPC server

### Usage

the wrap HTTP server

```go
package main

import (
 "context"
 "log"
 "net/http"

 "github.com/realjf/gracefulshut"
)

func main() {
 s := &http.Server{
  Addr: "127.0.0.1:5555",
 }
 g := gracefulshut.WrapHttpServer(s, context.Background())
 g.Setup()
 if err := g.Shutdown(); err != nil {
  log.Fatal(err)
 }
}

```

the wrap GRPC server

```go
package main

import (
 "context"
 "log"
 "net"

 "google.golang.org/grpc"

 "github.com/realjf/gracefulshut"
 "github.com/realjf/gracefulshut/pb"
)

func main() {
 server := grpc.NewServer()
 listener, err := net.Listen("tcp", "127.0.0.1:5555")
 if err != nil {
  panic(err)
 }
 calcServer := &gracefulshut.CalculatorServer{}
 pb.RegisterCalculatorServiceServer(server, calcServer)
 g := gracefulshut.WrapGrpcServer(server, listener, context.Background())
 g.Setup()
 calcServer.Add(context.Background(), &pb.AddRequest{
  Num1: 1,
  Num2: 2,
 })
 if err := g.Shutdown(); err != nil {
  log.Panic(err)
 }
}

```
