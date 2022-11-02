proto:
	protoc --go_out=./pkg/rpc --go_opt=paths=source_relative --go-grpc_out=./pkg/rpc --go-grpc_opt=paths=source_relative proto/cropper.proto