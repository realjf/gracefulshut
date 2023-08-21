// #############################################################################
// # File: http_test.go                                                        #
// # Project: gracefulshut                                                     #
// # Created Date: 2023/08/21 17:27:40                                         #
// # Author: realjf                                                            #
// # -----                                                                     #
// # Last Modified: 2023/08/21 18:09:42                                        #
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
	"os/exec"
	"testing"
	"time"

	"github.com/realjf/gracefulshut"
)

func TestHttp(t *testing.T) {
	g := gracefulshut.NewGracefulHttpServer("127.0.0.1:55555", context.Background())
	g.Setup()
	go func() {
		for range time.After(3 * time.Second) {
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
