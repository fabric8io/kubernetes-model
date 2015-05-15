package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	kerrors "github.com/GoogleCloudPlatform/kubernetes/pkg/api/errors"
	resourceapi "github.com/GoogleCloudPlatform/kubernetes/pkg/api/resource"
	kapi "github.com/GoogleCloudPlatform/kubernetes/pkg/api/v1beta3"
	kutil "github.com/GoogleCloudPlatform/kubernetes/pkg/util"
	buildapi "github.com/openshift/origin/pkg/build/api/v1beta3"
	deployapi "github.com/openshift/origin/pkg/deploy/api/v1beta3"
	imageapi "github.com/openshift/origin/pkg/image/api/v1beta3"
	oauthapi "github.com/openshift/origin/pkg/oauth/api/v1beta3"
	routeapi "github.com/openshift/origin/pkg/route/api/v1beta3"
	templateapi "github.com/openshift/origin/pkg/template/api/v1beta3"

	"github.com/fabric8io/origin-schema-generator/pkg/schemagen"
)

type Schema struct {
	ObjectMeta                   kapi.ObjectMeta
	PodList                      kapi.PodList
	ReplicationControllerList    kapi.ReplicationControllerList
	ServiceList                  kapi.ServiceList
	Endpoints                    kapi.Endpoints
	EndpointsList                kapi.EndpointsList
	Node                         kapi.Node
	NodeList                     kapi.NodeList
	BaseKubernetesList           kapi.List
	EnvVar                       kapi.EnvVar
	Namespace                    kapi.Namespace
	NamespaceList                kapi.NamespaceList
	Secret                       kapi.Secret
	Quantity                     resourceapi.Quantity
	StatusError                  kerrors.StatusError
	BuildList                    buildapi.BuildList
	BuildConfigList              buildapi.BuildConfigList
	ImageList                    imageapi.ImageList
	ImageStreamList              imageapi.ImageStreamList
	DeploymentConfigList         deployapi.DeploymentConfigList
	RouteList                    routeapi.RouteList
	ContainerStatus              kapi.ContainerStatus
	Template                     templateapi.Template
	TagEvent                     imageapi.TagEvent
	OAuthClient                  oauthapi.OAuthClient
	OAuthAccessToken             oauthapi.OAuthAccessToken
	OAuthAuthorizeToken          oauthapi.OAuthAuthorizeToken
	OAuthClientAuthorization     oauthapi.OAuthClientAuthorization
	OAuthAccessTokenList         oauthapi.OAuthAccessTokenList
	OAuthAuthorizeTokenList      oauthapi.OAuthAuthorizeTokenList
	OAuthClientList              oauthapi.OAuthClientList
	OAuthClientAuthorizationList oauthapi.OAuthClientAuthorizationList
}

func main() {
	packages := []schemagen.PackageDescriptor{
		{"github.com/GoogleCloudPlatform/kubernetes/pkg/api/v1beta3", "io.fabric8.kubernetes.api.model", "kubernetes_"},
		{"github.com/GoogleCloudPlatform/kubernetes/pkg/api/v1beta3", "io.fabric8.kubernetes.api.model", "kubernetes_"},
		{"github.com/GoogleCloudPlatform/kubernetes/pkg/api/v1beta3", "io.fabric8.kubernetes.api.model", "kubernetes_"},
		{"github.com/GoogleCloudPlatform/kubernetes/pkg/api/resource", "io.fabric8.kubernetes.api.model.resource", "kubernetes_resource_"},
		{"github.com/GoogleCloudPlatform/kubernetes/pkg/runtime", "io.fabric8.kubernetes.api.model.runtime", "kubernetes_runtime_"},
		{"github.com/GoogleCloudPlatform/kubernetes/pkg/util", "io.fabric8.kubernetes.api.model.util", "kubernetes_util_"},
		{"github.com/GoogleCloudPlatform/kubernetes/pkg/api/errors", "io.fabric8.kubernetes.api.model.errors", "kubernetes_errors_"},
		{"github.com/GoogleCloudPlatform/kubernetes/pkg/api", "io.fabric8.kubernetes.api.model.base", "kubernetes_base_"},
		{"github.com/fsouza/go-dockerclient", "io.fabric8.docker.client.dockerclient", "docker_"},
		{"speter.net/go/exp/math/dec/inf", "io.fabric8.openshift.client.util", "speter_inf_"},
		{"github.com/openshift/origin/pkg/build/api/v1beta3", "io.fabric8.openshift.api.model", "os_build_"},
		{"github.com/openshift/origin/pkg/deploy/api/v1beta3", "io.fabric8.openshift.api.model", "os_deploy_"},
		{"github.com/openshift/origin/pkg/image/api/v1beta3", "io.fabric8.openshift.api.model", "os_image_"},
		{"github.com/openshift/origin/pkg/oauth/api/v1beta3", "io.fabric8.openshift.api.model", "os_oauth_"},
		{"github.com/openshift/origin/pkg/route/api/v1beta3", "io.fabric8.openshift.api.model", "os_route_"},
		{"github.com/openshift/origin/pkg/template/api/v1beta3", "io.fabric8.openshift.api.model.template", "os_template_"},
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
