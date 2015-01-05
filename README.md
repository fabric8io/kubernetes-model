JSON Schema Generator for OpenShift v3 Origin API Objects
=========================================================

Uses Go reflection to generate a JSON schema that describes one or more 
API resources in Openshift Origin.

Pre-requisits
-------------

Have an up-to-date cloned copy of Kuberenetes following the [Development Guide](https://github.com/GoogleCloudPlatform/kubernetes/blob/master/docs/devel/development.md#development-guide) as the project needs to reside in the correct Go Workspace folder structure.  Ensure that [godep](https://github.com/GoogleCloudPlatform/kubernetes/blob/master/docs/devel/development.md#godep-and-dependency-management) is installed and the steps in [Using godep](https://github.com/GoogleCloudPlatform/kubernetes/blob/master/docs/devel/development.md#using-godep) have been executed.  For reference when you 'Set up your GOPATH' you should choose Option A.

Create a `csrwng` directory under the go working structure and clone this project into it.  You should have a folder structure similar to below..


> ~/go/src/github.com   
> > GoogleCloudPlatform  
> > > kubernetes  

> > csrwng  

> > > origin-schema-generator  



Building
--------
To build, run:  

```
cd origin-schema-generator  
godep go build ./cmd/generate/generate.go  
generate > kube-schema.json  
```

You should now be able to view the generated schema in `kube-schema.json`
