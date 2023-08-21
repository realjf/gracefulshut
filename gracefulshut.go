// #############################################################################
// # File: gracefulshut.go                                                     #
// # Project: gracefulshut                                                     #
// # Created Date: 2023/08/21 16:34:55                                         #
// # Author: realjf                                                            #
// # -----                                                                     #
// # Last Modified: 2023/08/21 17:50:57                                        #
// # Modified By: realjf                                                       #
// # -----                                                                     #
// # Copyright (c) 2023                                                        #
// #############################################################################
package gracefulshut

type GracefulShut interface {
	Start()
	Shutdown() error
	GetPid() (pid int)
}
