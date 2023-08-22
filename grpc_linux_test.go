// #############################################################################
// # File: grpc_linux_test.go                                                  #
// # Project: gracefulshut                                                     #
// # Created Date: 2023/08/21 18:16:20                                         #
// # Author: realjf                                                            #
// # -----                                                                     #
// # Last Modified: 2023/08/22 09:31:41                                        #
// # Modified By: realjf                                                       #
// # -----                                                                     #
// # Copyright (c) 2023                                                        #
// #############################################################################
package gracefulshut_test

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net"
	"os/exec"
	"testing"
	"time"

	"github.com/realjf/gracefulshut"
	"github.com/realjf/gracefulshut/pb"
	"google.golang.org/grpc"
)

func TestGrpcServer(t *testing.T) {
	server := grpc.NewServer()
	listener, err := net.Listen("tcp", "127.0.0.1:5555")
	if err != nil {
		t.Fatal(err)
	} else {
		calcServer := &gracefulshut.CalculatorServer{}
		pb.RegisterCalculatorServiceServer(server, calcServer)
		g := gracefulshut.WrapGrpcServer(server, listener, context.Background())
		g.Setup()
		go func() {
			// kill after 3s
			for range time.After(3 * time.Second) {
				res, err := calcServer.Add(context.Background(), &pb.AddRequest{
					Num1: 1,
					Num2: 2,
				})
				if err != nil {
					log.Printf("error: %v", err)
				} else {
					log.Printf("result: %v", res)
				}
				cmd := exec.Command("kill", "-2", fmt.Sprintf("%d", g.GetPid()))
				var out bytes.Buffer
				cmd.Stdout = &out
				err = cmd.Run()
				if err != nil {
					log.Fatal(err)
				} else {
					return
				}
			}
		}()
		if err := g.Shutdown(); err != nil {
			t.Fatal(err)
		} else {
			t.Log("done")
		}
	}
}
