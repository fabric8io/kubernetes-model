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
import io.ucosty.kubernetes.api.model.RootPaths;
import io.ucosty.kubernetes.api.model.Secret;
import io.ucosty.kubernetes.api.model.SecretList;
import io.ucosty.kubernetes.api.model.Service;
import io.ucosty.kubernetes.api.model.ServiceAccount;
import io.ucosty.kubernetes.api.model.ServiceAccountList;
import io.ucosty.kubernetes.api.model.ServiceList;
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
import okhttp3.OkHttpClient;

import java.io.InputStream;
import java.net.URL;

public class AutoAdaptableKubernetesClient extends DefaultKubernetesClient {

  private KubernetesClient delegate;

  public AutoAdaptableKubernetesClient() throws KubernetesClientException {
    delegate = adapt(new DefaultKubernetesClient());
  }

  public AutoAdaptableKubernetesClient(OkHttpClient httpClient, Config config) throws KubernetesClientException {
    delegate = adapt(new DefaultKubernetesClient(httpClient, config));
  }

  public AutoAdaptableKubernetesClient(Config config) throws KubernetesClientException {
    delegate = adapt(new DefaultKubernetesClient(config));
  }

  public AutoAdaptableKubernetesClient(String masterUrl) throws KubernetesClientException {
    delegate = adapt(new DefaultKubernetesClient(masterUrl));
  }


  public static KubernetesClient adapt(KubernetesClient initial) {
    KubernetesClient result = initial;
    for (ExtensionAdapter<? extends KubernetesClient> adapter : Adapters.list(KubernetesClient.class)) {
      if (adapter.isAdaptable(result)) {
        result = adapter.adapt(result);
      }
    }
    return result;
  }

  @Override
  public NamespacedKubernetesClient inNamespace(String namespace) {
    Config updated = new ConfigBuilder(getConfiguration()).withNamespace(namespace).build();
    return new AutoAdaptableKubernetesClient(httpClient, updated);
  }

  @Override
  public NamespacedKubernetesClient inAnyNamespace() {
    return inNamespace(null);
  }

  @Override
  public ExtensionsAPIGroupDSL extensions() {
    return delegate.extensions();
  }

  @Override
  public AppsAPIGroupDSL apps() {
    return delegate.apps();
  }

  @Override
  public AutoscalingAPIGroupDSL autoscaling() {
    return delegate.autoscaling();
  }

  @Override
  public MixedOperation<ComponentStatus, ComponentStatusList, DoneableComponentStatus, Resource<ComponentStatus, DoneableComponentStatus>> componentstatuses() {
    return delegate.componentstatuses();
  }

  @Override
  public ParameterNamespaceListVisitFromServerGetDeleteRecreateWaitApplicable<HasMetadata, Boolean> load(InputStream is) {
    return delegate.load(is);
  }

  @Override
  public NamespaceListVisitFromServerGetDeleteRecreateWaitApplicable<HasMetadata, Boolean> resourceList(KubernetesResourceList is) {
    return delegate.resourceList(is);
  }

  @Override
  public ParameterNamespaceListVisitFromServerGetDeleteRecreateWaitApplicable<HasMetadata, Boolean> resourceList(String s) {
    return delegate.resourceList(s);
  }

  @Override
  public NamespaceVisitFromServerGetWatchDeleteRecreateWaitApplicable<HasMetadata, Boolean> resource(HasMetadata is) {
    return delegate.resource(is);
  }

  @Override
  public NamespaceVisitFromServerGetWatchDeleteRecreateWaitApplicable<HasMetadata, Boolean> resource(String s) {
    return delegate.resource(s);
  }

  @Override
  public MixedOperation<Endpoints, EndpointsList, DoneableEndpoints, Resource<Endpoints, DoneableEndpoints>> endpoints() {
    return delegate.endpoints();
  }

  @Override
  public MixedOperation<Event, EventList, DoneableEvent, Resource<Event, DoneableEvent>> events() {
    return delegate.events();
  }

  @Override
  public NonNamespaceOperation<Namespace, NamespaceList, DoneableNamespace, Resource<Namespace, DoneableNamespace>> namespaces() {
    return delegate.namespaces();
  }

  @Override
  public NonNamespaceOperation<Node, NodeList, DoneableNode, Resource<Node, DoneableNode>> nodes() {
    return delegate.nodes();
  }

  @Override
  public NonNamespaceOperation<PersistentVolume, PersistentVolumeList, DoneablePersistentVolume, Resource<PersistentVolume, DoneablePersistentVolume>> persistentVolumes() {
    return delegate.persistentVolumes();
  }

  @Override
  public MixedOperation<PersistentVolumeClaim, PersistentVolumeClaimList, DoneablePersistentVolumeClaim, Resource<PersistentVolumeClaim, DoneablePersistentVolumeClaim>> persistentVolumeClaims() {
    return delegate.persistentVolumeClaims();
  }

  @Override
  public MixedOperation<Pod, PodList, DoneablePod, PodResource<Pod, DoneablePod>> pods() {
    return delegate.pods();
  }

  @Override
  public MixedOperation<ReplicationController, ReplicationControllerList, DoneableReplicationController, RollableScalableResource<ReplicationController, DoneableReplicationController>> replicationControllers() {
    return delegate.replicationControllers();
  }

  @Override
  public MixedOperation<ResourceQuota, ResourceQuotaList, DoneableResourceQuota, Resource<ResourceQuota, DoneableResourceQuota>> resourceQuotas() {
    return delegate.resourceQuotas();
  }

  @Override
  public MixedOperation<Secret, SecretList, DoneableSecret, Resource<Secret, DoneableSecret>> secrets() {
    return delegate.secrets();
  }

  @Override
  public MixedOperation<Service, ServiceList, DoneableService, Resource<Service, DoneableService>> services() {
    return delegate.services();
  }

  @Override
  public MixedOperation<ServiceAccount, ServiceAccountList, DoneableServiceAccount, Resource<ServiceAccount, DoneableServiceAccount>> serviceAccounts() {
    return delegate.serviceAccounts();
  }

  @Override
  public KubernetesListMixedOperation lists() {
    return delegate.lists();
  }

  @Override
  public MixedOperation<ConfigMap, ConfigMapList, DoneableConfigMap, Resource<ConfigMap, DoneableConfigMap>> configMaps() {
    return delegate.configMaps();
  }

  @Override
  public MixedOperation<LimitRange, LimitRangeList, DoneableLimitRange, Resource<LimitRange, DoneableLimitRange>> limitRanges() {
    return delegate.limitRanges();
  }

  @Override
  public <C> Boolean isAdaptable(Class<C> type) {
    return delegate.isAdaptable(type);
  }

  @Override
  public <C> C adapt(Class<C> type) {
    return delegate.adapt(type);
  }

  @Override
  public URL getMasterUrl() {
    return delegate.getMasterUrl();
  }

  @Override
  public String getApiVersion() {
    return delegate.getApiVersion();
  }

  @Override
  public String getNamespace() {
    return delegate.getNamespace();
  }

  @Override
  public RootPaths rootPaths() {
    return delegate.rootPaths();
  }

  @Override
  public void close() {
    delegate.close();
  }

  @Override
  public Config getConfiguration() {
    return delegate.getConfiguration();
  }
}
