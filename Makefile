SHELL := /bin/bash
tag := $(shell cat .openshift-version)

build:
	godep go build ./cmd/generate/generate.go
	./generate > kubernetes-model/src/main/resources/schema/kube-schema.json
	mvn clean install

update-deps:
	echo $(tag) > .openshift-version && \
		pushd $(GOPATH)/src/github.com/openshift/origin && \
		git fetch origin && \
		git checkout -B $(tag) refs/tags/$(tag) && \
		godep restore && \
		popd && \
		godep save cmd/generate/generate.go && \
		godep update ...
