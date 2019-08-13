.PHONY: protos

protos:
	protoc -I protos/ protos/dda.proto --go_out=plugins=grpc:protos/dda