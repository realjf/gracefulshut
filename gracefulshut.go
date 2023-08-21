// #############################################################################
// # File: gracefulshut.go                                                     #
// # Project: gracefulshut                                                     #
// # Created Date: 2023/08/21 16:34:55                                         #
// # Author: realjf                                                            #
// # -----                                                                     #
// # Last Modified: 2023/08/21 18:09:21                                        #
// # Modified By: realjf                                                       #
// # -----                                                                     #
// # Copyright (c) 2023                                                        #
// #############################################################################
package gracefulshut

type GracefulShut interface {
	Setup() // not block
	Shutdown() error
	GetPid() (pid int)
}
