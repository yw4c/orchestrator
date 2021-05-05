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

k8s-nats:
	#helm install rabbitmq --set auth.username=guest,auth.password=guest bitnami/rabbitmq -n orchestrator
	# helm install nats bitnami/nats -n orchestrator
 	# https://nats-io.github.io/k8s/
	# kubectl create ns nats
	helm repo add nats https://nats-io.github.io/k8s/helm/charts/
	helm repo update
	helm install my-nats nats/nats -n nats 
	helm install my-stan nats/stan --set stan.nats.url=nats://my-nats:4222 -n nats 
k8s-redis:
	helm install redis bitnami/redis -n orchestrator

rm-k8s-nats:
	helm uninstall my-nats nats/nats -n nats
	helm uninstall my-stan nats/stan -n nats
