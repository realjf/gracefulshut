// #############################################################################
// # File: program.go                                                          #
// # Project: gracefulshut                                                     #
// # Created Date: 2023/08/21 21:20:10                                         #
// # Author: realjf                                                            #
// # -----                                                                     #
// # Last Modified: 2023/08/22 06:25:27                                        #
// # Modified By: realjf                                                       #
// # -----                                                                     #
// # Copyright (c) 2023                                                        #
// #############################################################################
package gracefulshut

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type gracefulProgram struct {
	program    Program
	quit       chan os.Signal
	ctx        context.Context
	cancelFunc context.CancelFunc
}

func newGracefulProgram(parent context.Context, program Program) GracefulShut {
	ctx, cancel := context.WithTimeout(parent, 30*time.Second)
	return &gracefulProgram{
		program:    program,
		quit:       make(chan os.Signal, 1),
		ctx:        ctx,
		cancelFunc: cancel,
	}
}

func NewGracefulProgram(program Program) GracefulShut {
	return newGracefulProgram(context.Background(), program)
}

func WrapProgram(ctx context.Context, program Program) GracefulShut {
	return newGracefulProgram(ctx, program)
}

func (g *gracefulProgram) Setup() {
	signal.Notify(g.quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		if err := g.program.Setup(g.ctx); err != nil {
			log.Fatal(err)
		}
	}()

	log.Println("program start...")
}

func (g *gracefulProgram) Shutdown() (err error) {
	defer g.cancelFunc()
	sig := <-g.quit
	log.Printf("Received signal: %v\n", sig)

	err = g.program.Shutdown(g.ctx)
	close(g.quit)
	return nil
}

func (g *gracefulProgram) GetPid() (pid int) {
	return os.Getpid()
}
