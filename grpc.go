// #############################################################################
// # File: grpc.go                                                             #
// # Project: gracefulshut                                                     #
// # Created Date: 2023/08/21 18:16:07                                         #
// # Author: realjf                                                            #
// # -----                                                                     #
// # Last Modified: 2023/08/21 18:47:44                                        #
// # Modified By: realjf                                                       #
// # -----                                                                     #
// # Copyright (c) 2023                                                        #
// #############################################################################
package gracefulshut

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
)

type gracefulGrpc struct {
	server      *grpc.Server
	grpcServer  interface{}
	serviceDesc *grpc.ServiceDesc
	quit        chan os.Signal
	ctx         context.Context
	cancelFunc  context.CancelFunc
}

func newGracefulGrpc(serviceDesc *grpc.ServiceDesc, grpcServer interface{}, parent context.Context) *gracefulGrpc {
	ctx, cancel := context.WithTimeout(parent, 30*time.Second)
	return &gracefulGrpc{
		server:      grpc.NewServer(),
		grpcServer:  grpcServer,
		serviceDesc: serviceDesc,
		quit:        make(chan os.Signal, 1),
		ctx:         ctx,
		cancelFunc:  cancel,
	}
}

func NewGracefulGrpcServer(addr string) (GracefulShut, error) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	server := grpc.NewServer()

	return
}

func WrapGrpc(serviceDesc *grpc.ServiceDesc, grpcServer interface{}, ctx context.Context) GracefulShut {
	return newGracefulGrpc(serviceDesc, grpcServer, ctx)
}

func (g *gracefulGrpc) Setup() {
	signal.Notify(g.quit, syscall.SIGINT, syscall.SIGTERM)

	g.server.RegisterService(g.serviceDesc, g.grpcServer)
	go func() {
		if err := g.server.Serve(g.listener); err != nil {
			log.Fatal(err)
		}
	}()
	log.Println("Start Serve...")
	for n, s := range g.server.GetServiceInfo() {
		log.Printf("%s: %v", n, s)
	}
}

func (g *gracefulGrpc) Shutdown() error {
	defer g.cancelFunc()
	sig := <-g.quit
	log.Printf("Received signal: %v\n", sig)

	g.server.GracefulStop()
	close(g.quit)
	return nil
}

func (g *gracefulGrpc) GetPid() (pid int) {
	return os.Getpid()
}
