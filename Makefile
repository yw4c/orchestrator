gen_pb:
	protoc --gofast_out=plugins=grpc:. ./pb/booking.proto

test_example:
	grpcurl -plaintext localhost:10000 list
	grpcurl -rpc-header x-request-id:example-request-id -plaintext -d '{"ProductID": "1", "FaultInject": "true"}' localhost:10000 pb.BookingService/HandleSyncBooking
	grpcurl -rpc-header x-request-id:example-request-id -plaintext -d '{"ProductID": "1", "FaultInject": "true"}' localhost:10000 pb.BookingService/HandleAsyncBooking

run_local:
	docker build -t orchestrator:latest .
	docker run orchestrator

deploy-dev:
	kubectl create ns orchestrator
	kustomize build ./deployment/k8s/dev | kubectl apply -f - -n orchestrator

deploy-relative:
	#helm install rabbitmq --set auth.username=guest,auth.password=guest bitnami/rabbitmq -n orchestrator
	helm install nats bitnami/nats -n orchestrator

undeploy-relative:
	#helm uninstall rabbitmq -n orchestrator
	helm uninstall nats bitnami/nats -n orchestrator