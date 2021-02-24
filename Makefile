BINARY_NAME = faas-monitor
IMAGE = smvfal/faas-monitor
TAG = latest

publish: docker-build docker-push

docker-build:
	DOCKER_BUILDKIT=1 docker build -t ${IMAGE}:${TAG} .

docker-push:
	docker push ${IMAGE}:${TAG}

install:
	kubectl apply -f kubernetes/rbac.yml
	kubectl apply -f kubernetes/deployment.yml

vendor:
	go mod vendor -v

build:
	go build -o bin/${BINARY_NAME}

clean:
	rm -r vendor/ bin/
