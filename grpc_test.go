// #############################################################################
// # File: grpc_test.go                                                        #
// # Project: gracefulshut                                                     #
// # Created Date: 2023/08/21 18:16:20                                         #
// # Author: realjf                                                            #
// # -----                                                                     #
// # Last Modified: 2023/08/21 21:05:18                                        #
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

	"google.golang.org/grpc"

	"github.com/realjf/gracefulshut"
	"github.com/realjf/gracefulshut/pb"
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
				calcServer.Add(context.Background(), &pb.AddRequest{
					Num1: 1,
					Num2: 2,
				})
				cmd := exec.Command("kill", "-2", fmt.Sprintf("%d", g.GetPid()))
				var out bytes.Buffer
				cmd.Stdout = &out
				err := cmd.Run()
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
