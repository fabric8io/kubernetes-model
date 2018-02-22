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

	authapi "github.com/openshift/origin/pkg/authorization/apis/authorization/v1"
	buildapi "github.com/openshift/origin/pkg/build/apis/build/v1"
	deployapi "github.com/openshift/origin/pkg/deploy/apis/apps/v1"
	imageapi "github.com/openshift/origin/pkg/image/apis/image/v1"
	oauthapi "github.com/openshift/origin/pkg/oauth/apis/oauth/v1"
	projectapi "github.com/openshift/origin/pkg/project/apis/project/v1"
	routeapi "github.com/openshift/origin/pkg/route/apis/route/v1"
	securityapi "github.com/openshift/origin/pkg/security/apis/security/v1"
	templateapi "github.com/openshift/origin/pkg/template/apis/template/v1"
	userapi "github.com/openshift/origin/pkg/user/apis/user/v1"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	resource "k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apimachineryversion "k8s.io/apimachinery/pkg/version"
	configapi "k8s.io/client-go/tools/clientcmd/api/v1"
	rapi "k8s.io/kubernetes/pkg/api/unversioned"
	kapi "k8s.io/kubernetes/pkg/api/v1"
	appsapi "k8s.io/kubernetes/pkg/apis/apps/v1beta1"
	authenticationapi "k8s.io/kubernetes/pkg/apis/authentication/v1"
	k8sauthapi "k8s.io/kubernetes/pkg/apis/authorization/v1"
	autoscalingapi "k8s.io/kubernetes/pkg/apis/autoscaling/v1"
	batchapiv1 "k8s.io/kubernetes/pkg/apis/batch/v1"
	batchapiv2alpha1 "k8s.io/kubernetes/pkg/apis/batch/v2alpha1"
	extensions "k8s.io/kubernetes/pkg/apis/extensions/v1beta1"
	policy "k8s.io/kubernetes/pkg/apis/policy/v1beta1"
	storageclassapi "k8s.io/kubernetes/pkg/apis/storage/v1"
	watch "k8s.io/kubernetes/pkg/watch/json"

	"github.com/fabric8io/kubernetes-model/pkg/schemagen"
	//"os"
	"os"
)

type Schema struct {
	Info                              apimachineryversion.Info
	BaseKubernetesList                kapi.List
	ObjectMeta                        metav1.ObjectMeta
	PodList                           kapi.PodList
	PodTemplateList                   kapi.PodTemplateList
	ReplicationControllerList         kapi.ReplicationControllerList
	ServiceList                       kapi.ServiceList
	Endpoints                         kapi.Endpoints
	EndpointsList                     kapi.EndpointsList
	EventList                         kapi.EventList
	Node                              kapi.Node
	NodeList                          kapi.NodeList
	EnvVar                            kapi.EnvVar
	Namespace                         kapi.Namespace
	NamespaceList                     kapi.NamespaceList
	PersistentVolume                  kapi.PersistentVolume
	PersistentVolumeList              kapi.PersistentVolumeList
	PersistentVolumeClaim             kapi.PersistentVolumeClaim
	PersistentVolumeClaimList         kapi.PersistentVolumeClaimList
	ResourceQuota                     kapi.ResourceQuota
	ResourceQuotaList                 kapi.ResourceQuotaList
	Secret                            kapi.Secret
	SecretList                        kapi.SecretList
	SecurityContextConstraints        securityapi.SecurityContextConstraints
	SecurityContextConstraintsList    securityapi.SecurityContextConstraintsList
	ServiceAccount                    kapi.ServiceAccount
	ServiceAccountList                kapi.ServiceAccountList
	Status                            metav1.Status
	Patch                             metav1.Patch
	Binding                           kapi.Binding
	LimitRangeList                    kapi.LimitRangeList
	DeleteOptions                     kapi.DeleteOptions
	Quantity                          resource.Quantity
	BuildRequest                      buildapi.BuildRequest
	BuildList                         buildapi.BuildList
	BuildConfigList                   buildapi.BuildConfigList
	ImageList                         imageapi.ImageList
	ImageStreamList                   imageapi.ImageStreamList
	ImageStreamTagList                imageapi.ImageStreamTagList
	DeploymentConfig                  deployapi.DeploymentConfig
	DeploymentConfigList              deployapi.DeploymentConfigList
	Route                             routeapi.Route
	RouteList                         routeapi.RouteList
	ComponentStatusList               kapi.ComponentStatusList
	ContainerStatus                   kapi.ContainerStatus
	Template                          templateapi.Template
	TemplateList                      templateapi.TemplateList
	TagEvent                          imageapi.TagEvent
	OAuthClient                       oauthapi.OAuthClient
	OAuthAccessToken                  oauthapi.OAuthAccessToken
	OAuthAuthorizeToken               oauthapi.OAuthAuthorizeToken
	OAuthClientAuthorization          oauthapi.OAuthClientAuthorization
	OAuthAccessTokenList              oauthapi.OAuthAccessTokenList
	OAuthAuthorizeTokenList           oauthapi.OAuthAuthorizeTokenList
	OAuthClientList                   oauthapi.OAuthClientList
	OAuthClientAuthorizationList      oauthapi.OAuthClientAuthorizationList
	TokenReview                       authenticationapi.TokenReview
	K8sSubjectAccessReview            k8sauthapi.SubjectAccessReview
	K8sLocalSubjectAccessReview       k8sauthapi.LocalSubjectAccessReview
	ClusterPolicy                     authapi.ClusterPolicy
	ClusterPolicyList                 authapi.ClusterPolicyList
	ClusterPolicyBinding              authapi.ClusterPolicyBinding
	ClusterPolicyBindingList          authapi.ClusterPolicyBindingList
	Policy                            authapi.Policy
	PolicyList                        authapi.PolicyList
	PolicyBinding                     authapi.PolicyBinding
	PolicyBindingList                 authapi.PolicyBindingList
	Role                              authapi.Role
	RoleList                          authapi.RoleList
	RoleBinding                       authapi.RoleBinding
	RoleBindingList                   authapi.RoleBindingList
	RoleBindingRestriction            authapi.RoleBindingRestriction
	LocalSubjectAccessReview          authapi.LocalSubjectAccessReview
	SubjectAccessReview               authapi.SubjectAccessReview
	SubjectAccessReviewResponse       authapi.SubjectAccessReviewResponse
	ClusterRoleBinding                authapi.ClusterRoleBinding
	ClusterRoleBindingList            authapi.ClusterRoleBindingList
	User                              userapi.User
	UserList                          userapi.UserList
	Group                             userapi.Group
	GroupList                         userapi.GroupList
	Identity                          userapi.Identity
	IdentityList                      userapi.IdentityList
	Config                            configapi.Config
	WatchEvent                        watch.WatchEvent
	RootPaths                         metav1.RootPaths
	Project                           projectapi.Project
	ProjectList                       projectapi.ProjectList
	ProjectRequest                    projectapi.ProjectRequest
	ListMeta                          rapi.ListMeta
	Job                               batchapiv1.Job
	JobList                           batchapiv1.JobList
	CronJob                           batchapiv2alpha1.CronJob
	CronJobList                       batchapiv2alpha1.CronJobList
	Scale                             extensions.Scale
	HorizontalPodAutoscaler           autoscalingapi.HorizontalPodAutoscaler
	HorizontalPodAutoscalerList       autoscalingapi.HorizontalPodAutoscalerList
	ThirdPartyResource                extensions.ThirdPartyResource
	ThirdPartyResourceList            extensions.ThirdPartyResourceList
	Deployment                        extensions.Deployment
	DeploymentList                    extensions.DeploymentList
	DeploymentRollback                extensions.DeploymentRollback
	PodSecurityPolicy                 extensions.PodSecurityPolicy
	PodSecurityPolicyList             extensions.PodSecurityPolicyList
	PodDisruptionBudget               policy.PodDisruptionBudget
	PodDisruptionBudgetList           policy.PodDisruptionBudgetList
	StatefulSet                       appsapi.StatefulSet
	StatefulSetList                   appsapi.StatefulSetList
	DaemonSet                         extensions.DaemonSet
	DaemonSetList                     extensions.DaemonSetList
	Ingress                           extensions.Ingress
	IngressList                       extensions.IngressList
	ReplicaSet                        extensions.ReplicaSet
	ReplicaSetList                    extensions.ReplicaSetList
	NetworkPolicy                     extensions.NetworkPolicy
	NetworkPolicyList                 extensions.NetworkPolicyList
	ConfigMap                         kapi.ConfigMap
	ConfigMapList                     kapi.ConfigMapList
	Toleration                        kapi.Toleration
	CustomResourceDefinition          apiextensions.CustomResourceDefinition
	CustomResourceDefinitionList      apiextensions.CustomResourceDefinitionList
	CustomResourceDefinitionSpec      apiextensions.CustomResourceDefinitionSpec
	CustomResourceDefinitionNames     apiextensions.CustomResourceDefinitionNames
	CustomResourceDefinitionCondition apiextensions.CustomResourceDefinitionCondition
	CustomResourceDefinitionStatus    apiextensions.CustomResourceDefinitionStatus
	StorageClass                      storageclassapi.StorageClass
	StorageClassList                  storageclassapi.StorageClassList
}

func main() {
	packages := []schemagen.PackageDescriptor{
		{"k8s.io/kubernetes/pkg/api", "", "io.fabric8.kubernetes.api.model", "kubernetes_"},
		{"k8s.io/kubernetes/pkg/api/v1", "", "io.fabric8.kubernetes.api.model", "kubernetes_"},
		{"k8s.io/apimachinery/pkg/api/resource", "", "io.fabric8.kubernetes.api.model", "kubernetes_resource_"},
		{"k8s.io/apimachinery/pkg/util/intstr", "", "io.fabric8.kubernetes.api.model", "k8s_io_apimachinery_pkg_util_intstr_"},
		{"k8s.io/apimachinery/pkg/runtime", "", "io.fabric8.kubernetes.api.model.runtime", "k8s_io_apimachinery_pkg_runtime_"},
		{"k8s.io/apimachinery/pkg/version", "", "io.fabric8.kubernetes.api.model.version", "k8s_io_apimachinery_pkg_version_"},
		{"k8s.io/kubernetes/pkg/util", "", "io.fabric8.kubernetes.api.model", "kubernetes_util_"},
		{"k8s.io/kubernetes/pkg/watch/json", "", "io.fabric8.kubernetes.api.model", "kubernetes_watch_"},
		{"k8s.io/kubernetes/pkg/api/errors", "", "io.fabric8.kubernetes.api.model", "kubernetes_errors_"},
		{"k8s.io/client-go/tools/clientcmd/api/v1", "", "io.fabric8.kubernetes.api.model", "kubernetes_config_"},
		{"github.com/openshift/origin/pkg/build/apis/build/v1", "", "io.fabric8.openshift.api.model", "os_build_"},
		{"github.com/openshift/origin/pkg/deploy/apis/apps/v1", "", "io.fabric8.openshift.api.model", "os_deploy_"},
		{"github.com/openshift/origin/pkg/image/apis/image/v1", "", "io.fabric8.openshift.api.model", "os_image_"},
		{"github.com/openshift/origin/pkg/oauth/apis/oauth/v1", "", "io.fabric8.openshift.api.model", "os_oauth_"},
		{"github.com/openshift/origin/pkg/route/apis/route/v1", "", "io.fabric8.openshift.api.model", "os_route_"},
		{"github.com/openshift/origin/pkg/template/apis/template/v1", "", "io.fabric8.openshift.api.model", "os_template_"},
		{"github.com/openshift/origin/pkg/user/apis/user/v1", "", "io.fabric8.openshift.api.model", "os_user_"},
		{"github.com/openshift/origin/pkg/authorization/apis/authorization/v1", "", "io.fabric8.openshift.api.model", "os_authorization_"},
		{"github.com/openshift/origin/pkg/project/apis/project/v1", "", "io.fabric8.openshift.api.model", "os_project_"},
		{"github.com/openshift/origin/pkg/security/apis/security/v1", "", "io.fabric8.openshift.api.model", "os_security_"},
		{"k8s.io/kubernetes/pkg/api/unversioned", "", "io.fabric8.kubernetes.api.model", "api_"},
		{"k8s.io/kubernetes/pkg/apis/extensions/v1beta1", "", "io.fabric8.kubernetes.api.model.extensions", "kubernetes_extensions_"},
		{"k8s.io/kubernetes/pkg/apis/policy/v1beta1", "", "io.fabric8.kubernetes.api.model.policy", "kubernetes_policy_"},
		{"k8s.io/kubernetes/pkg/apis/authentication/v1", "authentication.k8s.io", "io.fabric8.kubernetes.api.model.authentication", "kubernetes_authentication_"},
		{"k8s.io/kubernetes/pkg/apis/authorization/v1", "authorization.k8s.io", "io.fabric8.kubernetes.api.model.authorization", "kubernetes_authorization_"},
		{"k8s.io/kubernetes/pkg/apis/apps/v1beta1", "", "io.fabric8.kubernetes.api.model.extensions", "kubernetes_apps_"},
		{"k8s.io/kubernetes/pkg/apis/batch/v2alpha1", "", "io.fabric8.kubernetes.api.model", "kubernetes_batch_"},
		{"k8s.io/kubernetes/pkg/apis/batch/v1", "", "io.fabric8.kubernetes.api.model", "kubernetes_batch_"},
		{"k8s.io/kubernetes/pkg/apis/autoscaling/v1", "", "io.fabric8.kubernetes.api.model", "kubernetes_autoscaling_"},
		{"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1", "", "io.fabric8.kubernetes.api.model.apiextensions", "k8s_io_apiextensions_"},
		{"k8s.io/apimachinery/pkg/apis/meta/v1", "", "io.fabric8.kubernetes.api.model", "k8s_io_apimachinery_"},
		{"k8s.io/kubernetes/pkg/apis/storage/v1", "", "io.fabric8.kubernetes.api.model", "kubernetes_storageclass_"},
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

	args := os.Args[1:]
	if len(args) < 1 || args[0] != "validation" {
		schema.Resources = nil
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
