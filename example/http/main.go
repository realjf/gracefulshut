// #############################################################################
// # File: main.go                                                             #
// # Project: http                                                             #
// # Created Date: 2023/08/21 21:08:00                                         #
// # Author: realjf                                                            #
// # -----                                                                     #
// # Last Modified: 2023/08/21 21:08:19                                        #
// # Modified By: realjf                                                       #
// # -----                                                                     #
// # Copyright (c) 2023                                                        #
// #############################################################################
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
