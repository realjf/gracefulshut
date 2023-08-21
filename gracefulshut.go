// #############################################################################
// # File: gracefulshut.go                                                     #
// # Project: gracefulshut                                                     #
// # Created Date: 2023/08/21 16:34:55                                         #
// # Author: realjf                                                            #
// # -----                                                                     #
// # Last Modified: 2023/08/22 06:24:58                                        #
// # Modified By: realjf                                                       #
// # -----                                                                     #
// # Copyright (c) 2023                                                        #
// #############################################################################
package gracefulshut

import "context"

type GracefulShut interface {
	Setup() // not block
	Shutdown() error
	GetPid() (pid int)
}

type Program interface {
	Setup(ctx context.Context) error    // run your program
	Shutdown(ctx context.Context) error // block
}
