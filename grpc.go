// #############################################################################
// # File: grpc.go                                                             #
// # Project: gracefulshut                                                     #
// # Created Date: 2023/08/21 18:16:07                                         #
// # Author: realjf                                                            #
// # -----                                                                     #
// # Last Modified: 2023/08/21 20:04:30                                        #
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
	server     *grpc.Server
	listener   net.Listener
	quit       chan os.Signal
	ctx        context.Context
	cancelFunc context.CancelFunc
}

func newGracefulGrpc(server *grpc.Server, listener net.Listener, parent context.Context) *gracefulGrpc {
	ctx, cancel := context.WithTimeout(parent, 30*time.Second)
	return &gracefulGrpc{
		server:     server,
		listener:   listener,
		quit:       make(chan os.Signal, 1),
		ctx:        ctx,
		cancelFunc: cancel,
	}
}

func NewGracefulGrpcServer(addr string, server *grpc.Server) (GracefulShut, error) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	return newGracefulGrpc(server, listener, context.Background()), nil
}

func WrapGrpc(server *grpc.Server, listener net.Listener, ctx context.Context) GracefulShut {
	return newGracefulGrpc(server, listener, ctx)
}

func (g *gracefulGrpc) Setup() {
	signal.Notify(g.quit, syscall.SIGINT, syscall.SIGTERM)

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
