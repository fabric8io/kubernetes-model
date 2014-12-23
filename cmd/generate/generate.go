package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	kapi "github.com/GoogleCloudPlatform/kubernetes/pkg/api/v1beta2"
	kutil "github.com/GoogleCloudPlatform/kubernetes/pkg/util"

	buildapi "github.com/openshift/origin/pkg/build/api"
	deployapi "github.com/openshift/origin/pkg/deploy/api"
	imageapi "github.com/openshift/origin/pkg/image/api"
	routeapi "github.com/openshift/origin/pkg/route/api"

	"github.com/csrwng/origin-schema-generator/pkg/schemagen"
)

type Schema struct {
	PodList                   kapi.PodList
	ReplicationControllerList kapi.ReplicationControllerList
	ServiceList               kapi.ServiceList
	BuildList                 buildapi.BuildList
	BuildConfigList           buildapi.BuildConfigList
	ImageList                 imageapi.ImageList
	ImageRepositoryList       imageapi.ImageRepositoryList
	DeploymentList            deployapi.DeploymentList
	DeploymentConfigList      deployapi.DeploymentConfigList
	RouteList                 routeapi.RouteList
	ContainerStatus           kapi.ContainerStatus
}

func main() {
	packages := []schemagen.PackageDescriptor{
		{"github.com/GoogleCloudPlatform/kubernetes/pkg/api", "com.openshift.client.kubernetes", "kubernetes_"},
		{"github.com/GoogleCloudPlatform/kubernetes/pkg/util", "com.openshift.client.kubernetes.util", "kubernetes_util_"},
		{"github.com/fsouza/go-dockerclient", "com.openshift.client.dockerclient", "docker_"},
		{"github.com/openshift/origin/pkg/build/api", "com.openshift.client.openshift.build", "os_build_"},
		{"github.com/openshift/origin/pkg/deploy/api", "com.openshift.client.openshift.deploy", "os_deploy_"},
		{"github.com/openshift/origin/pkg/image/api", "com.openshift.client.openshift.image", "os_image_"},
		{"github.com/openshift/origin/pkg/route/api", "com.openshift.client.openshift.route", "os_route_"},
	}

	typeMap := map[reflect.Type]reflect.Type{
		reflect.TypeOf(kutil.Time{}): reflect.TypeOf(""),
		reflect.TypeOf(time.Time{}):  reflect.TypeOf(""),
		reflect.TypeOf(struct{}{}):   reflect.TypeOf(""),
	}
	schema, err := schemagen.GenerateSchema(reflect.TypeOf(Schema{}), packages, typeMap)
	if err != nil {
		fmt.Errorf("An error occurred: %v", err)
		return
	}
	b, _ := json.Marshal(&schema)
	result := string(b)
	result = strings.Replace(result, "\"additionalProperty\":", "\"additionalProperties\":", -1)
	fmt.Println(result)
}
