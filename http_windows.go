// #############################################################################
// # File: http_windows.go                                                     #
// # Project: gracefulshut                                                     #
// # Created Date: 2023/08/22 09:15:26                                         #
// # Author: realjf                                                            #
// # -----                                                                     #
// # Last Modified: 2023/08/22 10:52:41                                        #
// # Modified By: realjf                                                       #
// # -----                                                                     #
// # Copyright (c) 2023                                                        #
// #############################################################################
// +build: windows
package gracefulshut

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type gracefulHttp struct {
	server     *http.Server
	kernel32   *syscall.LazyDLL
	quit       chan os.Signal
	ctx        context.Context
	cancelFunc context.CancelFunc
	TLSConf    *TLSConf
}

type TLSConf struct {
	CertFile string
	KeyFile  string
}

func NewGracefulHttpServer(addr string, ctx context.Context) GracefulShut {
	s := &http.Server{
		Addr: addr,
	}
	return newGracefulHttp(s, ctx, nil)
}

func WrapHttpServer(s *http.Server, ctx context.Context) GracefulShut {
	return newGracefulHttp(s, ctx, nil)
}

func WrapHttpServerWithTLS(s *http.Server, ctx context.Context, conf *TLSConf) GracefulShut {
	return newGracefulHttp(s, ctx, conf)
}

func newGracefulHttp(s *http.Server, parent context.Context, conf *TLSConf) *gracefulHttp {
	ctx, cancel := context.WithTimeout(parent, 30*time.Second)
	return &gracefulHttp{
		server:     s,
		kernel32:   syscall.NewLazyDLL("kernel32.dll"),
		quit:       make(chan os.Signal, 1),
		ctx:        ctx,
		cancelFunc: cancel,
		TLSConf:    conf,
	}
}

func (g *gracefulHttp) Setup() {
	signal.Notify(g.quit, os.Interrupt, syscall.SIGTERM)
	if g.server.TLSConfig != nil {
		go func() {
			if err := g.server.ListenAndServeTLS(g.TLSConf.CertFile, g.TLSConf.KeyFile); err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Fatal(err)
			}
		}()
	} else {
		go func() {
			if err := g.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Fatal(err)
			}
		}()
	}

	log.Println("listen on: ", g.server.Addr)
}

func (g *gracefulHttp) Shutdown() (err error) {
	defer g.cancelFunc()
	sig := <-g.quit
	log.Printf("Received signal: %v\n", sig)

	err = g.server.Shutdown(g.ctx)
	close(g.quit)
	return
}

func (g *gracefulHttp) GetPid() (pid int) {
	return os.Getpid()
}
