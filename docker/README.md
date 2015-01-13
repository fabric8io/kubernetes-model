CI JSON Schema Generator for Kubernetes and OpenShift v3 Origin API Objects
===========================================================================

`rawlingsj/origin-schema-generator` is a docker image that contains a Jenkins job that will..

- trigger on updates to [OpenShift v3](https://github.com/openshift/origin)
- pull the latest [Schema Generator](https://github.com/fabric8io/origin-schema-generator) 
- update Kubernetes and Origin pkg dependencies
- generate new JSON schema
- generate new fabric8 java types
- run fabric8 unit test suite
- notify IRC of CI job result
- if successful creates a PR for fabric8/origin-schema-generator

To run..

```
docker run -p 8080:8080  rawlingsj/origin-schema-generator
```