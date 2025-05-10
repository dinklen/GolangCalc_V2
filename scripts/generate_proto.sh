protoc -I ../api/proto \
	--go_out=paths=source_relative:../api/proto/generated \
	--go-grpc_out=paths=source_relative:../api/proto/generated \
	../api/proto/calculator_service.proto
