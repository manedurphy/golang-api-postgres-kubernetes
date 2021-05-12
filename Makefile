build:
	docker build -t k8s/golang-start .

cluster:
	kind create cluster --config kind.yaml

load: build
	kind load docker-image k8s/golang-start

docker-push:
	docker build -t golang-postgres-api .
	docker tag golang-postgres-api manedurphy/golang-postgres-api
	docker push manedurphy/golang-postgres-api

secret:
	kubectl create secret generic db-connection --from-literal DSN="host=pg-service user=user password=password dbname=golang port=5432"

deploy:
	kubectl apply -f kubernetes/config
	kubectl apply -f kubernetes/services
	kubectl apply -f kubernetes/statefulsets
	kubectl wait --for=condition=Ready --timeout=5m pod -l pg=pgdb
	kubectl apply -f kubernetes/deployments

destroy:
	kubectl delete -f kubernetes/config
	kubectl delete -f kubernetes/services
	kubectl delete -f kubernetes/statefulsets
	kubectl delete -f kubernetes/deployments

forward:
	kubectl port-forward service/goapp-service 8080:8080

linode:
	docker run --rm -it -v $(shell pwd):/work -w /work --entrypoint /bin/bash manedurphy/linode-cli
