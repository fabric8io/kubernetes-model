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

	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	resource "k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	configapi "k8s.io/client-go/tools/clientcmd/api/v1"
	rapi "k8s.io/kubernetes/pkg/api/unversioned"
	kapi "k8s.io/kubernetes/pkg/api/v1"
	appsapi "k8s.io/kubernetes/pkg/apis/apps/v1beta1"
	authenticationapi "k8s.io/kubernetes/pkg/apis/authentication/v1"
	autoscalingapi "k8s.io/kubernetes/pkg/apis/autoscaling/v1"
	batchapiv1 "k8s.io/kubernetes/pkg/apis/batch/v1"
	batchapiv2alpha1 "k8s.io/kubernetes/pkg/apis/batch/v2alpha1"
	extensions "k8s.io/kubernetes/pkg/apis/extensions/v1beta1"
	storageclassapi "k8s.io/kubernetes/pkg/apis/storage/v1"
	certificatesapi "k8s.io/kubernetes/pkg/apis/certificates/v1beta1"
	watch "k8s.io/kubernetes/pkg/watch/json"
	rbac "k8s.io/kubernetes/pkg/apis/rbac/v1beta1"

	"github.com/fabric8io/kubernetes-model/pkg/schemagen"
	//"os"
	"os"
)

type Schema struct {
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
	ServiceAccount                    kapi.ServiceAccount
	ServiceAccountList                kapi.ServiceAccountList
	Status                            metav1.Status
	Patch                             metav1.Patch
	Binding                           kapi.Binding
	LimitRangeList                    kapi.LimitRangeList
	DeleteOptions                     kapi.DeleteOptions
	Quantity                          resource.Quantity
	ComponentStatusList               kapi.ComponentStatusList
	ContainerStatus                   kapi.ContainerStatus
	TokenReview                       authenticationapi.TokenReview
	Role                              rbac.Role
	RoleList                          rbac.RoleList
	RoleBinding                       rbac.RoleBinding
	RoleBindingList                   rbac.RoleBindingList
  ClusterRole                       rbac.ClusterRole
  ClusterRoleList                   rbac.ClusterRoleList
  ClusterRoleBinding                rbac.ClusterRoleBinding
  ClusterRoleBindingList            rbac.ClusterRoleBindingList
	Config                            configapi.Config
	WatchEvent                        watch.WatchEvent
	RootPaths                         metav1.RootPaths
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
	CertificateSigningRequest         certificatesapi.CertificateSigningRequest
	CertificateSigningRequestList     certificatesapi.CertificateSigningRequestList
	StorageClass                      storageclassapi.StorageClass
	StorageClassList                  storageclassapi.StorageClassList
}

func main() {
	packages := []schemagen.PackageDescriptor{
		{"k8s.io/kubernetes/pkg/api", "io.ucosty.kubernetes.api.model", "kubernetes_"},
		{"k8s.io/kubernetes/pkg/api/v1", "io.ucosty.kubernetes.api.model", "kubernetes_"},
		{"k8s.io/apimachinery/pkg/api/resource", "io.ucosty.kubernetes.api.model", "kubernetes_resource_"},
		{"k8s.io/apimachinery/pkg/util/intstr", "io.ucosty.kubernetes.api.model", "k8s_io_apimachinery_pkg_util_intstr_"},
		{"k8s.io/apimachinery/pkg/runtime", "io.ucosty.kubernetes.api.model.runtime", "k8s_io_apimachinery_pkg_runtime_"},
		{"k8s.io/kubernetes/pkg/util", "io.ucosty.kubernetes.api.model", "kubernetes_util_"},
		{"k8s.io/kubernetes/pkg/watch/json", "io.ucosty.kubernetes.api.model", "kubernetes_watch_"},
		{"k8s.io/kubernetes/pkg/api/errors", "io.ucosty.kubernetes.api.model", "kubernetes_errors_"},
		{"k8s.io/client-go/tools/clientcmd/api/v1", "io.ucosty.kubernetes.api.model", "kubernetes_config_"},
		{"k8s.io/kubernetes/pkg/api/unversioned", "io.ucosty.kubernetes.api.model", "api_"},
		{"k8s.io/kubernetes/pkg/apis/extensions/v1beta1", "io.ucosty.kubernetes.api.model.extensions", "kubernetes_extensions_"},
		{"k8s.io/kubernetes/pkg/apis/authentication/v1", "io.ucosty.kubernetes.api.model.authentication", "kubernetes_authentication_"},
		{"k8s.io/kubernetes/pkg/apis/certificates/v1beta1", "io.ucosty.kubernetes.api.model", "kubernetes_certificates_"},
		{"k8s.io/kubernetes/pkg/apis/apps/v1beta1", "io.ucosty.kubernetes.api.model.extensions", "kubernetes_apps_"},
		{"k8s.io/kubernetes/pkg/apis/batch/v2alpha1", "io.ucosty.kubernetes.api.model", "kubernetes_batch_"},
		{"k8s.io/kubernetes/pkg/apis/batch/v1", "io.ucosty.kubernetes.api.model", "kubernetes_batch_"},
		{"k8s.io/kubernetes/pkg/apis/autoscaling/v1", "io.ucosty.kubernetes.api.model", "kubernetes_autoscaling_"},
		{"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1", "io.ucosty.kubernetes.api.model.apiextensions", "k8s_io_apiextensions_"},
		{"k8s.io/apimachinery/pkg/apis/meta/v1", "io.ucosty.kubernetes.api.model", "k8s_io_apimachinery_"},
    {"k8s.io/kubernetes/pkg/apis/storage/v1", "io.ucosty.kubernetes.api.model", "kubernetes_storageclass_"},
		{"k8s.io/kubernetes/pkg/apis/rbac/v1beta1", "io.ucosty.kubernetes.api.model", "kubernetes_rbac_"},
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
