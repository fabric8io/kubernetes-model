#!/usr/bin/python
import json
import sys
import urllib
import os

def all_same(items):
    return all(x == items[0] for x in items)
	
kubernetes_rev = None
kubernetes_comment = None
dependencies_list = []
openshift_rev = os.environ.get('GIT_COMMIT')

# Get the Kubernetes revision and comment that OpenShift Origin is using
with open("./Godeps/Godeps.json") as json_data:
    origin_data = json.load(json_data)

OriginDeps = origin_data["Deps"]
for od in OriginDeps:
	if ("github.com/GoogleCloudPlatform/kubernetes/" in od["ImportPath"]) :
		kubernetes_rev = od["Rev"]
		kubernetes_comment = od["Comment"]
		dependencies_list.append(kubernetes_rev)

# Check all the Kubernetes package dependencies are the same version
if (not all_same(dependencies_list)) :
	sys.exit("Not all Kubernetes package dependency revisions are the same in OpenShift Origin")

# Update the schema generator dependencies to the correct Origin and Kubernetes versions
with open("./Godeps/Godeps.json") as json_data:
    schema_data = json.load(json_data)

SchemaDeps = schema_data["Deps"]
for sd in SchemaDeps:
	if ("github.com/GoogleCloudPlatform/kubernetes/" in sd["ImportPath"]) :
		sd["Rev"] = kubernetes_rev
		sd["Comment"] = kubernetes_comment
	if ("github.com/openshift/origin/pkg/" in sd["ImportPath"]) :
		# Set Origin dependency to current Git commit that triggered this Job
		sd["Rev"] = openshift_rev
		# Clean up comment as not relavant to this commit anymore
		if "Comment" in sd:
			del sd["Comment"]

with open('./Godeps/Godeps.json', 'w') as outfile:
	json.dump(schema_data, outfile)

print "Updated Schema GoDeps.json..."
print "Latest OpenShift Origin revision:",openshift_rev
print "OpenShift rebased Kubernetes revision:",kubernetes_rev



