# JSON Schema Generator for OpenShift v3 Origin API Objects

Uses Go reflection to generate a JSON schema that describes one or more 
API resources in Openshift Origin.

## Pre-requisits

Install [go](https://golang.org/doc/install)   
Install [glide](https://github.com/Masterminds/glide#install)   


## Getting the code

```
git clone https://github.com/fabric8io/kubernetes-model $GOPATH/src/github.com/fabric8io/kubernetes-model
```


## Building

To build, clone repo and run:  

```
make
```

You should now be able to view the generated schema in `kube-schema.json`

## Update dependency API's

To update openshift/kubernetes dependencies, run:
```
make vendoring
```