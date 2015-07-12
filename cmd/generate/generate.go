package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	kerrors "github.com/GoogleCloudPlatform/kubernetes/pkg/api/errors"
	resourceapi "github.com/GoogleCloudPlatform/kubernetes/pkg/api/resource"
	kapi "github.com/GoogleCloudPlatform/kubernetes/pkg/api/v1"
	configapi "github.com/GoogleCloudPlatform/kubernetes/pkg/client/clientcmd/api/v1"
	kutil "github.com/GoogleCloudPlatform/kubernetes/pkg/util"
	watch "github.com/GoogleCloudPlatform/kubernetes/pkg/watch/json"
	buildapi "github.com/openshift/origin/pkg/build/api/v1"
	deployapi "github.com/openshift/origin/pkg/deploy/api/v1"
	imageapi "github.com/openshift/origin/pkg/image/api/v1"
	oauthapi "github.com/openshift/origin/pkg/oauth/api/v1"
	routeapi "github.com/openshift/origin/pkg/route/api/v1"
	templateapi "github.com/openshift/origin/pkg/template/api/v1"

	"github.com/fabric8io/origin-schema-generator/pkg/schemagen"
)

type Schema struct {
	ObjectMeta                   kapi.ObjectMeta
	PodList                      kapi.PodList
	ReplicationControllerList    kapi.ReplicationControllerList
	ServiceList                  kapi.ServiceList
	Endpoints                    kapi.Endpoints
	EndpointsList                kapi.EndpointsList
	EventList                    kapi.EventList
	Node                         kapi.Node
	NodeList                     kapi.NodeList
	BaseKubernetesList           kapi.List
	EnvVar                       kapi.EnvVar
	Namespace                    kapi.Namespace
	NamespaceList                kapi.NamespaceList
	Secret                       kapi.Secret
	SecretList                   kapi.SecretList
	ServiceAccount               kapi.ServiceAccount
	ServiceAccountList           kapi.ServiceAccountList
	Quantity                     resourceapi.Quantity
	StatusError                  kerrors.StatusError
	BuildRequest                 buildapi.BuildRequest
	BuildList                    buildapi.BuildList
	BuildConfigList              buildapi.BuildConfigList
	ImageList                    imageapi.ImageList
	ImageStreamList              imageapi.ImageStreamList
	DeploymentConfigList         deployapi.DeploymentConfigList
	RouteList                    routeapi.RouteList
	ContainerStatus              kapi.ContainerStatus
	Template                     templateapi.Template
	TemplateList                 templateapi.TemplateList
	TagEvent                     imageapi.TagEvent
	OAuthClient                  oauthapi.OAuthClient
	OAuthAccessToken             oauthapi.OAuthAccessToken
	OAuthAuthorizeToken          oauthapi.OAuthAuthorizeToken
	OAuthClientAuthorization     oauthapi.OAuthClientAuthorization
	OAuthAccessTokenList         oauthapi.OAuthAccessTokenList
	OAuthAuthorizeTokenList      oauthapi.OAuthAuthorizeTokenList
	OAuthClientList              oauthapi.OAuthClientList
	OAuthClientAuthorizationList oauthapi.OAuthClientAuthorizationList
	Config                       configapi.Config
	WatchEvent                   watch.WatchEvent
}

func main() {
	packages := []schemagen.PackageDescriptor{
		{"github.com/GoogleCloudPlatform/kubernetes/pkg/api/v1", "io.fabric8.kubernetes.api.model", "kubernetes_"},
		{"github.com/GoogleCloudPlatform/kubernetes/pkg/api/resource", "io.fabric8.kubernetes.api.model.resource", "kubernetes_resource_"},
		{"github.com/GoogleCloudPlatform/kubernetes/pkg/runtime", "io.fabric8.kubernetes.api.model.runtime", "kubernetes_runtime_"},
		{"github.com/GoogleCloudPlatform/kubernetes/pkg/util", "io.fabric8.kubernetes.api.model.util", "kubernetes_util_"},
		{"github.com/GoogleCloudPlatform/kubernetes/pkg/watch/json", "io.fabric8.kubernetes.api.watch", "kubernetes_watch_"},
		{"github.com/GoogleCloudPlatform/kubernetes/pkg/api/errors", "io.fabric8.kubernetes.api.model.errors", "kubernetes_errors_"},
		{"github.com/GoogleCloudPlatform/kubernetes/pkg/api", "io.fabric8.kubernetes.api.model.base", "kubernetes_base_"},
		{"github.com/GoogleCloudPlatform/kubernetes/pkg/client/clientcmd/api/v1", "io.fabric8.kubernetes.api.model.config", "kubernetes_config_"},
		{"github.com/fsouza/go-dockerclient", "io.fabric8.docker.client.dockerclient", "docker_"},
		{"speter.net/go/exp/math/dec/inf", "io.fabric8.openshift.client.util", "speter_inf_"},
		{"github.com/openshift/origin/pkg/build/api/v1", "io.fabric8.openshift.api.model", "os_build_"},
		{"github.com/openshift/origin/pkg/deploy/api/v1", "io.fabric8.openshift.api.model", "os_deploy_"},
		{"github.com/openshift/origin/pkg/image/api/v1", "io.fabric8.openshift.api.model", "os_image_"},
		{"github.com/openshift/origin/pkg/oauth/api/v1", "io.fabric8.openshift.api.model", "os_oauth_"},
		{"github.com/openshift/origin/pkg/route/api/v1", "io.fabric8.openshift.api.model", "os_route_"},
		{"github.com/openshift/origin/pkg/template/api/v1", "io.fabric8.openshift.api.model.template", "os_template_"},
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

	b, err := json.Marshal(&schema)
	if err != nil {
		log.Fatal(err)
	}
	result := string(b)
	result = strings.Replace(result, "\"additionalProperty\":", "\"additionalProperties\":", -1)
	var out bytes.Buffer
	err = json.Indent(&out, []byte(result), "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(out.String())
}
