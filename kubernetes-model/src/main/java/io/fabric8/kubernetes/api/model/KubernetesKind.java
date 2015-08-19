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
package io.fabric8.kubernetes.api.model;

import io.fabric8.openshift.api.model.BuildConfigList;
import io.fabric8.openshift.api.model.BuildList;
import io.fabric8.openshift.api.model.ClusterPolicy;
import io.fabric8.openshift.api.model.ClusterPolicyBinding;
import io.fabric8.openshift.api.model.ClusterPolicyBindingList;
import io.fabric8.openshift.api.model.ClusterPolicyList;
import io.fabric8.openshift.api.model.ClusterRoleBinding;
import io.fabric8.openshift.api.model.ClusterRoleBindingList;
import io.fabric8.openshift.api.model.DeploymentConfigList;
import io.fabric8.openshift.api.model.Group;
import io.fabric8.openshift.api.model.GroupList;
import io.fabric8.openshift.api.model.Identity;
import io.fabric8.openshift.api.model.IdentityList;
import io.fabric8.openshift.api.model.ImageList;
import io.fabric8.openshift.api.model.ImageStreamList;
import io.fabric8.openshift.api.model.OAuthAccessToken;
import io.fabric8.openshift.api.model.OAuthAccessTokenList;
import io.fabric8.openshift.api.model.OAuthAuthorizeToken;
import io.fabric8.openshift.api.model.OAuthAuthorizeTokenList;
import io.fabric8.openshift.api.model.OAuthClient;
import io.fabric8.openshift.api.model.OAuthClientAuthorization;
import io.fabric8.openshift.api.model.OAuthClientAuthorizationList;
import io.fabric8.openshift.api.model.OAuthClientList;
import io.fabric8.openshift.api.model.Policy;
import io.fabric8.openshift.api.model.PolicyBinding;
import io.fabric8.openshift.api.model.PolicyBindingList;
import io.fabric8.openshift.api.model.PolicyList;
import io.fabric8.openshift.api.model.Role;
import io.fabric8.openshift.api.model.RoleBinding;
import io.fabric8.openshift.api.model.RoleBindingList;
import io.fabric8.openshift.api.model.RoleList;
import io.fabric8.openshift.api.model.RouteList;
import io.fabric8.openshift.api.model.TagEvent;
import io.fabric8.openshift.api.model.Template;
import io.fabric8.openshift.api.model.TemplateList;
import io.fabric8.openshift.api.model.User;
import io.fabric8.openshift.api.model.UserList;

import java.util.HashMap;
import java.util.Map;

public enum KubernetesKind {

    List(KubernetesList.class),
    ObjectMeta(ObjectMeta.class),
    PodList(PodList.class),
    ReplicationControllerList(ReplicationControllerList.class),
    ServiceList(ServiceList.class),
    Endpoints(Endpoints.class),
    EndpointsList(EndpointsList.class),
    EventList(EventList.class),
    Node(Node.class),
    NodeList(NodeList.class),
    EnvVar(EnvVar.class),
    Namespace(Namespace.class),
    NamespaceList(NamespaceList.class),
    PersistentVolume(PersistentVolume.class),
    PersistentVolumeList(PersistentVolumeList.class),
    PersistentVolumeClaim(PersistentVolumeClaim.class),
    PersistentVolumeClaimList(PersistentVolumeClaimList.class),
    ResourceQuota(ResourceQuota.class),
    ResourceQuotaList(ResourceQuotaList.class),
    Secret(Secret.class),
    SecretList(SecretList.class),
    SecurityContextConstraints(SecurityContextConstraints.class),
    SecurityContextConstraintsList(SecurityContextConstraintsList.class),
    ServiceAccount(ServiceAccount.class),
    ServiceAccountList(ServiceAccountList.class),
    Status(Status.class),
    Quantity(Quantity.class),
    BuildRequest(io.fabric8.openshift.api.model.BuildRequest.class),
    BuildList(BuildList.class),
    BuildConfigList(BuildConfigList.class),
    ImageList(ImageList.class),
    ImageStreamList(ImageStreamList.class),
    DeploymentConfigList(DeploymentConfigList.class),
    RouteList(RouteList.class),
    ContainerStatus(ContainerStatus.class),
    Template(Template.class),
    TemplateList(TemplateList.class),
    TagEvent(TagEvent.class),
    OAuthClient(OAuthClient.class),
    OAuthAccessToken(OAuthAccessToken.class),
    OAuthAuthorizeToken(OAuthAuthorizeToken.class),
    OAuthClientAuthorization(OAuthClientAuthorization.class),
    OAuthAccessTokenList(OAuthAccessTokenList.class),
    OAuthAuthorizeTokenList(OAuthAuthorizeTokenList.class),
    OAuthClientList(OAuthClientList.class),
    OAuthClientAuthorizationList(OAuthClientAuthorizationList.class),
    ClusterPolicy(ClusterPolicy.class),
    ClusterPolicyList(ClusterPolicyList.class),
    ClusterPolicyBinding(ClusterPolicyBinding.class),
    ClusterPolicyBindingList(ClusterPolicyBindingList.class),
    Policy(Policy.class),
    PolicyList(PolicyList.class),
    PolicyBinding(PolicyBinding.class),
    PolicyBindingList(PolicyBindingList.class),
    Role(Role.class),
    RoleList(RoleList.class),
    RoleBinding(RoleBinding.class),
    RoleBindingList(RoleBindingList.class),
    ClusterRoleBinding(ClusterRoleBinding.class),
    ClusterRoleBindingList(ClusterRoleBindingList.class),
    User(User.class),
    UserList(UserList.class),
    Group(Group.class),
    GroupList(GroupList.class),
    Identity(Identity.class),
    IdentityList(IdentityList.class);

    private static final Map<String, Class<? extends KubernetesResource>> map = new HashMap<>();

    static {
        for (KubernetesKind kind : KubernetesKind.values()) {
            map.put(kind.name(), kind.type);
        }
    }

    private final Class<? extends KubernetesResource> type;

    KubernetesKind(Class type) {
        this.type = type;
    }

    public Class getType() {
        return type;
    }

    public static Class<? extends KubernetesResource> getTypeForName(String name) {
        return map.get(name);
    }
}
