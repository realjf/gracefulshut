// #############################################################################
// # File: main.go                                                             #
// # Project: example                                                          #
// # Created Date: 2023/08/21 21:03:35                                         #
// # Author: realjf                                                            #
// # -----                                                                     #
// # Last Modified: 2023/08/22 10:25:24                                        #
// # Modified By: realjf                                                       #
// # -----                                                                     #
// # Copyright (c) 2023                                                        #
// #############################################################################
// +build: linux || darwin
package main

import (
	"context"
	"log"
	"net"

	"github.com/realjf/gracefulshut"
	"github.com/realjf/gracefulshut/pb"
	"google.golang.org/grpc"
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
	res, err := calcServer.Add(context.Background(), &pb.AddRequest{
		Num1: 1,
		Num2: 2,
	})
	if err != nil {
		log.Printf("error: %v", err)
	} else {
		log.Printf("result: %v", res)
	}
	if err := g.Shutdown(); err != nil {
		log.Panic(err)
	}
}
