#
# Copyright (C) 2011 Red Hat, Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#         http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

SHELL := /bin/bash
tag := $(shell cat .openshift-version)

build:
	CGO_ENABLED=0 godep go build -a ./cmd/generate/generate.go
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
		godep update ... && \
		rm -rf Godeps/_workspace/src/github.com/GoogleCloudPlatform/kubernetes && \
		cp -r $(GOPATH)/src/github.com/openshift/origin/Godeps/_workspace/src/github.com/GoogleCloudPlatform/kubernetes Godeps/_workspace/src/github.com/GoogleCloudPlatform/kubernetes
