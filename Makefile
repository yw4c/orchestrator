IMAGE_NAME = "yw4code/orchestrator"

gen_pb:
	protoc --gofast_out=plugins=grpc:. ./pb/booking.proto

test_example:
	grpcurl -plaintext localhost:10000 list
	grpcurl -rpc-header x-request-id:example-request-id -plaintext -d '{"ProductID": "1", "FaultInject": "true"}' localhost:10000 pb.BookingService/HandleSyncBooking
	grpcurl -rpc-header x-request-id:example-request-id -plaintext -d '{"ProductID": "1", "FaultInject": "true"}' localhost:10000 pb.BookingService/HandleAsyncBooking

update_image:
	docker build -t "${IMAGE_NAME}:latest" .
	docker image push "${IMAGE_NAME}:latest"

deploy-dev:
	kustomize build ./deployment/k8s/dev | kubectl apply -f - -n orchestrator

