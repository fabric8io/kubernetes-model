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

package io.ucosty.kubernetes.client.dsl;

import io.ucosty.kubernetes.api.model.extensions.DaemonSet;
import io.ucosty.kubernetes.api.model.extensions.DaemonSetList;
import io.ucosty.kubernetes.api.model.extensions.Deployment;
import io.ucosty.kubernetes.api.model.extensions.DeploymentList;
import io.ucosty.kubernetes.api.model.extensions.DoneableDaemonSet;
import io.ucosty.kubernetes.api.model.extensions.DoneableDeployment;
import io.ucosty.kubernetes.api.model.extensions.DoneableIngress;
import io.ucosty.kubernetes.api.model.extensions.DoneableNetworkPolicy;
import io.ucosty.kubernetes.api.model.extensions.NetworkPolicy;
import io.ucosty.kubernetes.api.model.extensions.NetworkPolicyList;
import io.ucosty.kubernetes.api.model.DoneableJob;
import io.ucosty.kubernetes.api.model.extensions.DoneableReplicaSet;
import io.ucosty.kubernetes.api.model.extensions.DoneableThirdPartyResource;
import io.ucosty.kubernetes.api.model.extensions.Ingress;
import io.ucosty.kubernetes.api.model.extensions.IngressList;
import io.ucosty.kubernetes.api.model.Job;
import io.ucosty.kubernetes.api.model.JobList;
import io.ucosty.kubernetes.api.model.extensions.ReplicaSet;
import io.ucosty.kubernetes.api.model.extensions.ReplicaSetList;
import io.ucosty.kubernetes.api.model.extensions.ThirdPartyResource;
import io.ucosty.kubernetes.api.model.extensions.ThirdPartyResourceList;
import io.ucosty.kubernetes.client.Client;

public interface ExtensionsAPIGroupDSL extends Client {

  MixedOperation<Job, JobList, DoneableJob, ScalableResource<Job, DoneableJob>> jobs();

  MixedOperation<Deployment, DeploymentList, DoneableDeployment, ScalableResource<Deployment, DoneableDeployment>> deployments();

  @Deprecated
  MixedOperation<Ingress, IngressList, DoneableIngress, Resource<Ingress, DoneableIngress>> ingress();

  MixedOperation<Ingress, IngressList, DoneableIngress, Resource<Ingress, DoneableIngress>> ingresses();

  MixedOperation<NetworkPolicy, NetworkPolicyList, DoneableNetworkPolicy, Resource<NetworkPolicy, DoneableNetworkPolicy>> networkPolicies();

  MixedOperation<DaemonSet, DaemonSetList, DoneableDaemonSet, Resource<DaemonSet, DoneableDaemonSet>> daemonSets();

  NonNamespaceOperation<ThirdPartyResource, ThirdPartyResourceList, DoneableThirdPartyResource, Resource<ThirdPartyResource, DoneableThirdPartyResource>> thirdPartyResources();

  MixedOperation<ReplicaSet, ReplicaSetList, DoneableReplicaSet, RollableScalableResource<ReplicaSet, DoneableReplicaSet>> replicaSets();
}
