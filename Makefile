# ##############################################################################
# # File: Makefile                                                             #
# # Project: gracefulshut                                                      #
# # Created Date: 2023/08/21 17:11:57                                          #
# # Author: realjf                                                             #
# # -----                                                                      #
# # Last Modified: 2023/08/21 18:01:34                                         #
# # Modified By: realjf                                                        #
# # -----                                                                      #
# # Copyright (c) 2023                                                         #
# ##############################################################################


.PHONY: lint
lint:
	@golangci-lint run -v ./...


.PHONY: push
push:
	@git add -A && git commit -m "update" && git push origin master
