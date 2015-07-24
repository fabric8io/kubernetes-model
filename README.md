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
make
```

You should now be able to view the generated schema in `kube-schema.json`

Update dependency API's
-----------------------

To update openshift/kubernetes dependencies, run:

    make [tag=v0.5.2] update-deps

Where the optional tag value is the tagged version of OpenShift. This will update all
dependencies to those consistent with Openshift dependencies, including Kubernetes.
This command should also be run when you need to have any new dependencies included
in Godeps workspace, e.g. adding a new package for schema generation.

If you do not specify a tag value then the tag value will be read from .openshift-version
in the root of the schema generator repo.

So if you're just looking to ensure all dependencies are vendored in the godep workspace
run:

    make update-deps

If you're looking to update the version of openshift & it's dependencies run:

    make tag=<openshift_tag> update-deps
