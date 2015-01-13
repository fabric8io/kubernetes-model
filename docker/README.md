CI JSON Schema Generator for Kubernetes and OpenShift v3 Origin API Objects
===========================================================================

`rawlingsj/origin-schema-generator` is a docker image used by the fabric8 project to automatically generate and test internals used to interact with OpenShift and Kubernetes API's.

The [Schema Generator](https://github.com/fabric8io/origin-schema-generator) generates a JSON Schema based on updated go package dependencies from [Origin](https://github.com/openshift/origin) and [Kubernetes](https://github.com/GoogleCloudPlatform/kubernetes), which in turn is used to generate Java types in fabric8 representing each API and enables communication with both via REST at runtime.  These API's are updated when a Kubernetes rebase occurs in OpenShift.  fabric8 is required to generate its updated Java types, compile, run, unit test, integration test and if successuful submit a Pull Request on the [fabric8](https://github.com/fabric8io/fabric8) project.

This image is a building block towards Continuous Delivery for [fabric8](https://github.com/fabric8io/fabric8) aiming to automate the process when updating integration points of upstream projects.

- trigger on updates to [OpenShift v3](https://github.com/openshift/origin)
- pull the latest [Schema Generator](https://github.com/fabric8io/origin-schema-generator) 
- update Kubernetes and Origin pkg dependencies
- generate new JSON schema and copy to [fabric8/components/kubernetes-api](https://github.com/fabric8io/fabric8/blob/master/components/kubernetes-api/src/main/kubernetes/api/doc/kube-schema.json)
- generate new fabric8 java types
- run fabric8 unit test suite
- notify IRC of CI job result - _not yet implemented_
- if successful creates a PR for [fabric8](https://github.com/fabric8io/fabric8) - _not yet implemented_

## To run...

```
docker run -p 8080:8080  rawlingsj/origin-schema-generator
```