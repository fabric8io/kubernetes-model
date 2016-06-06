/**
 * Copyright (C) 2011 Red Hat, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *         http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	authapi "github.com/openshift/origin/pkg/authorization/api/v1"
	buildapi "github.com/openshift/origin/pkg/build/api/v1"
	deployapi "github.com/openshift/origin/pkg/deploy/api/v1"
	imageapi "github.com/openshift/origin/pkg/image/api/v1"
	oauthapi "github.com/openshift/origin/pkg/oauth/api/v1"
	projectapi "github.com/openshift/origin/pkg/project/api/v1"
	routeapi "github.com/openshift/origin/pkg/route/api/v1"
	templateapi "github.com/openshift/origin/pkg/template/api/v1"
	userapi "github.com/openshift/origin/pkg/user/api/v1"
	resourceapi "k8s.io/kubernetes/pkg/api/resource"
	rapi "k8s.io/kubernetes/pkg/api/unversioned"
	kapi "k8s.io/kubernetes/pkg/api/v1"
	extensions "k8s.io/kubernetes/pkg/apis/extensions/v1beta1"
	configapi "k8s.io/kubernetes/pkg/client/unversioned/clientcmd/api/v1"
	watch "k8s.io/kubernetes/pkg/watch/json"

	"github.com/fabric8io/kubernetes-model/pkg/schemagen"
)

type Schema struct {
	BaseKubernetesList             kapi.List
	ObjectMeta                     kapi.ObjectMeta
	PodList                        kapi.PodList
	PodTemplateList                kapi.PodTemplateList
	ReplicationControllerList      kapi.ReplicationControllerList
	ServiceList                    kapi.ServiceList
	Endpoints                      kapi.Endpoints
	EndpointsList                  kapi.EndpointsList
	EventList                      kapi.EventList
	Node                           kapi.Node
	NodeList                       kapi.NodeList
	EnvVar                         kapi.EnvVar
	Namespace                      kapi.Namespace
	NamespaceList                  kapi.NamespaceList
	PersistentVolume               kapi.PersistentVolume
	PersistentVolumeList           kapi.PersistentVolumeList
	PersistentVolumeClaim          kapi.PersistentVolumeClaim
	PersistentVolumeClaimList      kapi.PersistentVolumeClaimList
	ResourceQuota                  kapi.ResourceQuota
	ResourceQuotaList              kapi.ResourceQuotaList
	Secret                         kapi.Secret
	SecretList                     kapi.SecretList
	SecurityContextConstraints     kapi.SecurityContextConstraints
	SecurityContextConstraintsList kapi.SecurityContextConstraintsList
	ServiceAccount                 kapi.ServiceAccount
	ServiceAccountList             kapi.ServiceAccountList
	Status                         rapi.Status
	Patch                          rapi.Patch
	Binding                        kapi.Binding
	LimitRangeList                 kapi.LimitRangeList
	DeleteOptions                  kapi.DeleteOptions
	Quantity                       resourceapi.Quantity
	BuildRequest                   buildapi.BuildRequest
	BuildList                      buildapi.BuildList
	BuildConfigList                buildapi.BuildConfigList
	ImageList                      imageapi.ImageList
	ImageStreamList                imageapi.ImageStreamList
	DeploymentConfigList           deployapi.DeploymentConfigList
	RouteList                      routeapi.RouteList
	ComponentStatusList            kapi.ComponentStatusList
	ContainerStatus                kapi.ContainerStatus
	Template                       templateapi.Template
	TemplateList                   templateapi.TemplateList
	TagEvent                       imageapi.TagEvent
	OAuthClient                    oauthapi.OAuthClient
	OAuthAccessToken               oauthapi.OAuthAccessToken
	OAuthAuthorizeToken            oauthapi.OAuthAuthorizeToken
	OAuthClientAuthorization       oauthapi.OAuthClientAuthorization
	OAuthAccessTokenList           oauthapi.OAuthAccessTokenList
	OAuthAuthorizeTokenList        oauthapi.OAuthAuthorizeTokenList
	OAuthClientList                oauthapi.OAuthClientList
	OAuthClientAuthorizationList   oauthapi.OAuthClientAuthorizationList
	ClusterPolicy                  authapi.ClusterPolicy
	ClusterPolicyList              authapi.ClusterPolicyList
	ClusterPolicyBinding           authapi.ClusterPolicyBinding
	ClusterPolicyBindingList       authapi.ClusterPolicyBindingList
	Policy                         authapi.Policy
	PolicyList                     authapi.PolicyList
	PolicyBinding                  authapi.PolicyBinding
	PolicyBindingList              authapi.PolicyBindingList
	Role                           authapi.Role
	RoleList                       authapi.RoleList
	RoleBinding                    authapi.RoleBinding
	RoleBindingList                authapi.RoleBindingList
	LocalSubjectAccessReview       authapi.LocalSubjectAccessReview
	SubjectAccessReview            authapi.SubjectAccessReview
	SubjectAccessReviewResponse    authapi.SubjectAccessReviewResponse
	ClusterRoleBinding             authapi.ClusterRoleBinding
	ClusterRoleBindingList         authapi.ClusterRoleBindingList
	User                           userapi.User
	UserList                       userapi.UserList
	Group                          userapi.Group
	GroupList                      userapi.GroupList
	Identity                       userapi.Identity
	IdentityList                   userapi.IdentityList
	Config                         configapi.Config
	WatchEvent                     watch.WatchEvent
	RootPaths                      rapi.RootPaths
	Project                        projectapi.Project
	ProjectList                    projectapi.ProjectList
	ProjectRequest                 projectapi.ProjectRequest
	ListMeta                       rapi.ListMeta
	Job                            extensions.Job
	JobList                        extensions.JobList
	Scale                          extensions.Scale
	HorizontalPodAutoscaler        extensions.HorizontalPodAutoscaler
	HorizontalPodAutoscalerList    extensions.HorizontalPodAutoscalerList
	ThirdPartyResource             extensions.ThirdPartyResource
	ThirdPartyResourceList         extensions.ThirdPartyResourceList
	Deployment                     extensions.Deployment
	DeploymentList                 extensions.DeploymentList
	DeploymentRollback             extensions.DeploymentRollback
	DaemonSet                      extensions.DaemonSet
	DaemonSetList                  extensions.DaemonSetList
	Ingress                        extensions.Ingress
	IngressList                    extensions.IngressList
	ReplicaSet                     extensions.ReplicaSet
	ReplicaSetList                 extensions.ReplicaSetList
	ConfigMap                      kapi.ConfigMap
	ConfigMapList                  kapi.ConfigMapList
}

func main() {
	packages := []schemagen.PackageDescriptor{
		{"k8s.io/kubernetes/pkg/api/v1", "io.fabric8.kubernetes.api.model", "kubernetes_"},
		{"k8s.io/kubernetes/pkg/api/resource", "io.fabric8.kubernetes.api.model", "kubernetes_resource_"},
		{"k8s.io/kubernetes/pkg/util/intstr", "io.fabric8.kubernetes.api.model", "kubernetes_intstr_"},
		{"k8s.io/kubernetes/pkg/runtime", "io.fabric8.kubernetes.api.model.runtime", "kubernetes_runtime_"},
		{"k8s.io/kubernetes/pkg/util", "io.fabric8.kubernetes.api.model", "kubernetes_util_"},
		{"k8s.io/kubernetes/pkg/watch/json", "io.fabric8.kubernetes.api.model", "kubernetes_watch_"},
		{"k8s.io/kubernetes/pkg/api/errors", "io.fabric8.kubernetes.api.model", "kubernetes_errors_"},
		{"k8s.io/kubernetes/pkg/client/clientcmd/api/v1", "io.fabric8.kubernetes.api.model", "kubernetes_config_"},
		{"speter.net/go/exp/math/dec/inf", "io.fabric8.openshift.api.model", "speter_inf_"},
		{"github.com/openshift/origin/pkg/build/api/v1", "io.fabric8.openshift.api.model", "os_build_"},
		{"github.com/openshift/origin/pkg/deploy/api/v1", "io.fabric8.openshift.api.model", "os_deploy_"},
		{"github.com/openshift/origin/pkg/image/api/v1", "io.fabric8.openshift.api.model", "os_image_"},
		{"github.com/openshift/origin/pkg/oauth/api/v1", "io.fabric8.openshift.api.model", "os_oauth_"},
		{"github.com/openshift/origin/pkg/route/api/v1", "io.fabric8.openshift.api.model", "os_route_"},
		{"github.com/openshift/origin/pkg/template/api/v1", "io.fabric8.openshift.api.model", "os_template_"},
		{"github.com/openshift/origin/pkg/user/api/v1", "io.fabric8.openshift.api.model", "os_user_"},
		{"github.com/openshift/origin/pkg/authorization/api/v1", "io.fabric8.openshift.api.model", "os_authorization_"},
		{"github.com/openshift/origin/pkg/project/api/v1", "io.fabric8.openshift.api.model", "os_project_"},
		{"k8s.io/kubernetes/pkg/api/unversioned", "io.fabric8.kubernetes.api.model", "api_"},
		{"k8s.io/kubernetes/pkg/client/unversioned/clientcmd/api/v1", "io.fabric8.kubernetes.api.model", "clientcmd_api_"},
		{"k8s.io/kubernetes/pkg/apis/extensions/v1beta1", "io.fabric8.kubernetes.api.model.extensions", "kubernetes_extensions_"},
	}

	typeMap := map[reflect.Type]reflect.Type{
		reflect.TypeOf(rapi.Time{}): reflect.TypeOf(""),
		reflect.TypeOf(time.Time{}): reflect.TypeOf(""),
		reflect.TypeOf(struct{}{}):  reflect.TypeOf(""),
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
