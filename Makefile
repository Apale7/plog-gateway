grpc:
	protoc --go_out=. proto/*.proto
	protoc --go-grpc_out=. proto/*.proto