SHELL := /bin/bash

build:
	godep go build ./cmd/generate/generate.go
	./generate > kubernetes-model/src/main/resources/schema/kube-schema.json
	mvn clean install

update-deps:
	pushd $(GOPATH)/src/github.com/openshift/origin && \
		godep restore && \
		popd && \
		godep save cmd/generate/generate.go && \
		godep update ...
