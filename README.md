JSON Schema Generator for OpenShift v3 Origin API Objects
=========================================================

Uses Go reflection to generate a JSON schema that describes one or more 
API resources in Openshift Origin.

Pre-requisits
-------------

Install [go](https://golang.org/doc/install)   
Install [godep](https://github.com/GoogleCloudPlatform/kubernetes/blob/master/docs/devel/development.md#godep-and-dependency-management)   

Building
--------
To build, clone repo and run:  

```
cd origin-schema-generator
godep restore
godep go build ./cmd/generate/generate.go  
./generate > kube-schema.json  
```

You should now be able to view the generated schema in `kube-schema.json`

Update dependency API's
-----------------------

Following [godep](https://github.com/tools/godep/blob/master/Readme.md)   

To update Kubernetes   
go get -u github.com/GoogleCloudPlatform/kubernetes/   
godep update github.com/GoogleCloudPlatform/kubernetes/...   

To update Openshift   
go get -u github.com/openshift/origin/   
godep update github.com/openshift/origin/...   
