/**
 * Copyright (C) 2015 Red Hat, Inc.
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

package io.ucosty.kubernetes.client;

import io.ucosty.kubernetes.api.model.ComponentStatus;
import io.ucosty.kubernetes.api.model.ComponentStatusList;
import io.ucosty.kubernetes.api.model.ConfigMap;
import io.ucosty.kubernetes.api.model.ConfigMapList;
import io.ucosty.kubernetes.api.model.Doneable;
import io.ucosty.kubernetes.api.model.DoneableComponentStatus;
import io.ucosty.kubernetes.api.model.DoneableConfigMap;
import io.ucosty.kubernetes.api.model.DoneableEndpoints;
import io.ucosty.kubernetes.api.model.DoneableEvent;
import io.ucosty.kubernetes.api.model.DoneableLimitRange;
import io.ucosty.kubernetes.api.model.DoneableNamespace;
import io.ucosty.kubernetes.api.model.DoneableNode;
import io.ucosty.kubernetes.api.model.DoneablePersistentVolume;
import io.ucosty.kubernetes.api.model.DoneablePersistentVolumeClaim;
import io.ucosty.kubernetes.api.model.DoneablePod;
import io.ucosty.kubernetes.api.model.DoneableReplicationController;
import io.ucosty.kubernetes.api.model.DoneableResourceQuota;
import io.ucosty.kubernetes.api.model.DoneableSecret;
import io.ucosty.kubernetes.api.model.DoneableService;
import io.ucosty.kubernetes.api.model.DoneableServiceAccount;
import io.ucosty.kubernetes.api.model.Endpoints;
import io.ucosty.kubernetes.api.model.EndpointsList;
import io.ucosty.kubernetes.api.model.Event;
import io.ucosty.kubernetes.api.model.EventList;
import io.ucosty.kubernetes.api.model.HasMetadata;
import io.ucosty.kubernetes.api.model.KubernetesResourceList;
import io.ucosty.kubernetes.api.model.LimitRange;
import io.ucosty.kubernetes.api.model.LimitRangeList;
import io.ucosty.kubernetes.api.model.Namespace;
import io.ucosty.kubernetes.api.model.NamespaceList;
import io.ucosty.kubernetes.api.model.Node;
import io.ucosty.kubernetes.api.model.NodeList;
import io.ucosty.kubernetes.api.model.PersistentVolume;
import io.ucosty.kubernetes.api.model.PersistentVolumeClaim;
import io.ucosty.kubernetes.api.model.PersistentVolumeClaimList;
import io.ucosty.kubernetes.api.model.PersistentVolumeList;
import io.ucosty.kubernetes.api.model.Pod;
import io.ucosty.kubernetes.api.model.PodList;
import io.ucosty.kubernetes.api.model.ReplicationController;
import io.ucosty.kubernetes.api.model.ReplicationControllerList;
import io.ucosty.kubernetes.api.model.ResourceQuota;
import io.ucosty.kubernetes.api.model.ResourceQuotaList;
import io.ucosty.kubernetes.api.model.Secret;
import io.ucosty.kubernetes.api.model.SecretList;
import io.ucosty.kubernetes.api.model.Service;
import io.ucosty.kubernetes.api.model.ServiceAccount;
import io.ucosty.kubernetes.api.model.ServiceAccountList;
import io.ucosty.kubernetes.api.model.ServiceList;
import io.ucosty.kubernetes.api.model.apiextensions.CustomResourceDefinition;
import io.ucosty.kubernetes.api.model.apiextensions.CustomResourceDefinitionList;
import io.ucosty.kubernetes.api.model.apiextensions.DoneableCustomResourceDefinition;
import io.ucosty.kubernetes.client.dsl.AppsAPIGroupDSL;
import io.ucosty.kubernetes.client.dsl.AutoscalingAPIGroupDSL;
import io.ucosty.kubernetes.client.dsl.ExtensionsAPIGroupDSL;
import io.ucosty.kubernetes.client.dsl.KubernetesListMixedOperation;
import io.ucosty.kubernetes.client.dsl.MixedOperation;
import io.ucosty.kubernetes.client.dsl.NamespaceListVisitFromServerGetDeleteRecreateWaitApplicable;
import io.ucosty.kubernetes.client.dsl.NamespaceVisitFromServerGetWatchDeleteRecreateWaitApplicable;
import io.ucosty.kubernetes.client.dsl.NonNamespaceOperation;
import io.ucosty.kubernetes.client.dsl.ParameterNamespaceListVisitFromServerGetDeleteRecreateWaitApplicable;
import io.ucosty.kubernetes.client.dsl.PodResource;
import io.ucosty.kubernetes.client.dsl.Resource;
import io.ucosty.kubernetes.client.dsl.RollableScalableResource;

import java.io.InputStream;
import java.util.Collection;

public interface KubernetesClient extends Client {

  NonNamespaceOperation<CustomResourceDefinition, CustomResourceDefinitionList, DoneableCustomResourceDefinition, Resource<CustomResourceDefinition, DoneableCustomResourceDefinition>> customResourceDefinitions();

  <T extends HasMetadata, L extends KubernetesResourceList, D extends Doneable<T>> MixedOperation<T, L, D, Resource<T, D>> customResource(CustomResourceDefinition crd, Class<T> resourceType, Class<L> listClass, Class<D> doneClass);

  ExtensionsAPIGroupDSL extensions();

  AppsAPIGroupDSL apps();

  AutoscalingAPIGroupDSL autoscaling();

  MixedOperation<ComponentStatus, ComponentStatusList, DoneableComponentStatus, Resource<ComponentStatus, DoneableComponentStatus>> componentstatuses();

  ParameterNamespaceListVisitFromServerGetDeleteRecreateWaitApplicable<HasMetadata,Boolean> load(InputStream is);

  ParameterNamespaceListVisitFromServerGetDeleteRecreateWaitApplicable<HasMetadata, Boolean> resourceList(String s);

  NamespaceListVisitFromServerGetDeleteRecreateWaitApplicable<HasMetadata, Boolean> resourceList(KubernetesResourceList list);

  NamespaceListVisitFromServerGetDeleteRecreateWaitApplicable<HasMetadata, Boolean> resourceList(HasMetadata... items);

  NamespaceListVisitFromServerGetDeleteRecreateWaitApplicable<HasMetadata, Boolean> resourceList(Collection<HasMetadata> items);

  <T extends HasMetadata> NamespaceVisitFromServerGetWatchDeleteRecreateWaitApplicable<T ,Boolean> resource(T is);

  NamespaceVisitFromServerGetWatchDeleteRecreateWaitApplicable<HasMetadata,Boolean> resource(String s);

  MixedOperation<Endpoints, EndpointsList, DoneableEndpoints, Resource<Endpoints, DoneableEndpoints>> endpoints();

  MixedOperation<Event, EventList, DoneableEvent, Resource<Event, DoneableEvent>> events();

  NonNamespaceOperation< Namespace, NamespaceList, DoneableNamespace, Resource<Namespace, DoneableNamespace>> namespaces();

  NonNamespaceOperation<Node, NodeList, DoneableNode, Resource<Node, DoneableNode>> nodes();

  NonNamespaceOperation<PersistentVolume, PersistentVolumeList, DoneablePersistentVolume, Resource<PersistentVolume, DoneablePersistentVolume>> persistentVolumes();

  MixedOperation<PersistentVolumeClaim, PersistentVolumeClaimList, DoneablePersistentVolumeClaim, Resource<PersistentVolumeClaim, DoneablePersistentVolumeClaim>> persistentVolumeClaims();

  MixedOperation<Pod, PodList, DoneablePod, PodResource<Pod, DoneablePod>> pods();

  MixedOperation<ReplicationController, ReplicationControllerList, DoneableReplicationController, RollableScalableResource<ReplicationController, DoneableReplicationController>> replicationControllers();

  MixedOperation<ResourceQuota, ResourceQuotaList, DoneableResourceQuota, Resource<ResourceQuota, DoneableResourceQuota>> resourceQuotas();

  MixedOperation<Secret, SecretList, DoneableSecret, Resource<Secret, DoneableSecret>> secrets();

  MixedOperation<Service, ServiceList, DoneableService, Resource<Service, DoneableService>> services();

  MixedOperation<ServiceAccount, ServiceAccountList, DoneableServiceAccount, Resource<ServiceAccount, DoneableServiceAccount>> serviceAccounts();

  KubernetesListMixedOperation lists();
  
  MixedOperation<ConfigMap, ConfigMapList, DoneableConfigMap, Resource<ConfigMap, DoneableConfigMap>> configMaps();

  MixedOperation<LimitRange, LimitRangeList, DoneableLimitRange, Resource<LimitRange, DoneableLimitRange>> limitRanges();
}
