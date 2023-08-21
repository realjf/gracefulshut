// #############################################################################
// # File: grpc.go                                                             #
// # Project: gracefulshut                                                     #
// # Created Date: 2023/08/21 18:16:07                                         #
// # Author: realjf                                                            #
// # -----                                                                     #
// # Last Modified: 2023/08/21 18:28:31                                        #
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

func newGracefulGrpc(listener net.Listener, parent context.Context) *gracefulGrpc {
	ctx, cancel := context.WithTimeout(parent, 30*time.Second)
	return &gracefulGrpc{
		server:     grpc.NewServer(),
		listener:   listener,
		quit:       make(chan os.Signal, 1),
		ctx:        ctx,
		cancelFunc: cancel,
	}
}

func (g *gracefulGrpc) Setup() {
	signal.Notify(g.quit, syscall.SIGINT, syscall.SIGTERM)
	if err := g.server.Serve(g.listener); err != nil {
		log.Fatal(err)
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
