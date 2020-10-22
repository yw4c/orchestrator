gen_pb:
	protoc --gofast_out=plugins=grpc:. ./pb/booking.proto

test_example:
	grpcurl -plaintext localhost:10000 list
	grpcurl -rpc-header x-request-id:example-request-id -plaintext -d '{"ProductID": "1", "FaultInject": "true"}' localhost:10000 pb.BookingService/HandleSyncBooking
	grpcurl -rpc-header x-request-id:example-request-id -plaintext -d '{"ProductID": "1", "FaultInject": "true"}' localhost:10000 pb.BookingService/HandleAsyncBooking

run_local:
	docker build -t orchestrator:latest .
	docker run orchestrator