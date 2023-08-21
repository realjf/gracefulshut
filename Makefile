# ##############################################################################
# # File: Makefile                                                             #
# # Project: gracefulshut                                                      #
# # Created Date: 2023/08/21 17:11:57                                          #
# # Author: realjf                                                             #
# # -----                                                                      #
# # Last Modified: 2023/08/21 20:39:37                                         #
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


.PHONY: proto
proto:
	@protoc --proto_path=. --go_out=. --go-grpc_out=. calculator.proto
