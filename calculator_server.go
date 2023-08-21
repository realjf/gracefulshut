// #############################################################################
// # File: calculator_server.go                                                #
// # Project: gracefulshut                                                     #
// # Created Date: 2023/08/21 21:04:38                                         #
// # Author: realjf                                                            #
// # -----                                                                     #
// # Last Modified: 2023/08/21 21:05:13                                        #
// # Modified By: realjf                                                       #
// # -----                                                                     #
// # Copyright (c) 2023                                                        #
// #############################################################################
package gracefulshut

import (
	"context"
	"log"

	"github.com/realjf/gracefulshut/pb"
)

type CalculatorServer struct {
	*pb.UnimplementedCalculatorServiceServer
}

func (s *CalculatorServer) Add(ctx context.Context, req *pb.AddRequest) (*pb.AddResponse, error) {
	result := req.GetNum1() + req.GetNum2()
	log.Println("result: ", result)
	return &pb.AddResponse{Result: result}, nil
}
