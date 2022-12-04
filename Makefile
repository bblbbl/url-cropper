proto:
	protoc --go_out=./pkg/rpc --go_opt=paths=source_relative --go-grpc_out=./pkg/rpc --go-grpc_opt=paths=source_relative proto/cropper.proto

test:
	go test ./...

test-cover:
	go test ./... -cover -coverprofile=coverage.out

show-cover:
	go tool cover -html=coverage.out

migration:
	migrate create -ext sql -dir migrations -seq $(name)